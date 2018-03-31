package watchlist

import (
	"testing"

	. "github.com/macarrie/flemzerd/objects"
)

func TestAddWatchlist(t *testing.T) {
	watchlistLength := len(watchlistsCollection)
	w := MockWatchlist{}
	AddWatchlist(w)

	if len(watchlistsCollection) != watchlistLength+1 {
		t.Error("Expected ", watchlistLength+1, " providers, got ", len(watchlistsCollection))
	}
}

func TestStatus(t *testing.T) {
	w1 := MockWatchlist{}
	w2 := MockWatchlist{}

	watchlistsCollection = []Watchlist{w1, w2}

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 watchlist modules status, got %d instead", len(mods))
	}
	if err == nil {
		t.Error("Expected to have aggregated error for watchlist status")
	}
}

func TestReset(t *testing.T) {
	w := MockWatchlist{}
	AddWatchlist(w)
	Reset()

	if len(watchlistsCollection) != 0 {
		t.Error("Expected watchlist collection to be empty after reset")
	}
}

func TestGetShows(t *testing.T) {
	w1 := MockWatchlist{}
	w2 := MockWatchlist{}

	watchlistsCollection = []Watchlist{w1, w2}

	shows, err := GetTvShows()
	if err != nil {
		t.Errorf("Got error when getting tvshows from watchlist: %s", err.Error())
	}

	if len(shows) != 1 {
		t.Errorf("Expected 1 show in watchlists, got %d instead", len(shows))
	}
}

func TestGetMovies(t *testing.T) {
	w1 := MockWatchlist{}
	w2 := MockWatchlist{}

	watchlistsCollection = []Watchlist{w1, w2}

	movies, err := GetMovies()
	if err != nil {
		t.Errorf("Got error when getting movies from watchlist: %s", err.Error())
	}

	if len(movies) != 1 {
		t.Errorf("Expected 1 movie in watchlists, got %d instead", len(movies))
	}
}

func TestRemoveDuplicates(t *testing.T) {
	uniqueList := removeDuplicates([]MediaIds{
		MediaIds{
			Name: "test1",
		},
		MediaIds{
			Name: "test2",
		},
		MediaIds{
			Name: "test1",
		},
	})

	if len(uniqueList) != 2 {
		t.Errorf("Expected to have 2 element in list after removing duplicates, got %d instead", len(uniqueList))
	}
}

func TestGetWatchlist(t *testing.T) {
	w1 := MockWatchlist{}
	watchlistsCollection = []Watchlist{w1}

	if _, err := GetWatchlist("Unknown"); err == nil {
		t.Error("Expected to have error when getting unknown watchlist, got none")
	}

	if _, err := GetWatchlist("MockWatchlist"); err != nil {
		t.Errorf("Got error while retrieving known watchlist: %s", err.Error())
	}
}
