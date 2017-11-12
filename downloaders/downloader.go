package downloader

import . "github.com/macarrie/flemzerd/objects"

type Downloader interface {
	AddTorrent(t Torrent) error
	GetTorrentStatus(t Torrent) (int, error)
	Init() error
}
