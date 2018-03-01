package notifier

import (
	. "github.com/macarrie/flemzerd/objects"
)

type Notifier interface {
	Status() (Module, error)
	Send(title, content string) error
}
