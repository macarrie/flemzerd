package desktop

import (
	"fmt"

	"github.com/0xAX/notificator"
	log "github.com/macarrie/flemzerd/logging"
)

type DesktopNotifier struct {
	Name     string
	Notifier notificator.Notificator
}

func (d *DesktopNotifier) Init() error {
	d.Name = "Desktop notifier"
	d.Notifier = *notificator.New(notificator.Options{
		AppName: "flemzerd",
	})

	return nil
}

func (d *DesktopNotifier) Send(title, content string) error {
	log.WithFields(log.Fields{
		"title":   title,
		"content": content,
	}).Debug("Sending Desktop notification")

	fmt.Printf("%+v\n", d)
	d.Notifier.Push(title, content, "", notificator.UR_NORMAL)

	return nil
}

func New() *DesktopNotifier {
	var notifier DesktopNotifier
	notifier.Init()

	return &notifier
}
