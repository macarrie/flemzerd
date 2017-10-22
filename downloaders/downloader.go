package downloader

import (
	log "flemzerd/logging"
)

var downloaders []Downloader

type Downloader interface {
    AddTorrent(url string) error
    Init() error
}

func AddDownloader(d Downloader) {
    log.Debug(d)
    downloaders = append(downloaders, d)
}

func AddTorrent(url string) error {
    downloaders[0].AddTorrent(url)
    return nil
}
