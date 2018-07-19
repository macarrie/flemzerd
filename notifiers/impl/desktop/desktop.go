package desktop

import (
	"github.com/0xAX/notificator"
	notifier_helper "github.com/macarrie/flemzerd/helpers/notifiers"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

type DesktopNotifier struct {
	Notifier notificator.Notificator
}

func (d *DesktopNotifier) Init() error {
	d.Notifier = *notificator.New(notificator.Options{
		AppName: "flemzerd",
	})

	return nil
}

func (d *DesktopNotifier) Status() (Module, error) {
	// TODO
	log.Debug("Checking desktop notifier status")

	return Module{
		Name: d.GetName(),
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (d *DesktopNotifier) GetName() string {
	return "Desktop notifier"
}

func (d *DesktopNotifier) Send(notif Notification) error {
	log.Debug("Sending Desktop notification")

	title, content, err := notifier_helper.GetNotificationText(notif)
	if err != nil {
		return err
	}

	d.Notifier.Push(title, content, "", notificator.UR_NORMAL)

	return nil
}

func New() *DesktopNotifier {
	var notifier DesktopNotifier
	notifier.Init()

	return &notifier
}
