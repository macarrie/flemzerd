package downloader

import (
	"strconv"

	. "github.com/macarrie/flemzerd/objects"
)

type MockDownloader struct{}

var testTorrentsCount int

func (d MockDownloader) AddTorrent(t Torrent) (string, error) {
	testTorrentsCount += 1
	return "id", nil
}

func (d MockDownloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}

func (d MockDownloader) RemoveTorrent(t Torrent) error {
	if testTorrentsCount > 0 {
		testTorrentsCount -= 1
	}
	return nil
}

func (d MockDownloader) Init() error {
	return nil
}

func (d MockDownloader) GetTorrentStatus(t Torrent) (int, error) {
	status, _ := strconv.Atoi(t.Id)
	return status, nil
}
