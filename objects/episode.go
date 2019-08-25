package objects

import (
	"time"

	"github.com/jinzhu/gorm"

	log "github.com/macarrie/flemzerd/logging"
)

type Episode struct {
	gorm.Model
	MediaIds          MediaIds
	MediaIdsID        uint
	TvShow            TvShow
	TvShowID          uint
	AbsoluteNumber    int
	Number            int
	Season            int
	Title             string
	Date              time.Time
	Overview          string
	Notified          bool
	DownloadingItem   DownloadingItem
	DownloadingItemID uint
}

//////////////////////////////
// Downloadable implementation
//////////////////////////////

func (e *Episode) GetId() uint {
	return e.ID
}

func (e *Episode) GetTitle() string {
	return e.Title
}

func (e *Episode) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"show":    e.TvShow.GetTitle(),
		"season":  e.Season,
		"number":  e.Number,
		"episode": e.GetTitle(),
	})
}

func (e *Episode) GetDownloadingItem() DownloadingItem {
	return e.DownloadingItem
}

func (e *Episode) SetDownloadingItem(d DownloadingItem) {
	e.DownloadingItem = d
}
