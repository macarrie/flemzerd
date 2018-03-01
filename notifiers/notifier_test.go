package notifier

import (
	. "github.com/macarrie/flemzerd/objects"
)

type MockNotifier struct{}

var mockNotificationCounter int
var notifierInitialized bool

func (n MockNotifier) Status() (Module, error) {
	return Module{}, nil
}

func (n MockNotifier) Setup(authCredentials map[string]string) {
	return
}

func (n MockNotifier) Init() error {
	notifierInitialized = true
	return nil
}

func (n MockNotifier) Send(title, content string) error {
	mockNotificationCounter += 1
	return nil
}
