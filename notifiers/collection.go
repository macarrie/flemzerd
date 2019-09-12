// Package notifier groups methods for notifications handling in flemzerd.
// Multiple notifier types can be registered. Sending a notification will then send it with all registered notifiers.
// Helper methods are available to send notifications for specific purposes (download start, new item from watchlist, ...)
package notifier

import (
	"fmt"

	"github.com/macarrie/flemzerd/downloadable"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/stats"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

var notifiersCollection []Notifier

// AddNotifier registers a new notifier
func AddNotifier(notifier Notifier) {
	notifiersCollection = append(notifiersCollection, notifier)
	log.WithFields(log.Fields{
		"notifier": notifier.GetName(),
	}).Debug("Notifier loaded")
}

// Status checks registered notifiers status. A module list is returned, each module corresponds to a registered notifier. A non nil error is returned if at least one registered notifier is in error
func Status() ([]Module, error) {
	var modList []Module
	var errorList *multierror.Error

	for _, notifier := range notifiersCollection {
		mod, notifierAliveError := notifier.Status()
		if notifierAliveError != nil {
			log.WithFields(log.Fields{
				"error": notifierAliveError,
			}).Warning("Notifier is not alive")
			errorList = multierror.Append(errorList, notifierAliveError)
		}
		modList = append(modList, mod)
	}

	return modList, errorList.ErrorOrNil()
}

// Reset empties registered notifiers list
func Reset() {
	notifiersCollection = []Notifier{}
}

// NotifyRecentEpisode sends a notification on all registered notifiers to alert that a new episode for a tracked show has been released.
// The episode is then marked as notified and the notification will not be sent again if this method is called twice on the same episode.
func NotifyRecentEpisode(episode *Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyNewEpisode {
		return nil
	}

	if episode.Notified {
		return nil
	}

	if err := SendNotification(Notification{
		Type:    NOTIFICATION_NEW_EPISODE,
		Episode: *episode,
	}); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	episode.Notified = true
	db.Client.Save(&episode)

	return nil
}

// NotifyNewMovie sends a notification on all registered notifiers to alert that a new movie has been add in watchlists
// The movie is then marked as notified and the notification will not be sent again if this method is called twice on the same episode.
func NotifyNewMovie(m *Movie) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyNewMovie {
		return nil
	}

	if m.Notified {
		return nil
	}

	if err := SendNotification(Notification{
		Type:  NOTIFICATION_NEW_MOVIE,
		Movie: *m,
	}); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	m.Notified = true
	db.Client.Save(&m)

	return nil
}

// NotifyMovieDownloadStart sends a notification to alert that torrents have been found for movie and that download process is starting.
func NotifyDownloadStart(d downloadable.Downloadable) error {
	notification := Notification{}
	switch d.(type) {
	case *Movie:
		notification = Notification{
			Type:  NOTIFICATION_DOWNLOAD_START,
			Movie: *(d.(*Movie)),
		}
		stats.Stats.Movies.Downloading += 1
	case *Episode:
		notification = Notification{
			Type:    NOTIFICATION_DOWNLOAD_START,
			Episode: *(d.(*Episode)),
		}
		stats.Stats.Episodes.Downloading += 1
	}

	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadStart {
		return nil
	}

	if err := SendNotification(notification); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// NotifyDownloadedItem sends notification on registered notifiers to alert that the movie has been successfully downloaded
func NotifyDownloadedItem(d downloadable.Downloadable) error {
	notification := Notification{}
	switch d.(type) {
	case *Movie:
		notification = Notification{
			Type:  NOTIFICATION_DOWNLOAD_SUCCESS,
			Movie: *d.(*Movie),
		}
		stats.Stats.Movies.Downloaded += 1
		stats.Stats.Movies.Downloading -= 1
	case *Episode:
		notification = Notification{
			Type:    NOTIFICATION_DOWNLOAD_SUCCESS,
			Episode: *d.(*Episode),
		}
		stats.Stats.Episodes.Downloaded += 1
		stats.Stats.Episodes.Downloading -= 1
	}

	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadComplete {
		return nil
	}

	if err := SendNotification(notification); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// NotifyFailedDownload sends notification on registered notifiers to alert that the movie could not be downloaded.
// A movie is marked as failed when more that TorrentDownloadAttempts configuration parameter) torrent downloads have failed
func NotifyFailedDownload(d downloadable.Downloadable) error {
	notification := Notification{}
	dlItem := d.GetDownloadingItem()
	switch d.(type) {
	case *Movie:
		notification = Notification{
			Type:  NOTIFICATION_DOWNLOAD_FAILURE,
			Movie: *d.(*Movie),
		}
		if dlItem.Downloading {
			stats.Stats.Movies.Downloading -= 1
		}
	case *Episode:
		notification = Notification{
			Type:    NOTIFICATION_DOWNLOAD_FAILURE,
			Episode: *d.(*Episode),
		}
		if dlItem.Downloading {
			stats.Stats.Episodes.Downloading -= 1
		}
	}

	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyFailure {
		return nil
	}

	if err := SendNotification(notification); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// NotifyTorrentNotFound sends notification on registered notifiers to alert that torrents could not be found (either no torrents found or no available indexers)
func NotifyTorrentsNotFound(d downloadable.Downloadable) error {
	notification := Notification{}
	downloadingItem := d.GetDownloadingItem()

	switch d.(type) {
	case *Movie:
		notification = Notification{
			Type:  NOTIFICATION_NO_TORRENTS,
			Movie: *(d.(*Movie)),
		}
	case *Episode:
		notification = Notification{
			Type:    NOTIFICATION_NO_TORRENTS,
			Episode: *(d.(*Episode)),
		}
	default:
		d.GetLog().Debug("Unknown Downloadable object type wen sending torrent not found notification")
		return nil
	}

	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyFailure {
		return nil
	}

	downloadingItem.TorrentsNotFound = true
	d.SetDownloadingItem(downloadingItem)
	db.SaveDownloadable(&d)

	if err := SendNotification(notification); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Could not send 'no torrents found' notification")
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// SendNotification sends the notification with title and content using all registered notifiers.
// If at least one notifier returns an error when sending the notification, the method exists with a non nil error
func SendNotification(notif Notification) error {
	if !configuration.Config.Notifications.Enabled {
		return nil
	}

	var sendingErrors *multierror.Error
	var noNotificationSent bool
	noNotificationSent = true
	for _, notifier := range notifiersCollection {
		if err := notifier.Send(notif); err != nil {
			sendingErrors = multierror.Append(sendingErrors, err)
		} else {
			noNotificationSent = false
		}
	}

	if noNotificationSent {
		return sendingErrors
	}

	stats.Stats.Notifications.Unread += 1
	return nil
}

// GetNotifier returns the registered notifier with name "name". An non-nil error is returned if no registered notifier are found with the required name
func GetNotifier(name string) (Notifier, error) {
	for _, n := range notifiersCollection {
		mod, _ := n.Status()
		if mod.Name == name {
			return n, nil
		}
	}

	return nil, fmt.Errorf("Notifier %s not found in configuration", name)
}
