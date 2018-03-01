package indexer

import (
	. "github.com/macarrie/flemzerd/objects"
)

type MockIndexer struct{}

func (m MockIndexer) GetName() string {
	return "MockIndexer"
}

func (m MockIndexer) Status() (Module, error) {
	return Module{}, nil
}

func (m MockIndexer) GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	if episode == 0 {
		return []Torrent{}, nil
	}

	return []Torrent{
		Torrent{
			Name:    "Torrent1",
			Link:    "torrent1.torrent",
			Seeders: 1,
		},
		Torrent{
			Name:    "Torrent2",
			Link:    "torrent2.torrent",
			Seeders: 2,
		},
		Torrent{
			Name:    "Torrent3",
			Link:    "torrent3.torrent",
			Seeders: 3,
		},
	}, nil
}
