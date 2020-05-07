package objects

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	log "github.com/macarrie/flemzerd/logging"
)

type Movie struct {
	gorm.Model
	MediaIds          MediaIds
	MediaIdsID        int
	Title             string
	OriginalTitle     string
	CustomTitle       string
	Overview          string
	Poster            string
	Background        string
	Date              time.Time
	Notified          bool
	DownloadingItem   DownloadingItem
	DownloadingItemID int
	UseDefaultTitle   bool
}

//////////////////////////////
// Downloadable implementation
//////////////////////////////

func (m *Movie) GetId() uint {
	return m.ID
}

func (m *Movie) GetTitle() string {
	if m.CustomTitle != "" {
		return m.CustomTitle
	}

	if m.UseDefaultTitle {
		return m.Title
	}

	return m.OriginalTitle
}

func (m *Movie) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"movie": m.GetTitle(),
	})
}

func (m *Movie) GetDownloadingItem() DownloadingItem {
	return m.DownloadingItem
}

func (m *Movie) SetDownloadingItem(d DownloadingItem) {
	m.DownloadingItem = d
}

//////////////////////////////
// Cachable implementation
//////////////////////////////

func (m *Movie) IsCached(value string) bool {
	return strings.HasPrefix(value, "/cache")
}

func (m *Movie) ClearCache() {
	m.Poster = ""
	m.Background = ""
}
