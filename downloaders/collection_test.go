package downloader

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"

	log "github.com/macarrie/flemzerd/logging"
	mock "github.com/macarrie/flemzerd/mocks"
	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	log.Setup(true)

	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
}

func TestStatus(t *testing.T) {
	AddDownloader(mock.Downloader{})
	AddDownloader(mock.Downloader{})

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected status module length to be 2, got %d instead", len(mods))
	}
	if err != nil {
		t.Error("Expected not to have an error, got one instead")
	}

	AddDownloader(mock.ErrorDownloader{})
	_, err = Status()
	if err == nil {
		t.Error("Expected to have error whenchecking status of mock.ErrorDownloader")
	}
}

func TestReset(t *testing.T) {
	d := mock.Downloader{}
	AddDownloader(d)
	Reset()

	if len(downloadersCollection) != 0 {
		t.Error("Expected downloader collection to be empty after reset")
	}
}

func TestAddDownloader(t *testing.T) {
	downloadersLength := len(downloadersCollection)
	m := mock.Downloader{}
	AddDownloader(m)

	if len(downloadersCollection) != downloadersLength+1 {
		t.Error("Expected ", downloadersLength+1, " downloaders, got ", len(downloadersCollection))
	}
}

func TestAddTorrentWhenNoDownloadersAdded(t *testing.T) {
	downloadersCollection = []Downloader{}
	torrent := Torrent{
		Link: "test",
	}
	_, err := AddTorrent(torrent)
	if err == nil {
		t.Error("Got no downloaders configured, expected to have error when adding torrent")
	}
}

func TestAddTorrentWhenDownloadersAdded(t *testing.T) {
	m := mock.Downloader{}
	torrent := Torrent{
		Link: "test",
	}
	AddDownloader(m)

	count := m.GetTorrentCount()
	AddTorrent(torrent)

	if m.GetTorrentCount() != count+1 {
		t.Error("Expected ", count+1, " torrents, got ", m.GetTorrentCount())
	}

	downloadersCollection = []Downloader{mock.ErrorDownloader{}}
	_, err := AddTorrent(torrent)
	if err == nil {
		t.Error("Expected to have error when adding torrent to mock.ErrrDownloader")
	}
}

func TestRemoveTorrentWhenNoDownloadersAdded(t *testing.T) {
	downloadersCollection = []Downloader{}
	torrent := Torrent{
		Link: "test",
	}
	err := RemoveTorrent(torrent)
	if err == nil {
		t.Error("Got no downloaders configured, expected to have error when removing torrent")
	}
}

func TestRemoveTorrentWhenDownloadersAdded(t *testing.T) {
	m := mock.Downloader{}
	torrent := Torrent{
		Link: "test",
	}
	AddDownloader(m)

	AddTorrent(torrent)
	count := m.GetTorrentCount()
	RemoveTorrent(torrent)

	if m.GetTorrentCount() != count-1 {
		t.Error("Expected ", count-1, " torrents, got ", m.GetTorrentCount())
	}
}

func TestAddTorrentMapping(t *testing.T) {
	// Useless test because mapping handling is specific to each downloader

	downloadersCollection = []Downloader{mock.Downloader{}}

	AddTorrentMapping("test", "test")
}

func TestGetTorrentStatus(t *testing.T) {
	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_DOWNLOADING),
	}

	status, _ := GetTorrentStatus(testTorrent)

	if status != TORRENT_SEEDING {
		t.Errorf("Expected torrent status to be %d, got %d instead", TORRENT_SEEDING, status)
	}
}

func TestWaitForDownload(t *testing.T) {
	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}

	downloadersCollection = []Downloader{mock.ErrorDownloader{}}
	err, _ := WaitForDownload(context.Background(), testTorrent)
	if err == nil {
		t.Error("Expected to get an error when download is stopped, got none instead")
	}

	downloadersCollection = []Downloader{mock.Downloader{}}
	testTorrent.TorrentId = strconv.Itoa(TORRENT_SEEDING)
	err, _ = WaitForDownload(context.Background(), testTorrent)
	if err != nil {
		t.Error("Expected nil error to return when download is complete, got \"", err, "\" instead")
	}
}

func TestDownloadEpisode(t *testing.T) {
	show := TvShow{
		Title:         "test show",
		OriginalTitle: "test show",
	}

	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}

	testTorrent2 := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}

	episode := Episode{
		Title:  "Test episode",
		TvShow: show,
		Season: 4,
		Number: 10,
	}

	downloadersCollection = []Downloader{mock.ErrorDownloader{}}
	err := DownloadEpisode(episode, []Torrent{testTorrent}, false)
	if err == nil {
		t.Error("Expected stopped torrent to generate a download error, got none instead")
	}

	downloadersCollection = []Downloader{mock.Downloader{}}
	testTorrent.TorrentId = strconv.Itoa(TORRENT_SEEDING)
	err = DownloadEpisode(episode, []Torrent{testTorrent2, testTorrent}, false)
	if err != nil {
		t.Error("Expected seeding torrent to return no errors when downloading, got \"", err, "\" instead")
	}

	downloadersCollection = []Downloader{mock.ErrorDownloader{}}
	err = DownloadEpisode(episode, []Torrent{testTorrent2, testTorrent}, false)
	if err == nil {
		t.Error("Expected torrent download to return an error because torrent cannot be added to downloader")
	}

	downloadersCollection = []Downloader{mock.DLErrorDownloader{}}
	err = DownloadEpisode(episode, []Torrent{testTorrent2, testTorrent}, false)
	if err == nil {
		t.Error("Expected torrent download to return an error because torrent status are unknown")
	}
}

func TestDownloadMovie(t *testing.T) {
	testMovie := Movie{
		Title:         "test movie",
		OriginalTitle: "test movie",
	}

	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}
	testTorrent2 := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent 2",
		Link:      "test.torrent",
	}

	downloadersCollection = []Downloader{mock.ErrorDownloader{}}
	err := DownloadMovie(testMovie, []Torrent{testTorrent}, false)
	if err == nil {
		t.Error("Expected stopped torrent to generate a download error, got none instead")
	}

	downloadersCollection = []Downloader{mock.Downloader{}}
	testTorrent.TorrentId = strconv.Itoa(TORRENT_SEEDING)
	err = DownloadMovie(testMovie, []Torrent{testTorrent, testTorrent2, testTorrent}, false)
	if err != nil {
		t.Error("Expected seeding torrent to return no errors when downloading, got \"", err, "\" instead")
	}

	downloadersCollection = []Downloader{mock.ErrorDownloader{}}
	err = DownloadMovie(testMovie, []Torrent{testTorrent2, testTorrent}, false)
	if err == nil {
		t.Error("Expected torrent download to return an error because torrent cannot be added to downloader")
	}

	downloadersCollection = []Downloader{mock.DLErrorDownloader{}}
	err = DownloadMovie(testMovie, []Torrent{testTorrent2, testTorrent}, false)
	if err == nil {
		t.Error("Expected torrent download to return an error because torrent status are unknown")
	}
}

func TestMoveEpisodeToLibrary(t *testing.T) {
	configuration.Config.Library.ShowPath = "/tmp/flemzerd_test_tmp"

	os.RemoveAll(configuration.Config.Library.ShowPath)
	os.MkdirAll(configuration.Config.Library.ShowPath, 0755)

	tmpTorrentFile := fmt.Sprintf("%s/test_flemzerd_move_media", configuration.Config.Library.ShowPath)
	os.Create(tmpTorrentFile)

	testTorrent := Torrent{
		Name:        "test_torrent",
		DownloadDir: tmpTorrentFile,
	}

	testEpisode := Episode{
		Title: "test episode",
		TvShow: TvShow{
			Title:         "test show",
			OriginalTitle: "test show",
		},
		Season: 1,
		Number: 1,
		DownloadingItem: DownloadingItem{
			CurrentTorrent: testTorrent,
		},
	}

	err := MoveEpisodeToLibrary(&testEpisode)
	if err != nil {
		t.Errorf("Episode could not be moved to library: %s", err.Error())
	}
}

func TestMoveMovieToLibrary(t *testing.T) {
	configuration.Config.Library.MoviePath = "/tmp/flemzerd_test_tmp"

	os.RemoveAll(configuration.Config.Library.MoviePath)
	os.MkdirAll(configuration.Config.Library.MoviePath, 0755)

	tmpTorrentFile := fmt.Sprintf("%s/test_flemzerd_move_media", configuration.Config.Library.MoviePath)
	os.Create(tmpTorrentFile)

	testTorrent := Torrent{
		Name:        "test_torrent",
		DownloadDir: tmpTorrentFile,
	}

	testMovie := Movie{
		Title:         "test movie",
		OriginalTitle: "test movie",
		DownloadingItem: DownloadingItem{
			CurrentTorrent: testTorrent,
		},
	}

	err := MoveMovieToLibrary(&testMovie)
	if err != nil {
		t.Errorf("Movie could not be moved to library: %s", err.Error())
	}
}

func TestFillEpisodeToDownload(t *testing.T) {
	torrent1 := Torrent{
		TorrentId: "1",
	}
	torrent2 := Torrent{
		TorrentId: "2",
	}

	episode := Episode{
		Model: gorm.Model{
			ID: 1000,
		},
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}

	episode.DownloadingItem.FailedTorrents = append(episode.DownloadingItem.FailedTorrents, torrent1)
	torrentList := FillEpisodeToDownloadTorrentList(&episode, []Torrent{torrent1, torrent2})
	if len(torrentList) != 1 {
		t.Errorf("Expected torrent list to have 1 torrent, got %d instead", len(torrentList))
	}

	episode.DownloadingItem.FailedTorrents = []Torrent{}
	torrentList = FillEpisodeToDownloadTorrentList(&episode, []Torrent{
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
	})
	if len(torrentList) > 10 {
		t.Errorf("Expected torrent list no to be bigger than 10 items, got %d instead", len(torrentList))
	}
}

func TestFillMovieToDownload(t *testing.T) {
	torrent1 := Torrent{
		TorrentId: "1",
	}
	torrent2 := Torrent{
		TorrentId: "2",
	}

	movie := Movie{
		Model: gorm.Model{
			ID: 1000,
		},
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}

	movie.DownloadingItem.FailedTorrents = append(movie.DownloadingItem.FailedTorrents, torrent1)
	torrentList := FillMovieToDownloadTorrentList(&movie, []Torrent{torrent1, torrent2})
	if len(torrentList) != 1 {
		t.Errorf("Expected torrent list to have 1 torrent, got %d instead", len(torrentList))
	}

	movie.DownloadingItem.FailedTorrents = []Torrent{}
	torrentList = FillMovieToDownloadTorrentList(&movie, []Torrent{
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
		torrent1,
	})
	if len(torrentList) > 10 {
		t.Errorf("Expected torrent list no to be bigger than 10 items, got %d instead", len(torrentList))
	}
}

func TestAbortDownload(t *testing.T) {
	downloadersCollection = []Downloader{mock.StalledDownloader{}}

	torrent1 := Torrent{
		TorrentId: "1",
	}
	torrent2 := Torrent{
		TorrentId: "2",
	}

	movie := Movie{
		Model: gorm.Model{
			ID: 1000,
		},
		Title:         "test movie",
		OriginalTitle: "test movie",
	}
	episode := Episode{
		Model: gorm.Model{
			ID: 1000,
		},
		Title:  "test episode",
		Season: 1,
		Number: 1,
		TvShow: TvShow{
			Model: gorm.Model{
				ID: 1000,
			},
			Title:         "test show",
			OriginalTitle: "test show",
		},
	}

	go DownloadEpisode(episode, []Torrent{torrent1, torrent2}, false)
	time.Sleep(2 * time.Second)
	AbortEpisodeDownload(&episode)
	if episode.DownloadingItem.Downloading {
		t.Error("Expected download to be stopped")
	}

	go DownloadMovie(movie, []Torrent{torrent1, torrent2}, false)
	time.Sleep(2 * time.Second)
	AbortMovieDownload(&movie)
	if movie.DownloadingItem.Downloading {
		t.Error("Expected download to be stopped")
	}

	// Test same download with recovery mode enabled
	go DownloadEpisode(episode, []Torrent{torrent1, torrent2}, true)
	time.Sleep(2 * time.Second)
	AbortEpisodeDownload(&episode)
	if episode.DownloadingItem.Downloading {
		t.Error("Expected download to be stopped")
	}

	go DownloadMovie(movie, []Torrent{torrent1, torrent2}, true)
	time.Sleep(2 * time.Second)
	AbortMovieDownload(&movie)
	if movie.DownloadingItem.Downloading {
		t.Error("Expected download to be stopped")
	}
	downloadersCollection = []Downloader{mock.Downloader{}}
}

func TestGetDownloader(t *testing.T) {
	downloadersCollection = []Downloader{mock.Downloader{}}

	if _, err := GetDownloader("Unknown"); err == nil {
		t.Error("Expected to have error when getting unknown notifier, got none")
	}

	if _, err := GetDownloader("Downloader"); err != nil {
		t.Errorf("Got error while retrieving known notifier: %s", err.Error())
	}
}
