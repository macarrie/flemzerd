package downloader

import (
	"testing"
)

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
	err := AddTorrent("test")
	if err == nil {
		t.Error("Got no downloaders configured, expected to have error when adding torrent")
	}
}

func TestAddTorrentWhenDownloadersAdded(t *testing.T) {
	m := MockDownloader{}
	AddDownloader(m)

	testTorrentsCount = 0
	count := testTorrentsCount
	AddTorrent("test")

	if testTorrentsCount != count+1 {
		t.Error("Expected ", count+1, " torrents, got ", testTorrentsCount)
	}
}
