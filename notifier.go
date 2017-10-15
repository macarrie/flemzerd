package main

import (
	"errors"
	log "flemzerd/logging"
)

type Notifier interface {
	Setup(authCredentials map[string]string)
	Init() bool
	Send(title, content string) error
}

var notifiers []Notifier

func Init() {
	for _, notifier := range notifiers {
		notifier.Init()
	}
}

func AddNotifier(notifier Notifier) {
	notifiers = append(notifiers, notifier)
	log.WithFields(log.Fields{
		"notifier": notifier,
	}).Info("Notifier loaded")
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
