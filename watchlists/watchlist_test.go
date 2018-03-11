package watchlist

import . "github.com/macarrie/flemzerd/objects"

type MockWatchlist struct{}

func (p MockWatchlist) Status() (Module, error) {
	return Module{}, nil
}
