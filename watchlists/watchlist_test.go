package watchlist

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

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

func (w MockWatchlist) GetTvShows() ([]string, error) {
	return []string{"test show"}, nil
}

func (w MockWatchlist) GetMovies() ([]string, error) {
	return []string{"test movie"}, nil
}
