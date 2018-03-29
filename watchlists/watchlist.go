package watchlist

import . "github.com/macarrie/flemzerd/objects"

type Watchlist interface {
	Status() (Module, error)
	GetTvShows() ([]MediaIds, error)
	GetMovies() ([]MediaIds, error)
}
