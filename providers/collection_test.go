package provider

import (
	"testing"

	. "github.com/macarrie/flemzerd/objects"
	watchlist "github.com/macarrie/flemzerd/watchlists"
)

func TestAddProvider(t *testing.T) {
	providersLength := len(providersCollection)
	p := MockTVProvider{}
	AddProvider(p)

	if len(providersCollection) != providersLength+1 {
		t.Error("Expected ", providersLength+1, " providers, got ", len(providersCollection))
	}
}

func TestStatus(t *testing.T) {
	pr1 := MockTVProvider{}
	pr2 := MockMovieProvider{}

	providersCollection = []Provider{pr1, pr2}

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 provider modules status, got %d instead", len(mods))
	}
	if err == nil {
		t.Error("Expected to have aggregated error for provider status")
	}
}

func TestFindShow(t *testing.T) {
	p := MockTVProvider{}
	providersCollection = []Provider{p}

	show, err := FindShow("Test show")
	if err != nil {
		t.Error("Got error during FindShow: ", err)
	}
	if show.Id != 1000 {
		t.Errorf("Expected show with id 1000, got id %v instead\n", show.Id)
	}
}

func TestFindMovie(t *testing.T) {
	p := MockMovieProvider{}
	providersCollection = []Provider{p}

	movie, err := FindMovie("Test Movie")
	if err != nil {
		t.Error("Got error during FindShow: ", err)
	}
	if movie.Id != 1000 {
		t.Errorf("Expected movie with id 1000, got id %v instead\n", movie.Id)
	}
}

func TestFindRecentlyAiredEpisodesForShow(t *testing.T) {
	p := MockTVProvider{}
	providersCollection = []Provider{p}

	episodeList, err := FindRecentlyAiredEpisodesForShow(TvShow{
		Name: "Test show",
	})
	if err != nil {
		t.Error("Got error during FindRecentlyAiredEpisodesForShow: ", err)
	}
	if episodeList[0].Id != 1000 {
		t.Errorf("Expected episode with id 1000, got id %v instead\n", episodeList[0].Id)
	}
}

func TestGetInfoFromConfig(t *testing.T) {
	pr1 := MockTVProvider{}
	pr2 := MockMovieProvider{}

	providersCollection = []Provider{pr1, pr2}

	GetTVShowsInfoFromConfig()
	GetMoviesInfoFromConfig()

	if len(TVShows) != 0 {
		t.Error("Expected to have no elements in tvshows from watchlists")
	}

	if len(Movies) != 0 {
		t.Error("Expected to have no elements in movies from watchlists")
	}

	w1 := MockWatchlist{}
	w2 := MockWatchlist{}
	watchlist.AddWatchlist(w1)
	watchlist.AddWatchlist(w2)

	GetTVShowsInfoFromConfig()
	GetMoviesInfoFromConfig()

	if len(TVShows) != 1 {
		t.Errorf("Expected to have 1 elements in tvshows from watchlists, got %d instead", len(TVShows))
	}

	if len(Movies) != 1 {
		t.Errorf("Expected to have 1 elements in movies from watchlists, got %d instead", len(Movies))
	}
}

func TestRemoveDuplicate(t *testing.T) {
	testShow := TvShow{
		Id:   1000,
		Name: "test",
	}

	testMovie := Movie{
		Id:    1000,
		Title: "test",
	}

	uniqueTVList := removeDuplicateShows([]TvShow{testShow, testShow})
	uniqueMovieList := removeDuplicateMovies([]Movie{testMovie, testMovie})

	if len(uniqueTVList) != 1 {
		t.Errorf("Expected to have 1 show in list after removing duplicates, got %d instead", len(uniqueTVList))
	}
	if len(uniqueMovieList) != 1 {
		t.Errorf("Expected to have 1 movie in list after removing duplicates, got %d instead", len(uniqueMovieList))
	}
}

func TestGetProvider(t *testing.T) {
	providersCollection = []Provider{}

	p1 := MockTVProvider{}
	p2 := MockMovieProvider{}

	if tvProvider := getTVProvider(); tvProvider != nil {
		t.Error("Expected to not being able to retrieve tv provider")
	}
	if movieProvider := getMovieProvider(); movieProvider != nil {
		t.Error("Expected to not being able to retrieve movie provider")
	}

	AddProvider(p1)
	AddProvider(p2)

	tvProvider := getTVProvider()
	movieProvider := getMovieProvider()

	if tvProvider == nil {
		t.Error("Expected to be able to retrieve tv provider")
	}
	if movieProvider == nil {
		t.Error("Expected to be able to retrieve movie provider")
	}

	if _, ok := tvProvider.(TVProvider); !ok {
		t.Error("TvProvider retrieved from getTvProvider is not of type TVProvider")
	}
	if _, ok := movieProvider.(MovieProvider); !ok {
		t.Error("MovieProvider retrieved from getMovieProvider is not of type MovieProvider")
	}
}
