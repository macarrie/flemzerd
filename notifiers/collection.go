package notifier

import (
	"errors"
	"fmt"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

var notifiersCollection []Notifier
var Retention []int

func AddNotifier(notifier Notifier) {
	notifiersCollection = append(notifiersCollection, notifier)
	log.WithFields(log.Fields{
		"notifier": notifier,
	}).Debug("Notifier loaded")
}

func RemoveFromRetention(idToRemove int) {
	var newRetention []int

	for _, episodeId := range Retention {
		if episodeId != idToRemove {
			newRetention = append(newRetention, episodeId)
		}
	}

	Retention = newRetention
}

func NotifyRecentEpisode(show TvShow, episode Episode) error {
	if !configuration.Config.Notifications.Enabled || !configuration.Config.Notifications.NotifyNewEpisode {
		return nil
	}

	for _, episodeId := range Retention {
		if episode.Date.Before(time.Now().AddDate(0, 0, -14)) {
			RemoveFromRetention(episodeId)
		}
	}

	notificationTitle := fmt.Sprintf("%v: New episode aired (S%03dE%03d)", show.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("New episode aired on %v\n%v Season %03d Episode %03d: %v", episode.Date, show.Name, episode.Season, episode.Number, episode.Name)

	alreadyNotified := false
	for _, retentionEpisodeId := range Retention {
		if retentionEpisodeId == episode.Id {
			alreadyNotified = true

			break
		}
	}

	if alreadyNotified {
		log.Debug("Notifications already sent for episode. Nothing to do")
		return nil
	} else {
		err := SendNotification(notificationTitle, notificationContent)
		if err != nil {
			return err
		}

		Retention = append(Retention, episode.Id)

		return nil
	}
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
