package notifier

import (
	"errors"
	log "github.com/macarrie/flemzerd/logging"
)

type Notifier interface {
	Setup(authCredentials map[string]string)
	Init() error
	Send(title, content string) error
}

var notifiers []Notifier
var Retention []int

func Init() {
	for _, notifier := range notifiers {
		notifier.Init()
	}
}

func AddNotifier(notifier Notifier) {
	notifiers = append(notifiers, notifier)
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

func NotifyRecentEpisode(episodeId int, title string, content string) error {
	alreadyNotified := false

	for _, retentionEpisodeId := range Retention {
		if retentionEpisodeId == episodeId {
			alreadyNotified = true

			break
		}
	}

	if alreadyNotified {
		return errors.New("Notifications already sent for episode. Nothing to do")
	} else {
		err := SendNotification(title, content)
		if err != nil {
			return err
		}

		Retention = append(Retention, episodeId)

		return nil
	}
}

func SendNotification(title, content string) error {
	var sendingErrors bool
	for _, notifier := range notifiers {
		err := notifier.Send(title, content)
		if err != nil {
			sendingErrors = true
		}
	}

	if sendingErrors {
		return errors.New("Couldn't send all notifications")
	} else {
		return nil
	}
}
