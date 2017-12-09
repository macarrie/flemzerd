package downloader

import . "github.com/macarrie/flemzerd/objects"

type Downloader interface {
	AddTorrent(t Torrent) (string, error)
	AddTorrentMapping(string, string)
	RemoveTorrent(t Torrent) error
	GetTorrentStatus(t Torrent) (int, error)
	Init() error
}
