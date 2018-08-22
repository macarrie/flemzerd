package mock

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type Notifier struct{}
type ErrorNotifier struct{}

var mockNotificationCounter int
var notifierInitialized bool

func (n Notifier) Status() (Module, error) {
	return Module{
		Name: "Notifier",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (n ErrorNotifier) Status() (Module, error) {
	var err error = fmt.Errorf("Notifier error")
	return Module{
		Name: "ErrorNotifier",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (n Notifier) GetName() string {
	return "Notifier"
}
func (n ErrorNotifier) GetName() string {
	return "ErrorNotifier"
}

func (n Notifier) Setup(authCredentials map[string]string) {
	return
}
func (n ErrorNotifier) Setup(authCredentials map[string]string) {
	return
}

func (n Notifier) Init() error {
	notifierInitialized = true
	return nil
}

func (n ErrorNotifier) Init() error {
	return nil
}

func (n Notifier) Send(notif Notification) error {
	mockNotificationCounter += 1
	return nil
}

func (n ErrorNotifier) Send(notif Notification) error {
	return fmt.Errorf("Error when sending notification")
}

func (n Notifier) GetNotificationCount() int {
	return mockNotificationCounter
}
func (n ErrorNotifier) GetNotificationCount() int {
	return mockNotificationCounter
}
