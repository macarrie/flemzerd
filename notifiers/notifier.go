package notifier

import (
	. "github.com/macarrie/flemzerd/objects"
)

// Notifier is the generic interface that a struct has to implement in order to be used as a notifier in flemzerd
type Notifier interface {
	Status() (Module, error)
	GetName() string
	Send(title, content string) error
}
