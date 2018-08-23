package mock

import (
	"fmt"

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

func getTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	if episode == 0 {
		return []Torrent{}, nil
	}

	if season == 0 {
		return []Torrent{}, fmt.Errorf(" error")
	}

	return []Torrent{
		Torrent{
			Name:    "Torrent1.s01e01.720p",
			Link:    "torrent1.torrent",
			Seeders: 1,
		},
		Torrent{
			Name:    "Torrent2.s01e01.720p",
			Link:    "torrent2.torrent",
			Seeders: 2,
		},
		Torrent{
			Name:    "Torrent3.s01e01.720p",
			Link:    "torrent3.torrent",
			Seeders: 3,
		},
		Torrent{
			Name:    "Torrent4.s02e02",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
	}, nil
}
func (m TVIndexer) GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	return getTorrentForEpisode(show, season, episode)
}
func (m ErrorTVIndexer) GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	return getTorrentForEpisode(show, season, episode)
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
			Name:    "Torrent4",
			Link:    "torrent4.torrent",
			Seeders: 4,
		},
	}, nil
}
func (m MovieIndexer) GetTorrentForMovie(movieName string) ([]Torrent, error) {
	return getTorrentForMovie(movieName)
}
func (m ErrorMovieIndexer) GetTorrentForMovie(movieName string) ([]Torrent, error) {
	return getTorrentForMovie(movieName)
}
