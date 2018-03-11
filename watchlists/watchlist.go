package watchlist

import . "github.com/macarrie/flemzerd/objects"

type Watchlist interface {
	Status() (Module, error)
	GetTvShows() ([]string, error)
	GetMovies() ([]string, error)
}
