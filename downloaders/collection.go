package downloader

import (
	"errors"
)

var downloadersCollection []Downloader

func AddDownloader(d Downloader) {
	downloadersCollection = append(downloadersCollection, d)
}

func AddTorrent(url string) error {
	if len(downloadersCollection) == 0 {
		return errors.New("Cannot add torrents, no downloaders are configured")
	}

	err := downloadersCollection[0].AddTorrent(url)
	return err
}
