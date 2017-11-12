package indexer

import (
	. "github.com/macarrie/flemzerd/objects"
)

type Indexer interface {
	GetName() string
	GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error)
}
