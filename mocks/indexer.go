package mock

import (
	"fmt"
	"github.com/macarrie/flemzerd/downloadable"

	. "github.com/macarrie/flemzerd/objects"
)

type TVIndexer struct{}
type MovieIndexer struct{}
type ErrorTVIndexer struct{}
type ErrorMovieIndexer struct{}

func (m TVIndexer) GetName() string {
	return "TVIndexer"
}
func (m MovieIndexer) GetName() string {
	return "MovieIndexer"
}
func (m ErrorTVIndexer) GetName() string {
	return "ErrorTVIndexer"
}
func (m ErrorMovieIndexer) GetName() string {
	return "ErrorMovieIndexer"
}

func (m TVIndexer) Status() (Module, error) {
	return Module{
		Name: "TVIndexer",
		Type: "indexer",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (m MovieIndexer) Status() (Module, error) {
	return Module{
		Name: "MovieIndexer",
		Type: "indexer",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (m ErrorTVIndexer) Status() (Module, error) {
	var err error = fmt.Errorf("Indexer error")
	return Module{
		Name: "TVIndexer",
		Type: "indexer",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}
func (m ErrorMovieIndexer) Status() (Module, error) {
	var err error = fmt.Errorf("Indexer error")
	return Module{
		Name: "MovieIndexer",
		Type: "indexer",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func getTorrentForEpisode(episode Episode) ([]Torrent, error) {
	if episode.Number == 0 {
		return []Torrent{}, nil
	}

	if episode.Season == 0 {
		return []Torrent{}, fmt.Errorf(" error")
	}

	return []Torrent{
		Torrent{
			Name:    "Torrent1.s01e01.720p",
			Link:    "torrent1.torrent",
			Seeders: 1,
		},
		Torrent{
			Name:    "Torrent2.s01e01.1080p",
			Link:    "torrent2.torrent",
			Seeders: 2,
		},
		Torrent{
			Name:    "Torrent3.s01e01.720p",
			Link:    "torrent3.torrent",
			Seeders: 3,
		},
		Torrent{
			Name:    "Torrent4.s01e01.cam",
			Link:    "torrent4.cam.torrent",
			Seeders: 4,
		},
		Torrent{
			Name:    "Torrent4.s02e02.cam",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
		Torrent{
			Name:    "Torrent4.s02e02.480p",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
		Torrent{
			Name:    "Torrent4.s02e02",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
	}, nil
}
func (m TVIndexer) GetTorrentForEpisode(episode Episode) ([]Torrent, error) {
	return getTorrentForEpisode(episode)
}
func (m ErrorTVIndexer) GetTorrentForEpisode(episode Episode) ([]Torrent, error) {
	return getTorrentForEpisode(episode)
}

func getTorrentForMovie(movieName string) ([]Torrent, error) {
	if movieName == "" {
		return []Torrent{}, nil
	}

	if movieName == "error" {
		return []Torrent{}, fmt.Errorf(" error")
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
			Name:    "Torrent4.cam",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
		Torrent{
			Name:    "Torrent5.1994",
			Link:    "torrent5.torrent",
			Seeders: 4,
		},
		Torrent{
			Name:    "Torrent6.2018.screener",
			Link:    "torrent6.torrent",
			Seeders: 4,
		},
		Torrent{
			Name:    "Torrent7.480p.screener",
			Link:    "torrent7.torrent",
			Seeders: 2,
		},
	}, nil
}
func (m MovieIndexer) GetTorrent(d downloadable.Downloadable) ([]Torrent, error) {
	return getTorrentForMovie(d)
}
func (m ErrorMovieIndexer) GetTorrent(d downloadable.Downloadable) ([]Torrent, error) {
	return getTorrentForMovie(d)
}
