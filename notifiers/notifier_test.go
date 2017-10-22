package notifier

import (
	"testing"
)

type MockNotifier struct{}

func (n MockNotifier) Setup(authCredentials map[string]string) {
	return
}

func (n MockNotifier) Init() error {
	return nil
}

func (n MockNotifier) Send(title, content string) error {
	return nil
}

func TestAddNotifier(t *testing.T)         {}
func TestRemoveFromRetention(t *testing.T) {}
func TestSendNotification(t *testing.T)    {}
func TestNotifyRecentEpisode(t *testing.T) {}
