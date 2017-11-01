package downloader

import (
	"testing"
)

type MockDownloader struct{}

var testTorrentsCount int

func (d MockDownloader) AddTorrent(url string) error {
	testTorrentsCount += 1
	return nil
}

func (d MockDownloader) Init() error {
	return nil
}

func TestAddDownloader(t *testing.T) {
	downloadersLength := len(downloaders)
	m := MockDownloader{}
	AddDownloader(m)

	if len(downloaders) != downloadersLength+1 {
		t.Error("Expected ", downloadersLength+1, " downloaders, got ", len(downloaders))
	}
}

func TestAddTorrentWhenNoDownloadersAdded(t *testing.T) {
	downloaders = []Downloader{}
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
