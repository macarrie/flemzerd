package retention

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	resetRetention()
	retentionPath = "/tmp/flemzerd_retention_test.dat"
}

func resetRetention() {
	retentionData.NotifiedEpisodes = make(map[int]Episode)
	retentionData.DownloadedEpisodes = make(map[int]Episode)
	retentionData.DownloadingEpisodes = make(map[int]*DownloadingEpisode)
	retentionData.FailedEpisodes = make(map[int]Episode)

	retentionData.NotifiedMovies = make(map[int]Movie)
	retentionData.DownloadedMovies = make(map[int]Movie)
	retentionData.DownloadingMovies = make(map[int]*DownloadingMovie)
	retentionData.FailedMovies = make(map[int]Movie)

	retentionData.Watchlists = make(map[string][]MediaIds)
}

func TestSave(t *testing.T) {
	resetRetention()

	retentionSave := retentionData
	err := Save()
	if err != nil {
		t.Errorf("An error occured during retention save: %s", err.Error())
	}
	Load()

	fmt.Printf("Save: %+v\n", retentionSave)
	fmt.Printf("Load: %+v\n", retentionData)
	if !cmp.Equal(retentionData, retentionSave) {
		t.Error("Retention data loaded from file is not the same as retention data saved beforehand")
	}
}

func TestElementInRetention(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id: 1000,
	}
	testMovie := Movie{
		Id: 1000,
	}

	if EpisodeHasBeenNotified(testEpisode) {
		t.Error("Expected episode not to be in notified episodes retention")
	}
	if MovieHasBeenNotified(testMovie) {
		t.Error("Expected movie not to be in notified movies retention")
	}
	if EpisodeHasBeenDownloaded(testEpisode) {
		t.Error("Expected episode not to be in downloaded episodes retention")
	}
	if MovieHasBeenDownloaded(testMovie) {
		t.Error("Expected movie not to be in downloaded movies retention")
	}
	if EpisodeIsInDownloadProcess(testEpisode) {
		t.Error("Expected episode not to be in downloading episodes retention")
	}
	if MovieIsInDownloadProcess(testMovie) {
		t.Error("Expected movie not to be in downloading movies retention")
	}
	if EpisodeHasDownloadFailed(testEpisode) {
		t.Error("Expected episode not to be in failed episode torrents retention")
	}
	if MovieHasDownloadFailed(testMovie) {
		t.Error("Expected movie not to be in failed movies torrents retention")
	}

	retentionData = RetentionData{
		NotifiedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
		DownloadedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
		DownloadingEpisodes: map[int]*DownloadingEpisode{
			testEpisode.Id: &DownloadingEpisode{
				Episode: testEpisode,
			},
		},
		FailedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
		NotifiedMovies: map[int]Movie{
			testMovie.Id: testMovie,
		},
		DownloadedMovies: map[int]Movie{
			testMovie.Id: testMovie,
		},
		DownloadingMovies: map[int]*DownloadingMovie{
			testMovie.Id: &DownloadingMovie{
				Movie: testMovie,
			},
		},
		FailedMovies: map[int]Movie{
			testMovie.Id: testMovie,
		},
	}

	if !EpisodeHasBeenNotified(testEpisode) {
		t.Error("Expected test episode to be present in notified episodes retention")
	}
	if !MovieHasBeenNotified(testMovie) {
		t.Error("Expected test movie to be present in notified movie retention")
	}
	if !EpisodeHasBeenDownloaded(testEpisode) {
		t.Error("Expected test episode to be present in downloaded episodes retention")
	}
	if !MovieHasBeenDownloaded(testMovie) {
		t.Error("Expected test movie to be present in downloaded movie retention")
	}
	if !EpisodeIsInDownloadProcess(testEpisode) {
		t.Error("Expected test episode to be present in downloading episodes retention")
	}
	if !MovieIsInDownloadProcess(testMovie) {
		t.Error("Expected test movie to be present in downloading movie retention")
	}
	if !EpisodeHasDownloadFailed(testEpisode) {
		t.Error("Expected test episode to be present in failed episode torrents retention")
	}
	if !MovieHasDownloadFailed(testMovie) {
		t.Error("Expected test movie to be present in failed movies torrents retention")
	}
}

func TestAddElementInRetention(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id: 1000,
	}
	testMovie := Movie{
		Id: 1000,
	}

	AddNotifiedEpisode(testEpisode)
	AddNotifiedMovie(testMovie)
	AddDownloadedEpisode(testEpisode)
	AddDownloadedMovie(testMovie)
	AddDownloadingEpisode(testEpisode)
	AddDownloadingMovie(testMovie)
	fmt.Printf("%+v\n", retentionData)

	if len(retentionData.NotifiedEpisodes) != 1 {
		t.Error("Expected to have 1 item in notified episodes retention, got ", len(retentionData.NotifiedEpisodes), " items instead")
	}
	if len(retentionData.NotifiedMovies) != 1 {
		t.Error("Expected to have 1 item in notified movies retention, got ", len(retentionData.NotifiedMovies), " items instead")
	}
	if len(retentionData.DownloadedEpisodes) != 1 {
		t.Error("Expected to have 1 item in downloaded episodes retention, got ", len(retentionData.DownloadedEpisodes), " items instead")
	}
	if len(retentionData.DownloadedMovies) != 1 {
		t.Error("Expected to have 1 item in downloaded movies retention, got ", len(retentionData.DownloadedMovies), " items instead")
	}
	if len(retentionData.DownloadingEpisodes) != 1 {
		t.Error("Expected to have 1 item in downloading episodes retention, got ", len(retentionData.DownloadingEpisodes), " items instead")
	}
	if len(retentionData.DownloadingMovies) != 1 {
		t.Error("Expected to have 1 item in downloading movies retention, got ", len(retentionData.DownloadingMovies), " items instead")
	}
}

func TestCleanOldNotifiedElements(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id:   1000,
		Date: time.Time{},
	}
	testMovie := Movie{
		Id:   1000,
		Date: time.Time{},
	}

	AddNotifiedEpisode(testEpisode)
	AddNotifiedMovie(testMovie)

	CleanOldNotifiedEpisodes()
	CleanOldNotifiedMovies()

	if len(retentionData.NotifiedEpisodes) > 0 {
		t.Error("Expected old episode to have been removed from notified episodes retention but episode is still present in retention")
	}
	if len(retentionData.NotifiedMovies) > 0 {
		t.Error("Expected old movie to have been removed from notified movies retention but movie is still present in retention")
	}
}

func TestRemoveElementFromRetention(t *testing.T) {
	resetRetention()
	e1 := Episode{
		Id: 1,
	}
	e2 := Episode{
		Id: 2,
	}
	m1 := Movie{
		Id: 1,
	}
	m2 := Movie{
		Id: 2,
	}

	AddNotifiedEpisode(e1)
	AddNotifiedEpisode(e2)
	AddDownloadedEpisode(e1)
	AddDownloadedEpisode(e2)
	AddDownloadingEpisode(e1)
	AddDownloadingEpisode(e2)
	AddFailedEpisode(e1)
	AddFailedEpisode(e2)

	AddNotifiedMovie(m1)
	AddNotifiedMovie(m2)
	AddDownloadedMovie(m1)
	AddDownloadedMovie(m2)
	AddDownloadingMovie(m1)
	AddDownloadingMovie(m2)
	AddFailedMovie(m1)
	AddFailedMovie(m2)

	itemToRemove := 2

	RemoveNotifiedEpisode(e2)
	RemoveDownloadingEpisode(e2)
	RemoveDownloadedEpisode(e2)
	RemoveFailedEpisode(e2)

	RemoveNotifiedMovie(m2)
	RemoveDownloadingMovie(m2)
	RemoveDownloadedMovie(m2)
	RemoveFailedMovie(m2)

	if EpisodeHasBeenNotified(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from notified episodes retention but it is still present")
	}
	if MovieHasBeenNotified(m2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from notified movie retention but it is still present")
	}
	if EpisodeHasBeenDownloaded(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from downloaded episodes retention but it is still present")
	}
	if MovieHasBeenDownloaded(m2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from downloaded movies retention but it is still present")
	}
	if EpisodeIsInDownloadProcess(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from downloading episodes retention but it is still present")
	}
	if MovieIsInDownloadProcess(m2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from downloading movies retention but it is still present")
	}
	if EpisodeHasDownloadFailed(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from failed episodes retention but it is still present")
	}
	if MovieHasDownloadFailed(m2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from failed movies retention but it is still present")
	}
}

func TestTorrentsHandling(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id: 1000,
	}
	testMovie := Movie{
		Id: 1000,
	}
	testTorrent := Torrent{
		Id: "id",
	}

	retentionData = RetentionData{
		DownloadingEpisodes: map[int]*DownloadingEpisode{
			testEpisode.Id: &DownloadingEpisode{
				Episode:     testEpisode,
				Downloading: true,
			},
		},
		DownloadingMovies: map[int]*DownloadingMovie{
			testMovie.Id: &DownloadingMovie{
				Movie:       testMovie,
				Downloading: true,
			},
		},
	}

	if !EpisodeIsDownloading(testEpisode) {
		t.Error("Episode is supposed to be downloading")
	}
	if !MovieIsDownloading(testMovie) {
		t.Error("Movie is supposed to be downloading")
	}

	ChangeEpisodeDownloadingState(testEpisode, false)
	ChangeMovieDownloadingState(testMovie, false)

	if EpisodeIsDownloading(testEpisode) {
		t.Error("Episode should not be downloading")
	}
	if MovieIsDownloading(testMovie) {
		t.Error("Movie should not be downloading")
	}

	AddFailedEpisodeTorrent(testEpisode, testTorrent)
	AddFailedMovieTorrent(testMovie, testTorrent)

	if !EpisodeIsInFailedTorrents(testEpisode, testTorrent) {
		t.Error("Expected torrent to be in failed episode torrents")
	}
	if !MovieIsInFailedTorrents(testMovie, testTorrent) {
		t.Error("Expected torrent to be in failed movie torrents")
	}

	failedEpisodeTorrentsCount := EpisodeGetFailedTorrentsCount(testEpisode)
	if failedEpisodeTorrentsCount != 1 {
		t.Error("Expected failed episode torrents count to be 1, got ", failedEpisodeTorrentsCount, " instead")
	}

	failedMovieTorrentsCount := MovieGetFailedTorrentsCount(testMovie)
	if failedMovieTorrentsCount != 1 {
		t.Error("Expected failed movie torrents count to be 1, got ", failedMovieTorrentsCount, " instead")
	}
}

func TestCurrentDownloadTorrent(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id: 1000,
	}
	testMovie := Movie{
		Id: 1000,
	}
	testTorrent := Torrent{
		Id: "id",
	}

	AddDownloadingEpisode(testEpisode)
	AddDownloadingMovie(testMovie)
	SetCurrentEpisodeDownload(testEpisode, testTorrent, "")
	SetCurrentMovieDownload(testMovie, testTorrent, "")

	episodeTorrent, _, err := GetCurrentDownloadEpisodeFromRetention(testEpisode)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	movieTorrent, _, _ := GetCurrentDownloadMovieFromRetention(testMovie)
	if episodeTorrent != testTorrent {
		t.Errorf("Expected torrent for current episode download to be '%+v', got '%+v' instead", testTorrent, episodeTorrent)
	}
	if movieTorrent != testTorrent {
		t.Errorf("Expected torrent for current movie download to be '%+v', got '%+v' instead", testTorrent, movieTorrent)
	}
}

func TestGetDownloadingItems(t *testing.T) {
	resetRetention()
	e1 := Episode{
		Id: 1,
	}
	e2 := Episode{
		Id: 2,
	}
	m1 := Movie{
		Id: 1,
	}
	m2 := Movie{
		Id: 2,
	}

	AddDownloadingEpisode(e1)
	AddDownloadingEpisode(e2)
	AddDownloadingMovie(m1)
	AddDownloadingMovie(m2)

	episodeList, _ := GetDownloadingEpisodes()
	movieList, _ := GetDownloadingMovies()

	if len(episodeList) != 2 {
		t.Errorf("Expected 2 items in downloading episodes list, got %d instead", len(episodeList))
	}
	if len(movieList) != 2 {
		t.Errorf("Expected 2 items in downloading episodes list, got %d instead", len(movieList))
	}
}

func TestTraktToken(t *testing.T) {
	token := "test"

	SaveTraktToken(token)

	savedToken := LoadTraktToken()
	if savedToken != token {
		t.Errorf("Expected token from retention to be %s, got %s instead", token, savedToken)
	}
}

func TestWatchlistRetention(t *testing.T) {
	show := "test"
	watchlist := "trakt"

	shows := GetShowsFromWatchlist(watchlist)
	if len(shows) != 0 {
		t.Error("Expected watchlist in retention to be empty")
	}

	AddShowToWatchlistRetention(watchlist, show)

	shows = GetShowsFromWatchlist(watchlist)
	if shows[0].Name != show {
		t.Errorf("Expected show from watchlist retention to be %s, got %s instead", show, shows[0].Name)
	}
}
