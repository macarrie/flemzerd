package indexer

import (
	"github.com/macarrie/flemzerd/downloadable"
	. "github.com/macarrie/flemzerd/objects"
)

type Indexer interface {
	GetName() string
	Status() (Module, error)
	CheckCapabilities(d downloadable.Downloadable) bool
	GetTorrents(d downloadable.Downloadable) ([]Torrent, error)
}
