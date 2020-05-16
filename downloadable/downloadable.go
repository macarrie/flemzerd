package downloadable

import (
	log "github.com/macarrie/flemzerd/logging"

	. "github.com/macarrie/flemzerd/objects"
)

type Downloadable interface {
	GetId() uint
	GetTitle() string
	GetLog() *log.Entry
	GetMediaIds() MediaIds
	GetDownloadingItem() DownloadingItem
	SetDownloadingItem(DownloadingItem)
}
