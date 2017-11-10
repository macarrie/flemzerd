package downloader

import . "github.com/macarrie/flemzerd/objects"

type Downloader interface {
	AddTorrent(t Torrent) error
	Init() error
}
