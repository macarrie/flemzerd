package notifier

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/retention"
)

var notifiersCollection []Notifier
var Retention []int

func AddNotifier(notifier Notifier) {
	notifiersCollection = append(notifiersCollection, notifier)
	log.WithFields(log.Fields{
		"notifier": notifier,
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

func NotifyRecentEpisode(show TvShow, episode Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyNewEpisode {
		return nil
	}

	retention.CleanOldNotifiedEpisodes()

	if retention.HasBeenNotified(episode) {
		return nil
	}

	notificationTitle := fmt.Sprintf("%v: New episode aired (S%03dE%03d)", show.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("New episode aired on %v\n%v Season %03d Episode %03d: %v", episode.Date, show.Name, episode.Season, episode.Number, episode.Name)

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
	}

	retention.AddNotifiedEpisode(episode)

	return nil
}

func NotifyDownloadedEpisode(show TvShow, episode Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyDownloadComplete {
		return nil
	}

	notificationTitle := fmt.Sprintf("%v: Episode downloaded (S%03dE%03d)", show.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("New episode downloaded\n%v Season %03d Episode %03d: %v", show.Name, episode.Season, episode.Number, episode.Name)

	err := SendNotification(notificationTitle, notificationContent)
	if err != nil {
		return err
	}

	return nil
}

func NotifyFailedEpisode(show TvShow, episode Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyFailure {
		return nil
	}

	notificationTitle := fmt.Sprintf("%v: Episode download failed (S%03dE%03d)", show.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("Failed to download episode\n%v Season %03d Episode %03d: %v", show.Name, episode.Season, episode.Number, episode.Name)

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
