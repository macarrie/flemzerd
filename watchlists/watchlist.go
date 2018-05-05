package watchlist

import . "github.com/macarrie/flemzerd/objects"

type Watchlist interface {
	Status() (Module, error)
	GetName() string
	GetTvShows() ([]MediaIds, error)
	GetMovies() ([]MediaIds, error)
}
