package watchlist

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type MockWatchlist struct{}
type MockErrorWatchlist struct{}

func (w MockWatchlist) Status() (Module, error) {
	return Module{
		Name: "MockWatchlist",
		Type: "watchlist",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (w MockErrorWatchlist) Status() (Module, error) {
	var err error = fmt.Errorf("Watchlist error")
	return Module{
		Name: "MockErrorWatchlist",
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

func (w MockErrorWatchlist) GetTvShows() ([]MediaIds, error) {
	return []MediaIds{}, fmt.Errorf("Error while getting TV Shows from watchlist")
}

func (w MockWatchlist) GetMovies() ([]MediaIds, error) {
	return []MediaIds{
		MediaIds{
			Name: "test movie",
		},
	}, nil
}

func (w MockErrorWatchlist) GetMovies() ([]MediaIds, error) {
	return []MediaIds{}, fmt.Errorf("Error while getting movies from watchlist")
}
