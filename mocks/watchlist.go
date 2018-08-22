package mock

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type Watchlist struct{}
type ErrorWatchlist struct{}

func (w Watchlist) Status() (Module, error) {
	return Module{
		Name: "Watchlist",
		Type: "watchlist",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (w ErrorWatchlist) Status() (Module, error) {
	var err error = fmt.Errorf("Watchlist error")
	return Module{
		Name: "ErrorWatchlist",
		Type: "watchlist",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (w Watchlist) GetName() string {
	return "Watchlist"
}
func (w ErrorWatchlist) GetName() string {
	return "ErrorWatchlist"
}

func (w Watchlist) GetTvShows() ([]MediaIds, error) {
	return []MediaIds{
		MediaIds{
			Name: "test show",
		},
	}, nil
}

func (w ErrorWatchlist) GetTvShows() ([]MediaIds, error) {
	return []MediaIds{}, fmt.Errorf("Error while getting TV Shows from watchlist")
}

func (w Watchlist) GetMovies() ([]MediaIds, error) {
	return []MediaIds{
		MediaIds{
			Name: "test movie",
		},
	}, nil
}

func (w ErrorWatchlist) GetMovies() ([]MediaIds, error) {
	return []MediaIds{}, fmt.Errorf("Error while getting movies from watchlist")
}
