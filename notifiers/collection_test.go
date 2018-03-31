package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/retention"
)

func init() {
	// go test makes a cd into package directory when testing. We must go up by one level to load our testdata
	configuration.UseFile("../testdata/test_config.yaml")
	err := configuration.Load()
	configuration.Config.Notifications.Enabled = true

	if err != nil {
		fmt.Print("Could not load test configuration file: ", err)
	}

	retention.InitStruct()
}

func TestStatus(t *testing.T) {
	n1 := MockNotifier{}
	n2 := MockNotifier{}

	notifiersCollection = []Notifier{n1, n2}

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 notifier modules status, got %d instead", len(mods))
	}
	if err == nil {
		t.Error("Expected to have aggregated error for notifier status")
	}
}

func TestReset(t *testing.T) {
	n := MockNotifier{}
	AddNotifier(n)
	Reset()

	if len(notifiersCollection) != 0 {
		t.Error("Expected notifier collection to be empty after reset")
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

func TestNotifyEpisode(t *testing.T) {
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
		t.Error("Expected recent notification to be sent, got none")
	}

	NotifyDownloadedEpisode(show, episode)
	if mockNotificationCounter != 2 {
		t.Error("Expected downloaded notification to be sent, got none")
	}

	NotifyFailedEpisode(show, episode)
	if mockNotificationCounter != 3 {
		t.Error("Expected failed notification to be sent, got none")
	}

	NotifyRecentEpisode(show, episode)
	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}

	configuration.Config.Notifications.Enabled = false

	NotifyRecentEpisode(show, episode)
	NotifyDownloadedEpisode(show, episode)
	NotifyFailedEpisode(show, episode)
	NotifyRecentEpisode(show, episode)

	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because notifications are disabled, but notifications have been sent anyway")
	}

	configuration.Config.Notifications.Enabled = true
}

func TestNotifyMovie(t *testing.T) {
	notifiersCollection = []Notifier{}

	movie := Movie{
		Id:    1,
		Title: "Test movie",
	}

	m := MockNotifier{}
	AddNotifier(m)

	mockNotificationCounter = 0

	NotifyMovieDownload(movie)
	if mockNotificationCounter != 1 {
		t.Error("Expected movie notification to be sent, got none")
	}

	NotifyDownloadedMovie(movie)
	if mockNotificationCounter != 2 {
		t.Error("Expected downloaded movie notification to be sent, got none")
	}

	NotifyFailedMovie(movie)
	if mockNotificationCounter != 3 {
		t.Error("Expected failed notification to be sent, got none")
	}

	NotifyMovieDownload(movie)
	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}

	configuration.Config.Notifications.Enabled = false

	NotifyMovieDownload(movie)
	NotifyDownloadedMovie(movie)
	NotifyFailedMovie(movie)
	NotifyMovieDownload(movie)

	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because notifications are disabled, but notifications have been sent anyway")
	}

	configuration.Config.Notifications.Enabled = true
}
