package provider

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/macarrie/flemzerd/db"
	mock "github.com/macarrie/flemzerd/mocks"
	. "github.com/macarrie/flemzerd/objects"
	watchlist "github.com/macarrie/flemzerd/watchlists"
)

func init() {
	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
}

func TestAddProvider(t *testing.T) {
	providersLength := len(providersCollection)
	p := mock.TVProvider{}
	AddProvider(p)

	if len(providersCollection) != providersLength+1 {
		t.Error("Expected ", providersLength+1, " providers, got ", len(providersCollection))
	}
}

func TestStatus(t *testing.T) {
	pr1 := mock.TVProvider{}
	pr2 := mock.MovieProvider{}

	providersCollection = []Provider{pr1, pr2}

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 provider modules status, got %d instead", len(mods))
	}
	if err != nil {
		t.Error("Expected not to have error for provider status")
	}

	AddProvider(mock.ErrorProvider{})
	_, err = Status()
	if err == nil {
		t.Error("Expected to have aggregated error for provider status")
	}
}

func TestReset(t *testing.T) {
	p := mock.MovieProvider{}
	AddProvider(p)
	Reset()

	if len(providersCollection) != 0 {
		t.Error("Expected provider collection to be empty after reset")
	}
}

func TestFindShow(t *testing.T) {
	db.ResetDb()
	testShow := MediaIds{
		Title: "Test show",
	}

	providersCollection = []Provider{}
	_, err := FindShow(testShow)
	if err == nil {
		t.Error("Expected to have error when calling FindShow with no TV providers in collection")
	}

	AddProvider(mock.TVProvider{})

	show, err := FindShow(testShow)
	if err != nil {
		t.Error("Got error during FindShow: ", err)
	}
	if show.Model.ID != 1 {
		t.Errorf("Expected show with id 1, got id %v instead\n", show.Model.ID)
	}

	providersCollection = []Provider{mock.ErrorProvider{}}

	_, err = FindShow(testShow)
	if err == nil {
		t.Error("Expected to have error when calling FindShow with mock.ErrorProvider")
	}
}

func TestFindMovie(t *testing.T) {
	db.ResetDb()
	testMovie := MediaIds{
		Title: "Test Movie",
	}
	db.Client.Save(&testMovie)

	providersCollection = []Provider{}
	_, err := FindMovie(testMovie)
	if err == nil {
		t.Error("Expected to have error when calling FindMovie with no TV providers in collection")
	}

	p := mock.MovieProvider{}
	providersCollection = []Provider{p}

	movie, err := FindMovie(testMovie)
	if err != nil {
		t.Error("Got error during FindMovie: ", err)
	}
	if movie.Model.ID != 1 {
		t.Errorf("Expected movie with id 1, got id %v instead\n", movie.Model.ID)
	}

	//Calling find movie a second should get object from Db
	movie, err = FindMovie(testMovie)
	if err != nil {
		t.Error("Got error during FindMovie: ", err)
	}
	if movie.Model.ID != 1 {
		t.Errorf("Expected movie with id 1, got id %v instead\n", movie.Model.ID)
	}

	providersCollection = []Provider{mock.ErrorProvider{}}

	_, err = FindMovie(testMovie)
	if err == nil {
		t.Error("Expected to have error when calling FindMovie with mock.ErrorProvider")
	}
}

func TestFindRecentlyAiredEpisodesForShow(t *testing.T) {
	testShow := TvShow{
		Title:         "Test show",
		OriginalTitle: "Test show",
	}

	providersCollection = []Provider{}
	_, err := FindRecentlyAiredEpisodesForShow(testShow)
	if err == nil {
		t.Error("Expected to have error when calling FindRecentlyAiredEpisodesForShow with no TV providers in collection")
	}

	p := mock.TVProvider{}
	providersCollection = []Provider{p}

	episodeList, err := FindRecentlyAiredEpisodesForShow(testShow)
	if err != nil {
		t.Error("Got error during FindRecentlyAiredEpisodesForShow: ", err)
	}
	if episodeList[0].Model.ID != 1 {
		t.Errorf("Expected episode with id 1, got id %v instead\n", episodeList[0].Model.ID)
	}
}

func TestGetSeasonEpisodeList(t *testing.T) {
	testShow := TvShow{
		Title:         "Test show",
		OriginalTitle: "Test show",
	}

	providersCollection = []Provider{}
	_, err := GetSeasonEpisodeList(testShow, 1)
	if err == nil {
		t.Error("Expected to have error when calling GetSeasonEpisodeList with no TV providers in collection")
	}

	p := mock.TVProvider{}
	providersCollection = []Provider{p}

	episodeList, err := GetSeasonEpisodeList(testShow, 1)
	if err != nil {
		t.Error("Got error during GetSeasonEpisodeList: ", err)
	}
	if episodeList[0].Model.ID != 1 {
		t.Errorf("Expected episode with id 1, got id %v instead\n", episodeList[0].Model.ID)
	}

	// Calling the method a second times gets episodes from DB
	episodeList, err = GetSeasonEpisodeList(testShow, 1)
	if err != nil {
		t.Error("Got error during GetSeasonEpisodeList: ", err)
	}
	if episodeList[0].Model.ID != 1 {
		t.Errorf("Expected episode with id 1, got id %v instead\n", episodeList[0].Model.ID)
	}

	_, err = GetSeasonEpisodeList(testShow, 1000)
	if err == nil {
		t.Error("Expected to have error when asking for episodes of unknown season")
	}
}

func TestGetInfoFromConfig(t *testing.T) {
	db.ResetDb()
	pr1 := mock.TVProvider{}
	pr2 := mock.MovieProvider{}

	providersCollection = []Provider{pr1, pr2}

	GetTVShowsInfoFromConfig()
	GetMoviesInfoFromConfig()

	if len(TVShows) != 0 {
		t.Error("Expected to have no elements in tvshows from watchlists")
	}

	if len(Movies) != 0 {
		t.Error("Expected to have no elements in movies from watchlists")
	}

	db.ResetDb()

	w1 := mock.Watchlist{}
	w2 := mock.Watchlist{}
	watchlist.AddWatchlist(w1)
	watchlist.AddWatchlist(w2)

	providersCollection = []Provider{mock.ErrorProvider{}}

	GetTVShowsInfoFromConfig()
	GetMoviesInfoFromConfig()

	if len(TVShows) != 0 {
		t.Error("Expected to have no elements in tvshows from watchlists because only mock.ErrorProvider is defined")
	}

	if len(Movies) != 0 {
		t.Error("Expected to have no elements in movies from watchlists because only mock.ErrorProvider is defined")
	}

	providersCollection = []Provider{pr1, pr2}

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
		Model: gorm.Model{
			ID: 1000,
		},
		Title:         "test",
		OriginalTitle: "test",
	}

	testMovie := Movie{
		Model: gorm.Model{
			ID: 1000,
		},
		Title:         "test",
		OriginalTitle: "test",
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

	p1 := mock.TVProvider{}
	p2 := mock.MovieProvider{}

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

	if _, ok := (*tvProvider).(TVProvider); !ok {
		t.Error("TvProvider retrieved from getTvProvider is not of type TVProvider")
	}
	if _, ok := (*movieProvider).(MovieProvider); !ok {
		t.Error("MovieProvider retrieved from getMovieProvider is not of type MovieProvider")
	}
}

func TestMergeMediaIds(t *testing.T) {
	m1 := MediaIds{
		Title: "test",
		Tmdb:  1000,
	}

	m2 := MediaIds{
		Tvdb: 1000,
		Imdb: "test",
	}

	mergedIds := mergeMediaIds(m1, m2)

	if mergedIds.Title != "test" {
		t.Errorf("Expected merged IDs name to be 'test', got '%s' instead", mergedIds.Title)
	}
	if mergedIds.Tvdb != 1000 {
		t.Errorf("Expected merged IDs TVDB ID to be '1000', got '%d' instead", mergedIds.Tvdb)
	}
}

func TestGetSpecificProvider(t *testing.T) {
	p := mock.TVProvider{}
	providersCollection = []Provider{p}

	if _, err := GetProvider("Unknown"); err == nil {
		t.Error("Expected to have error when getting unknown provider, got none")
	}

	if _, err := GetProvider("TVProvider"); err != nil {
		t.Errorf("Got error while retrieving known provider: %s", err.Error())
	}
}
