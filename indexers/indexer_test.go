package indexer

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type MockTVIndexer struct{}
type MockMovieIndexer struct{}

func (m MockTVIndexer) GetName() string {
	return "MockIndexer"
}
func (m MockMovieIndexer) GetName() string {
	return "MockIndexer"
}

func (m MockTVIndexer) Status() (Module, error) {
	var err error = fmt.Errorf("Indexer error")
	return Module{
		Name: "MockTVIndexer",
		Type: "indexer",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (m MockMovieIndexer) Status() (Module, error) {
	var err error = fmt.Errorf("Indexer error")
	return Module{
		Name: "MockMovieIndexer",
		Type: "indexer",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (m MockTVIndexer) GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	if episode == 0 {
		return []Torrent{}, nil
	}

	return []Torrent{
		Torrent{
			Name:    "Torrent1.s01.e01",
			Link:    "torrent1.torrent",
			Seeders: 1,
		},
		Torrent{
			Name:    "Torrent2.s01.e01",
			Link:    "torrent2.torrent",
			Seeders: 2,
		},
		Torrent{
			Name:    "Torrent3.s01.e01",
			Link:    "torrent3.torrent",
			Seeders: 3,
		},
	}, nil
}

func (m MockMovieIndexer) GetTorrentForMovie(movieName string) ([]Torrent, error) {
	if movieName == "" {
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
