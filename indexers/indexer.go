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
	GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error)
}

type MovieIndexer interface {
	GetName() string
	Status() (Module, error)
	GetTorrentForMovie(movieName string) ([]Torrent, error)
}
