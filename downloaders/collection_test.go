package downloader

import (
	"strconv"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/macarrie/flemzerd/db"
	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
}

func TestStatus(t *testing.T) {
	m1 := MockDownloader{}
	m2 := MockDownloader{}
	AddDownloader(m1)
	AddDownloader(m2)

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected status module length to be 2, got %d instead", len(mods))
	}
	if err == nil {
		t.Error("Expected to have an error, got none instead")
	}
}

func TestReset(t *testing.T) {
	d := MockDownloader{}
	AddDownloader(d)
	Reset()

	if len(downloadersCollection) != 0 {
		t.Error("Expected downloader collection to be empty after reset")
	}
}

func TestAddDownloader(t *testing.T) {
	downloadersLength := len(downloadersCollection)
	m := MockDownloader{}
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
	m := MockDownloader{}
	torrent := Torrent{
		Link: "test",
	}
	AddDownloader(m)

	testTorrentsCount = 0
	count := testTorrentsCount
	AddTorrent(torrent)

	if testTorrentsCount != count+1 {
		t.Error("Expected ", count+1, " torrents, got ", testTorrentsCount)
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
	m := MockDownloader{}
	torrent := Torrent{
		Link: "test",
	}
	AddDownloader(m)

	AddTorrent(torrent)
	testTorrentsCount = 1
	count := testTorrentsCount
	RemoveTorrent(torrent)

	if testTorrentsCount != count-1 {
		t.Error("Expected ", count-1, " torrents, got ", testTorrentsCount)
	}
}

func TestGetTorrentStatus(t *testing.T) {
	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_DOWNLOADING),
	}

	status, _ := GetTorrentStatus(testTorrent)

	if status != TORRENT_DOWNLOADING {
		t.Errorf("Expected torrent status to be %d, got %d instead", TORRENT_DOWNLOADING, status)
	}
}

func TestWaitForDownload(t *testing.T) {
	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}

	err := WaitForDownload(testTorrent)
	if err == nil {
		t.Error("Expected to get an error when download is stopped, got none instead")
	}

	testTorrent.TorrentId = strconv.Itoa(TORRENT_SEEDING)
	err = WaitForDownload(testTorrent)
	if err != nil {
		t.Error("Expected nil error to return when download is complete, got \"", err, "\" instead")
	}
}

func TestDownloadEpisode(t *testing.T) {
	show := TvShow{
		Name: "test show",
	}

	episode := Episode{
		Model: gorm.Model{
			ID: 1000,
		},
		Name:   "Test episode",
		Season: 4,
		Number: 10,
	}

	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}

	err := DownloadEpisode(&show, &episode, []Torrent{testTorrent})
	if err == nil {
		t.Error("Expected stopped torrent to generate a download error, got none instead")
	}

	testTorrent.TorrentId = strconv.Itoa(TORRENT_SEEDING)
	err = DownloadEpisode(&show, &episode, []Torrent{testTorrent})
	if err != nil {
		t.Error("Expected seeding torrent to return no errors when downloading, got \"", err, "\" instead")
	}
}

func TestDownloadMovie(t *testing.T) {
	testMovie := Movie{
		Model: gorm.Model{
			ID: 1000,
		},
		Title: "test movie",
	}

	testTorrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}

	err := DownloadMovie(&testMovie, []Torrent{testTorrent})
	if err == nil {
		t.Error("Expected stopped torrent to generate a download error, got none instead")
	}

	testTorrent.TorrentId = strconv.Itoa(TORRENT_SEEDING)
	err = DownloadMovie(&testMovie, []Torrent{testTorrent})
	if err != nil {
		t.Error("Expected seeding torrent to return no errors when downloading, got \"", err, "\" instead")
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
}

func TestRecoverFromRetention(t *testing.T) {
	db.ResetDb()

	testTorrent := Torrent{
		TorrentId: "1000",
	}
	testTorrent2 := Torrent{
		TorrentId: "1001",
	}

	testEpisode := Episode{
		Model: gorm.Model{
			ID: 1000,
		},
		Name:   "testEpisode",
		Season: 1,
		Number: 4,
		DownloadingItem: DownloadingItem{
			Downloading:         true,
			CurrentTorrent:      testTorrent,
			CurrentDownloaderId: "id",
		},
	}

	testMovie := Movie{
		Model: gorm.Model{
			ID: 1001,
		},
		Title: "testMovie",
		DownloadingItem: DownloadingItem{
			Downloading:         true,
			CurrentTorrent:      testTorrent2,
			CurrentDownloaderId: "id",
		},
	}

	db.Client.Create(&testEpisode)
	db.Client.Create(&testMovie)

	RecoverFromRetention()
	// TODO: check objects states is correct
}
