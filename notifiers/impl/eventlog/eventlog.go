package eventlog

import (
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

type EventLogNotifier struct{}

func (d *EventLogNotifier) Status() (Module, error) {
	log.Debug("Checking event log notifier status")

	return Module{
		Name: d.GetName(),
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (d *EventLogNotifier) GetName() string {
	return "Event Log Notifier"
}

func (d *EventLogNotifier) Send(notif Notification) error {
	log.Debug("Sending notification to event log")

	db.Client.Create(&notif)

	return nil
}

func New() *EventLogNotifier {
	return &EventLogNotifier{}
}
