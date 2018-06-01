package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
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

	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
}

func TestStatus(t *testing.T) {
	n1 := MockNotifier{}
	notifiersCollection = []Notifier{n1}
	_, err := Status()
	if err != nil {
		t.Error("Expected to have no error for notifier status")
	}

	n2 := MockErrorNotifier{}
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

	AddNotifier(MockErrorNotifier{})
	err := SendNotification("Title", "Content")
	if err != nil {
		t.Error("Expected to have no error when sending notifications")
	}

	notifiersCollection = []Notifier{MockErrorNotifier{}}
	err = SendNotification("Title", "Content")
	if err == nil {
		t.Error("Expected to have an error when sending notifications")
	}

	prev := mockNotificationCounter
	configuration.Config.Notifications.Enabled = false
	SendNotification("Title", "Content")
	if mockNotificationCounter != prev {
		t.Error("Expected notifications not to be sent because notifications are disabled in configuration")
	}
	configuration.Config.Notifications.Enabled = true
}

func TestNotifyEpisode(t *testing.T) {
	notifiersCollection = []Notifier{}

	show := TvShow{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "Test TVShow",
	}

	episode := Episode{
		Model: gorm.Model{
			ID: 30,
		},
		TvShow: show,
		Season: 3,
		Number: 4,
		Name:   "Test Episode S03E04",
		Date:   time.Now(),
	}

	m := MockNotifier{}
	AddNotifier(m)

	mockNotificationCounter = 0

	NotifyRecentEpisode(&episode)
	if mockNotificationCounter != 1 {
		t.Error("Expected recent notification to be sent, got none")
	}

	NotifyDownloadedEpisode(&episode)
	if mockNotificationCounter != 2 {
		t.Error("Expected downloaded notification to be sent, got none")
	}

	NotifyFailedEpisode(&episode)
	if mockNotificationCounter != 3 {
		t.Error("Expected failed notification to be sent, got none")
	}

	NotifyRecentEpisode(&episode)
	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}

	configuration.Config.Notifications.Enabled = false

	NotifyRecentEpisode(&episode)
	NotifyDownloadedEpisode(&episode)
	NotifyFailedEpisode(&episode)
	NotifyRecentEpisode(&episode)

	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because notifications are disabled, but notifications have been sent anyway")
	}

	configuration.Config.Notifications.Enabled = true

	notifiersCollection = []Notifier{MockErrorNotifier{}}

	episode.Notified = false
	err := NotifyRecentEpisode(&episode)
	if err == nil {
		t.Error("Expected to have error while sending notification for recent episode, got none")
	}
	err = NotifyDownloadedEpisode(&episode)
	if err == nil {
		t.Error("Expected to have error while sending notification for downloaded episode, got none")
	}
	err = NotifyFailedEpisode(&episode)
	if err == nil {
		t.Error("Expected to have error while sending notification for failed episode, got none")
	}
}

func TestNotifyMovie(t *testing.T) {
	notifiersCollection = []Notifier{}

	movie := Movie{
		Model: gorm.Model{
			ID: 1,
		},
		Title: "Test movie",
	}

	m := MockNotifier{}
	AddNotifier(m)

	mockNotificationCounter = 0

	NotifyNewMovie(&movie)
	if mockNotificationCounter != 1 {
		t.Error("Expected movie notification to be sent, got none")
	}

	NotifyDownloadedMovie(&movie)
	if mockNotificationCounter != 2 {
		t.Error("Expected downloaded movie notification to be sent, got none")
	}

	NotifyFailedMovie(&movie)
	if mockNotificationCounter != 3 {
		t.Error("Expected failed notification to be sent, got none")
	}

	NotifyNewMovie(&movie)
	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}

	configuration.Config.Notifications.Enabled = false

	NotifyNewMovie(&movie)
	NotifyDownloadedMovie(&movie)
	NotifyFailedMovie(&movie)
	NotifyNewMovie(&movie)

	if mockNotificationCounter != 3 {
		t.Error("Expected notification not to be sent because notifications are disabled, but notifications have been sent anyway")
	}

	notifiersCollection = []Notifier{MockErrorNotifier{}}
	configuration.Config.Notifications.Enabled = true

	movie.Notified = false
	err := NotifyNewMovie(&movie)
	if err == nil {
		t.Error("Expected to have error while sending notification for recent movie, got none")
	}
	err = NotifyDownloadedMovie(&movie)
	if err == nil {
		t.Error("Expected to have error while sending notification for downloaded movie, got none")
	}
	err = NotifyFailedMovie(&movie)
	if err == nil {
		t.Error("Expected to have error while sending notification for failed movie, got none")
	}
}

func TestNotifyDownloadStart(t *testing.T) {
	notifiersCollection = []Notifier{}
	AddNotifier(MockNotifier{})

	show := TvShow{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "Test TVShow",
	}

	episode := Episode{
		Model: gorm.Model{
			ID: 30,
		},
		TvShow: show,
		Season: 3,
		Number: 4,
		Name:   "Test Episode S03E04",
		Date:   time.Now(),
	}

	movie := Movie{
		Model: gorm.Model{
			ID: 1,
		},
		Title: "Test movie",
	}

	count := mockNotificationCounter
	NotifyEpisodeDownloadStart(&episode)
	if mockNotificationCounter != count+1 {
		t.Error("Expected notification to be sent when notifying episode download start")
	}

	count = mockNotificationCounter
	NotifyMovieDownloadStart(&movie)
	if mockNotificationCounter != count+1 {
		t.Error("Expected notification to be sent when notifying movie download start")
	}

	configuration.Config.Notifications.NotifyDownloadStart = false
	count = mockNotificationCounter
	NotifyEpisodeDownloadStart(&episode)
	if mockNotificationCounter != count {
		t.Error("Expected notification not to be sent when notifying episode download start because of configuration params")
	}

	count = mockNotificationCounter
	NotifyMovieDownloadStart(&movie)
	if mockNotificationCounter != count {
		t.Error("Expected notification not to be sent when notifying movie download start because of configuration params")
	}

	configuration.Config.Notifications.NotifyDownloadStart = true
	notifiersCollection = []Notifier{MockErrorNotifier{}}
	err := NotifyEpisodeDownloadStart(&episode)
	if err == nil {
		t.Error("Expected error when notifying episode download start with MockErrorNotifier")
	}

	err = NotifyMovieDownloadStart(&movie)
	if err == nil {
		t.Error("Expected error when notifying movie download start with MockErrorNotifier")
	}
}
