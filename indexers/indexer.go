package indexer

import (
	. "github.com/macarrie/flemzerd/objects"
)

type Indexer interface {
	GetName() string
	Status() (Module, error)
}

type TVIndexer interface {
	GetName() string
	Status() (Module, error)
	GetTorrentForEpisode(episode Episode) ([]Torrent, error)
}

type MovieIndexer interface {
	GetName() string
	Status() (Module, error)
	GetTorrentForMovie(movieName string) ([]Torrent, error)
}
