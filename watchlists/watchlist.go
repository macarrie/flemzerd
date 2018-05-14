package watchlist

import . "github.com/macarrie/flemzerd/objects"

// Watchlist is the generic interface that a struct has to implement in order to be used as a watchlist in flemzerd
type Watchlist interface {
	Status() (Module, error)
	GetName() string
	GetTvShows() ([]MediaIds, error)
	GetMovies() ([]MediaIds, error)
}
