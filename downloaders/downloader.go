package downloader

import . "github.com/macarrie/flemzerd/objects"

type Downloader interface {
	Init() error
	Status() (Module, error)
	AddTorrent(t Torrent) (string, error)
	AddTorrentMapping(string, string)
	RemoveTorrent(t Torrent) error
	GetTorrentStatus(t Torrent) (int, error)
}
