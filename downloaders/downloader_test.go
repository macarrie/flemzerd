package downloader

import (
	. "github.com/macarrie/flemzerd/objects"
)

type MockDownloader struct{}

var testTorrentsCount int

func (d MockDownloader) AddTorrent(t Torrent) error {
	testTorrentsCount += 1
	return nil
}

func (d MockDownloader) Init() error {
	return nil
}

func (d MockDownloader) GetTorrentStatus(t Torrent) (int, error) {
	return 0, nil
}
