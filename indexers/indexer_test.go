package indexer

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type MockTVIndexer struct{}
type MockMovieIndexer struct{}
type MockOkTVIndexer struct{}

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

	if season == 0 {
		return []Torrent{}, fmt.Errorf("Mock error")
	}

	return []Torrent{
		Torrent{
			Name:    "Torrent1.s01.e01.720p",
			Link:    "torrent1.torrent",
			Seeders: 1,
		},
		Torrent{
			Name:    "Torrent2.s01.e01.720p",
			Link:    "torrent2.torrent",
			Seeders: 2,
		},
		Torrent{
			Name:    "Torrent3.s01.e01.720p",
			Link:    "torrent3.torrent",
			Seeders: 3,
		},
		Torrent{
			Name:    "Torrent4.s02.e02",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
	}, nil
}

func (m MockMovieIndexer) GetTorrentForMovie(movieName string) ([]Torrent, error) {
	if movieName == "" {
		return []Torrent{}, nil
	}

	if movieName == "error" {
		return []Torrent{}, fmt.Errorf("Mock error")
	}


	return []Torrent{
		Torrent{
			Name:    "Torrent1.720p.2018",
			Link:    "torrent1.torrent",
			Seeders: 1,
		},
		Torrent{
			Name:    "Torrent2.720p.2018",
			Link:    "torrent2.torrent",
			Seeders: 2,
		},
		Torrent{
			Name:    "Torrent3.720p.2018",
			Link:    "torrent3.torrent",
			Seeders: 3,
		},
		Torrent{
			Name:    "Torrent4",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
	}, nil
}

func (m MockOkTVIndexer) GetName() string {
	return "MockIndexer"
}
func (m MockOkTVIndexer) Status() (Module, error) {
	return Module{
		Name: "MockOkTVIndexer",
		Type: "indexer",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
