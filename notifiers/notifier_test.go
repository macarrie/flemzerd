package notifier

import (
	"testing"
)

type MockNotifier struct{}

var mockNotificationCounter int

func (n MockNotifier) Setup(authCredentials map[string]string) {
	return
}

func (n MockNotifier) Init() error {
	return nil
}

func (n MockNotifier) Send(title, content string) error {
	mockNotificationCounter += 1
	return nil
}

func TestAddNotifier(t *testing.T) {
	notifiersLength := len(notifiers)
	m := MockNotifier{}
	AddNotifier(m)

	if len(notifiers) != notifiersLength+1 {
		t.Error("Expected ", notifiersLength+1, " Notifiers, got ", len(notifiers))
	}
}

func TestRemoveFromRetention(t *testing.T) {
	Retention = []int{1, 2}
	itemToRemove := 2

	RemoveFromRetention(itemToRemove)

	removed := true
	for _, item := range Retention {
		if item == itemToRemove {
			removed = false
		}
	}

	if !removed {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from retention but it is still present")
	}
}

func TestSendNotification(t *testing.T) {
	n := 2

	notifiers = []Notifier{}
	mockNotificationCounter = 0
	mockNotifiers := make([]MockNotifier, n)

	for i := range mockNotifiers {
		mockNotifiers[i] = MockNotifier{}
		AddNotifier(mockNotifiers[i])
	}

	SendNotification("Title", "Content")

	if mockNotificationCounter != n {
		t.Error("Expected to send ", n, " notifications, but ", mockNotificationCounter, " notifications have been sent")
	}
}

func TestNotifyRecentEpisode(t *testing.T) {
	Retention = []int{}
	notifiers = []Notifier{}

	m := MockNotifier{}
	AddNotifier(m)

	mockNotificationCounter = 0

	NotifyRecentEpisode(1, "Test title", "Test content")

	if mockNotificationCounter != 1 {
		t.Error("Expected 1 notification to be sent, got ", mockNotificationCounter)
	}

	err := NotifyRecentEpisode(1, "Test title", "Test content")

	if mockNotificationCounter != 1 || err == nil {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}
}
