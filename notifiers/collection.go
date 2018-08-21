// Package notifier groups methods for notifications handling in flemzerd.
// Multiple notifier types can be registered. Sending a notification will then send it with all registered notifiers.
// Helper methods are available to send notifications for specific purposes (download start, new item from watchlist, ...)
package notifier

import (
	"fmt"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	multierror "github.com/hashicorp/go-multierror"
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

// NotifyEpisodeDownloadStart sends a notification to alert that torrents have been found for episode and that download process is starting.
func NotifyEpisodeDownloadStart(episode *Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadStart {
		return nil
	}

	if err := SendNotification(Notification{
		Type:    NOTIFICATION_DOWNLOAD_START,
		Episode: *episode,
	}); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

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
func NotifyMovieDownloadStart(m *Movie) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadStart {
		return nil
	}

	if err := SendNotification(Notification{
		Type:  NOTIFICATION_DOWNLOAD_START,
		Movie: *m,
	}); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// NotifyDownloadedEpisode sends notification on registered notifiers to alert that the episode has been successfully downloaded
func NotifyDownloadedEpisode(episode *Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadComplete {
		return nil
	}

	if err := SendNotification(Notification{
		Type:    NOTIFICATION_DOWNLOAD_SUCCESS,
		Episode: *episode,
	}); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// NotifyDownloadedMovie sends notification on registered notifiers to alert that the movie has been successfully downloaded
func NotifyDownloadedMovie(m *Movie) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadComplete {
		return nil
	}

	if err := SendNotification(Notification{
		Type:  NOTIFICATION_DOWNLOAD_SUCCESS,
		Movie: *m,
	}); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// NotifyFailedEpisode sends notification on registered notifiers to alert that the episode could not be downloaded.
// An episode is marked as failed when more that TorrentDownloadAttempts configuration parameter) torrent downloads have failed
func NotifyFailedEpisode(episode *Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyFailure {
		return nil
	}

	if err := SendNotification(Notification{
		Type:    NOTIFICATION_DOWNLOAD_FAILURE,
		Episode: *episode,
	}); err != nil {
		return errors.Wrap(err, "Errors detected when sending notification")
	}

	return nil
}

// NotifyFailedMovie sends notification on registered notifiers to alert that the movie could not be downloaded.
// A movie is marked as failed when more that TorrentDownloadAttempts configuration parameter) torrent downloads have failed
func NotifyFailedMovie(m *Movie) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyFailure {
		return nil
	}

	if err := SendNotification(Notification{
		Type:  NOTIFICATION_DOWNLOAD_FAILURE,
		Movie: *m,
	}); err != nil {
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
