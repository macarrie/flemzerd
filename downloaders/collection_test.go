package downloader

import (
	"strconv"
	"testing"

	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/retention"
)

func init() {
	retention.InitStruct()
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
	err := AddTorrent(torrent)
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
		Id: strconv.Itoa(TORRENT_DOWNLOADING),
	}

	status, _ := GetTorrentStatus(testTorrent)

	if status != TORRENT_DOWNLOADING {
		t.Errorf("Expected torrent status to be %d, got %d instead", TORRENT_DOWNLOADING, status)
	}
}

func TestWaitForDownload(t *testing.T) {
	testTorrent := Torrent{
		Id:   strconv.Itoa(TORRENT_STOPPED),
		Name: "Test torrent",
		Link: "test.torrent",
	}

	err := WaitForDownload(testTorrent)
	if err == nil {
		t.Error("Expected to get an error when download is stopped, got none instead")
	}

	testTorrent.Id = strconv.Itoa(TORRENT_SEEDING)
	err = WaitForDownload(testTorrent)
	if err != nil {
		t.Error("Expected nil error to return when download is complete, got \"", err, "\" instead")
	}
}

func TestDownload(t *testing.T) {
	show := TvShow{
		Name: "test show",
	}

	episode := Episode{
		Id:     1000,
		Name:   "Test episode",
		Season: 4,
		Number: 10,
	}

	testTorrent := Torrent{
		Id:   strconv.Itoa(TORRENT_STOPPED),
		Name: "Test torrent",
		Link: "test.torrent",
	}

	err := Download(show, episode, []Torrent{testTorrent})
	if err == nil {
		t.Error("Expected stopped torrent to generate a download error, got none instead")
	}

	retention.RemoveDownloadedEpisode(episode)
	retention.RemoveDownloadingEpisode(episode)
	testTorrent.Id = strconv.Itoa(TORRENT_SEEDING)
	err = Download(show, episode, []Torrent{testTorrent})
	if err != nil {
		t.Error("Expected seeding torrent to return no errors when downloading, got \"", err, "\" instead")
	}
}
