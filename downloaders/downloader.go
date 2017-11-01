package downloader

import (
	"errors"
)

var downloaders []Downloader

type Downloader interface {
	AddTorrent(url string) error
	Init() error
}

func AddDownloader(d Downloader) {
	downloaders = append(downloaders, d)
}

func AddTorrent(url string) error {
	if len(downloaders) == 0 {
		return errors.New("Cannot add torrents, no downloaders are configured")
	}

	err := downloaders[0].AddTorrent(url)
	return err
}
