package notifier

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	mock "github.com/macarrie/flemzerd/mocks"
	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	// go test makes a cd into package directory when testing. We must go up by one level to load our testdata
	configuration.UseFile("../testdata/test_config.toml")
	err := configuration.Load()
	configuration.Config.Notifications.Enabled = true

	if err != nil {
		fmt.Print("Could not load test configuration file: ", err)
	}

	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
}

func TestStatus(t *testing.T) {
	notifiersCollection = []Notifier{mock.Notifier{}}
	_, err := Status()
	if err != nil {
		t.Error("Expected to have no error for notifier status")
	}

	notifiersCollection = []Notifier{mock.Notifier{}, mock.ErrorNotifier{}}

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 notifier modules status, got %d instead", len(mods))
	}
	if err == nil {
		t.Error("Expected to have aggregated error for notifier status")
	}
}

func TestReset(t *testing.T) {
	n := mock.Notifier{}
	AddNotifier(n)
	Reset()

	if len(notifiersCollection) != 0 {
		t.Error("Expected notifier collection to be empty after reset")
	}
}

func TestAddNotifier(t *testing.T) {
	notifiersLength := len(notifiersCollection)
	m := mock.Notifier{}
	AddNotifier(m)

	if len(notifiersCollection) != notifiersLength+1 {
		t.Error("Expected ", notifiersLength+1, " Notifiers, got ", len(notifiersCollection))
	}
}

func TestSendNotification(t *testing.T) {
	n := 2

	notif := Notification{
		Type: NOTIFICATION_NEW_EPISODE,
		Movie: Movie{
			Title:         "Test Movie",
			OriginalTitle: "Test Movie",
		},
	}
	notifiersCollection = []Notifier{}
	mockNotifiers := make([]mock.Notifier, n)
	n1 := mockNotifiers[0]
	count := n1.GetNotificationCount()

	for i := range mockNotifiers {
		mockNotifiers[i] = mock.Notifier{}
		AddNotifier(mockNotifiers[i])
	}

	SendNotification(notif)

	if n1.GetNotificationCount() != count+n {
		t.Error("Expected to send ", n, " notifications, but ", n1.GetNotificationCount(), " notifications have been sent")
	}

	// If some notifications have been sent, do not return an error
	AddNotifier(mock.ErrorNotifier{})
	err := SendNotification(notif)
	if err != nil {
		t.Error("Expected to have no error when sending notifications")
	}

	notifiersCollection = []Notifier{mock.ErrorNotifier{}}
	if err := SendNotification(notif); err == nil {
		t.Error("Expected to have an error when sending notifications")
	}

	prev := n1.GetNotificationCount()
	configuration.Config.Notifications.Enabled = false
	SendNotification(notif)
	if n1.GetNotificationCount() != prev {
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
		Title:         "Test TVShow",
		OriginalTitle: "Test TVShow",
	}

	episode := Episode{
		Model: gorm.Model{
			ID: 30,
		},
		TvShow: show,
		Season: 3,
		Number: 4,
		Title:  "Test Episode S03E04",
		Date:   time.Now(),
	}

	m := mock.Notifier{}
	AddNotifier(m)

	count := m.GetNotificationCount()

	NotifyRecentEpisode(&episode)
	if m.GetNotificationCount() != count+1 {
		t.Error("Expected recent notification to be sent, got none")
	}

	NotifyDownloadedEpisode(&episode)
	if m.GetNotificationCount() != count+2 {
		t.Error("Expected downloaded notification to be sent, got none")
	}

	NotifyFailedEpisode(&episode)
	if m.GetNotificationCount() != count+3 {
		t.Error("Expected failed notification to be sent, got none")
	}

	NotifyRecentEpisode(&episode)
	if m.GetNotificationCount() != count+3 {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}

	configuration.Config.Notifications.Enabled = false

	NotifyRecentEpisode(&episode)
	NotifyDownloadedEpisode(&episode)
	NotifyFailedEpisode(&episode)
	NotifyRecentEpisode(&episode)

	if m.GetNotificationCount() != count+3 {
		t.Error("Expected notification not to be sent because notifications are disabled, but notifications have been sent anyway")
	}

	configuration.Config.Notifications.Enabled = true

	notifiersCollection = []Notifier{mock.ErrorNotifier{}}

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
		Title:         "Test movie",
		OriginalTitle: "Test movie",
	}

	m := mock.Notifier{}
	AddNotifier(m)

	count := m.GetNotificationCount()

	NotifyNewMovie(&movie)
	if m.GetNotificationCount() != count+1 {
		t.Error("Expected movie notification to be sent, got none")
	}

	NotifyDownloadedMovie(&movie)
	if m.GetNotificationCount() != count+2 {
		t.Error("Expected downloaded movie notification to be sent, got none")
	}

	NotifyFailedMovie(&movie)
	if m.GetNotificationCount() != count+3 {
		t.Error("Expected failed notification to be sent, got none")
	}

	NotifyNewMovie(&movie)
	if m.GetNotificationCount() != count+3 {
		t.Error("Expected notification not to be sent because episode is on retention, but a notification has been sent anyway")
	}

	configuration.Config.Notifications.Enabled = false

	NotifyNewMovie(&movie)
	NotifyDownloadedMovie(&movie)
	NotifyFailedMovie(&movie)
	NotifyNewMovie(&movie)

	if m.GetNotificationCount() != count+3 {
		t.Error("Expected notification not to be sent because notifications are disabled, but notifications have been sent anyway")
	}

	notifiersCollection = []Notifier{mock.ErrorNotifier{}}
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
	n := mock.Notifier{}
	AddNotifier(n)

	show := TvShow{
		Model: gorm.Model{
			ID: 1,
		},
		Title:         "Test TVShow",
		OriginalTitle: "Test TVShow",
	}

	episode := Episode{
		Model: gorm.Model{
			ID: 30,
		},
		TvShow: show,
		Season: 3,
		Number: 4,
		Title:  "Test Episode S03E04",
		Date:   time.Now(),
	}

	movie := Movie{
		Model: gorm.Model{
			ID: 1,
		},
		Title:         "Test movie",
		OriginalTitle: "Test movie",
	}

	count := n.GetNotificationCount()
	NotifyEpisodeDownloadStart(&episode)
	if n.GetNotificationCount() != count+1 {
		t.Error("Expected notification to be sent when notifying episode download start")
	}

	count = n.GetNotificationCount()
	NotifyMovieDownloadStart(&movie)
	if n.GetNotificationCount() != count+1 {
		t.Error("Expected notification to be sent when notifying movie download start")
	}

	configuration.Config.Notifications.NotifyDownloadStart = false
	count = n.GetNotificationCount()
	NotifyEpisodeDownloadStart(&episode)
	if n.GetNotificationCount() != count {
		t.Error("Expected notification not to be sent when notifying episode download start because of configuration params")
	}

	count = n.GetNotificationCount()
	NotifyMovieDownloadStart(&movie)
	if n.GetNotificationCount() != count {
		t.Error("Expected notification not to be sent when notifying movie download start because of configuration params")
	}

	configuration.Config.Notifications.NotifyDownloadStart = true
	notifiersCollection = []Notifier{mock.ErrorNotifier{}}
	err := NotifyEpisodeDownloadStart(&episode)
	if err == nil {
		t.Error("Expected error when notifying episode download start with mock.ErrorNotifier")
	}

	err = NotifyMovieDownloadStart(&movie)
	if err == nil {
		t.Error("Expected error when notifying movie download start with mock.ErrorNotifier")
	}
}

func TestGetNotifier(t *testing.T) {
	notifiersCollection = []Notifier{mock.Notifier{}}

	if _, err := GetNotifier("Unknown"); err == nil {
		t.Error("Expected to have error when getting unknown notifier, got none")
	}

	if _, err := GetNotifier("Notifier"); err != nil {
		t.Errorf("Got error while retrieving known notifier: %s", err.Error())
	}
}
