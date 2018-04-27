package notifier

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type MockNotifier struct{}
type MockErrorNotifier struct{}

var mockNotificationCounter int
var notifierInitialized bool

func (n MockNotifier) Status() (Module, error) {
	return Module{
		Name: "MockNotifier",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (n MockErrorNotifier) Status() (Module, error) {
	var err error = fmt.Errorf("Notifier error")
	return Module{
		Name: "MockErrorNotifier",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (n MockNotifier) Setup(authCredentials map[string]string) {
	return
}
func (n MockErrorNotifier) Setup(authCredentials map[string]string) {
	return
}

func (n MockNotifier) Init() error {
	notifierInitialized = true
	return nil
}

func (n MockErrorNotifier) Init() error {
	return nil
}

func (n MockNotifier) Send(title, content string) error {
	mockNotificationCounter += 1
	return nil
}

func (n MockErrorNotifier) Send(title, content string) error {
	return fmt.Errorf("Error when sending notification")
}
