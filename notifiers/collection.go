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

func AddNotifier(notifier Notifier) {
	notifiersCollection = append(notifiersCollection, notifier)
	log.WithFields(log.Fields{
		"notifier": notifier.GetName(),
	}).Debug("Notifier loaded")
}

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

func Reset() {
	notifiersCollection = []Notifier{}
}

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

func SendNotification(title, content string) error {
	if !configuration.Config.Notifications.Enabled {
		return nil
	}

	var sendingErrors bool
	for _, notifier := range notifiersCollection {
		err := notifier.Send(title, content)
		if err != nil {
			sendingErrors = true
		}
	}

	if sendingErrors {
		return errors.New("Couldn't send notifications for all notifiers")
	} else {
		return nil
	}
}
