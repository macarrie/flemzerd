package provider

import (
	"fmt"

	"github.com/jinzhu/gorm"
	. "github.com/macarrie/flemzerd/objects"
)

type MockTVProvider struct{}
type MockMovieProvider struct{}

func (p MockTVProvider) Status() (Module, error) {
	var err error = fmt.Errorf("Provider error")
	return Module{
		Name: "MockTVProvider",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}
func (p MockMovieProvider) Status() (Module, error) {
	var err error = fmt.Errorf("Provider error")
	return Module{
		Name: "MockTVProvider",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (p MockTVProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	return TvShow{
		Model: gorm.Model{
			ID: 1000,
		},
		Name: "Test show",
	}, nil
}
func (p MockTVProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockTVProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockTVProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{
		Episode{
			Number: 1,
			Season: 1,
			Name:   "Test episode",
			Model: gorm.Model{
				ID: 1000,
			},
		},
	}, nil
}
func (p MockMovieProvider) GetMovie(movie MediaIds) (Movie, error) {
	return Movie{
		Model: gorm.Model{
			ID: 1000,
		},
		Title: "Test Movie",
	}, nil
}

/////////////////////////////////////////////////////
// Define MockWatchlist to use with provider tests //
/////////////////////////////////////////////////////

type MockWatchlist struct{}

func (w MockWatchlist) Status() (Module, error) {
	var err error = fmt.Errorf("Watchlist error")
	return Module{
		Name: "MockWatchlist",
		Type: "watchlist",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (w MockWatchlist) GetTvShows() ([]MediaIds, error) {
	return []MediaIds{
		MediaIds{
			Name: "test show",
		},
	}, nil
}

func (w MockWatchlist) GetMovies() ([]MediaIds, error) {
	return []MediaIds{
		MediaIds{
			Name: "test movie",
		},
	}, nil
}
