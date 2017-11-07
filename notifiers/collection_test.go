package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	// go test makes a cd into package directory when testing. We must go up by one level to load our testdata
	configuration.UseFile("../testdata/test_config.yaml")
	err := configuration.Load()
	configuration.Config.Notifications.Enabled = true

	if err != nil {
		fmt.Print("Could not load test configuration file: ", err)
	}
}

func TestAddNotifier(t *testing.T) {
	notifiersLength := len(notifiersCollection)
	m := MockNotifier{}
	AddNotifier(m)

	if len(notifiersCollection) != notifiersLength+1 {
		t.Error("Expected ", notifiersLength+1, " Notifiers, got ", len(notifiersCollection))
	}
}

func TestSendNotification(t *testing.T) {
	n := 2

	notifiersCollection = []Notifier{}
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
	notifiersCollection = []Notifier{}

	show := TvShow{
		Id:   1,
		Name: "Test TVShow",
	}

	episode := Episode{
		Id:     30,
		Season: 3,
		Number: 4,
		Name:   "Test Episode S03E04",
		Date:   time.Now(),
	}

	m := MockNotifier{}
	AddNotifier(m)

	mockNotificationCounter = 0

	NotifyRecentEpisode(show, episode)

	if mockNotificationCounter != 1 {
		t.Error("Expected 1 notification to be sent, got ", mockNotificationCounter)
	}

	NotifyRecentEpisode(show, episode)

	if mockNotificationCounter != 1 {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}
}
