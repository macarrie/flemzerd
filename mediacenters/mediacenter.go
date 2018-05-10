package mediacenter

import . "github.com/macarrie/flemzerd/objects"

type MediaCenter interface {
	Status() (Module, error)
	GetName() string
	RefreshLibrary() error
}
