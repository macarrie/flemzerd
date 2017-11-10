package downloader

import (
	"errors"

	. "github.com/macarrie/flemzerd/objects"
)

var downloadersCollection []Downloader

func AddDownloader(d Downloader) {
	downloadersCollection = append(downloadersCollection, d)
}

func AddTorrent(t Torrent) error {
	if len(downloadersCollection) == 0 {
		return errors.New("Cannot add torrents, no downloaders are configured")
	}

	err := downloadersCollection[0].AddTorrent(t)
	return err
}
