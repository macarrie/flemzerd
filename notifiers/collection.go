// Package notifier groups methods for notifications handling in flemzerd.
// Multiple notifier types can be registered. Sending a notification will then send it with all registered notifiers.
// Helper methods are available to send notifications for specific purposes (download start, new item from watchlist, ...)
package notifier

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
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
	var aggregatedErrorMessage bytes.Buffer

	for _, notifier := range notifiersCollection {
		mod, notifierAliveError := notifier.Status()
		if notifierAliveError != nil {
			log.WithFields(log.Fields{
				"error": notifierAliveError,
			}).Warning("Notifier is not alive")
			aggregatedErrorMessage.WriteString(notifierAliveError.Error())
			aggregatedErrorMessage.WriteString("\n")
		}
		modList = append(modList, mod)
	}

	var retError error
	if aggregatedErrorMessage.Len() == 0 {
		retError = nil
	} else {
		retError = errors.New(aggregatedErrorMessage.String())
	}
	return modList, retError
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

	notificationTitle := fmt.Sprintf("%v: New episode aired (S%03dE%03d)", episode.TvShow.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("New episode aired on %v\n%v Season %03d Episode %03d: %v", episode.Date, episode.TvShow.Name, episode.Season, episode.Number, episode.Name)

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
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

	notificationTitle := fmt.Sprintf("%v: Download start (S%03dE%03d)", episode.TvShow.Name, episode.Season, episode.Number)
	notificationContent := "Torrents found for episode. Starting download"

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
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

	notificationTitle := fmt.Sprintf("%s", m.Title)
	notificationContent := "Movie found in watchlist, adding to tracked movies"

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
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

	notificationTitle := fmt.Sprintf("%v: Download start", m.Title)
	notificationContent := "Torrents found for movie. Starting download"

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
	}

	return nil
}

// NotifyDownloadedEpisode sends notification on registered notifiers to alert that the episode has been successfully downloaded
func NotifyDownloadedEpisode(episode *Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadComplete {
		return nil
	}

	notificationTitle := fmt.Sprintf("%v: Episode downloaded (S%03dE%03d)", episode.TvShow.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("New episode downloaded\n%v Season %03d Episode %03d: %v", episode.TvShow.Name, episode.Season, episode.Number, episode.Name)

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
	}

	return nil
}

// NotifyDownloadedMovie sends notification on registered notifiers to alert that the movie has been successfully downloaded
func NotifyDownloadedMovie(m *Movie) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadComplete {
		return nil
	}

	notificationTitle := fmt.Sprintf("%v: Movie downloaded", m.Title)
	notificationContent := "New movie downloaded\n"

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
	}

	return nil
}

// NotifyFailedEpisode sends notification on registered notifiers to alert that the episode could not be downloaded.
// An episode is marked as failed when more that TorrentDownloadAttempts configuration parameter) torrent downloads have failed
func NotifyFailedEpisode(episode *Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyFailure {
		return nil
	}

	notificationTitle := fmt.Sprintf("%v: Episode download failed (S%03dE%03d)", episode.TvShow.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("Failed to download episode\n%v Season %03d Episode %03d: %v", episode.TvShow.Name, episode.Season, episode.Number, episode.Name)

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
	}

	return nil
}

// NotifyFailedMovie sends notification on registered notifiers to alert that the movie could not be downloaded.
// A movie is marked as failed when more that TorrentDownloadAttempts configuration parameter) torrent downloads have failed
func NotifyFailedMovie(m *Movie) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyFailure {
		return nil
	}

	notificationTitle := fmt.Sprintf("%v: Movie download failed", m.Title)
	notificationContent := "Failed to download movie\n"

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
	}

	return nil
}

// SendNotification sends the notification with title and content using all registered notifiers.
// If at least one notifier returns an error when sending the notification, the method exists with a non nil error
func SendNotification(title, content string) error {
	if !configuration.Config.Notifications.Enabled {
		return nil
	}

	var noNotificationSent bool
	noNotificationSent = true
	for _, notifier := range notifiersCollection {
		err := notifier.Send(title, content)
		if err == nil {
			noNotificationSent = false
		}
	}

	if noNotificationSent {
		return errors.New("Could not send any notification")
	} else {
		return nil
	}
}
