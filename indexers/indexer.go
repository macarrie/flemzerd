package indexer

import (
	. "github.com/macarrie/flemzerd/objects"
)

type Indexer interface {
	GetName() string
	IsAlive() error
	GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error)
}
