package downloader

import (
//log "github.com/macarrie/flemzerd/logging"
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
	err := downloaders[0].AddTorrent(url)
	return err
}
