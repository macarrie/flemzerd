package objects

import "github.com/jinzhu/gorm"

type DownloadingItem struct {
	gorm.Model
	Pending             bool
	Downloading         bool
	Downloaded          bool
	TorrentList         []Torrent `gorm:"foreignkey:TorrentListID"`
	CurrentDownloaderId string
	DownloadFailed      bool
	TorrentsNotFound    bool
}

func (d *DownloadingItem) CurrentTorrent() Torrent {
	for i := 0; i < len(d.TorrentList); i += 1 {
		if !d.TorrentList[i].Failed {
			return d.TorrentList[i]
		}
	}

	return Torrent{}
}

func (d *DownloadingItem) FailedTorrents() []Torrent {
	var retList []Torrent

	for _, torrent := range d.TorrentList {
		if torrent.Failed {
			retList = append(retList, torrent)
		}
	}

	return retList
}
