package objects

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	log "github.com/macarrie/flemzerd/logging"
)

const (
	TVSHOW_RETURNING = 1
	TVSHOW_PLANNED   = 2
	TVSHOW_ENDED     = 3
	TVSHOW_UNKNOWN   = 4
)

type TvShow struct {
	gorm.Model
	MediaIds         MediaIds
	MediaIdsID       int
	Banner           string
	Poster           string
	Background       string
	FirstAired       time.Time
	Overview         string
	Title            string
	OriginalTitle    string
	CustomTitle      string
	Status           int
	NumberOfEpisodes int
	NumberOfSeasons  int
	Seasons          []TvSeason
	UseDefaultTitle  bool
	IsAnime          bool
}

type TvSeason struct {
	gorm.Model
	AirDate      *time.Time
	EpisodeCount int
	SeasonNumber int
	PosterPath   string
	TvShowID     int
}
type SeasonDetails struct {
	Info        TvSeason
	EpisodeList []Episode
}

func (s TvShow) GetTitle() string {
	if s.CustomTitle != "" {
		return s.CustomTitle
	}

	if s.UseDefaultTitle {
		return s.Title
	}

	return s.OriginalTitle
}

//////////////////////////////
// Cachable implementation
//////////////////////////////

func (s *TvShow) IsCached(value string) bool {
	return strings.HasPrefix(value, "/cache")
}

func (s *TvShow) GetId() uint {
	return s.ID
}

func (s *TvShow) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"show": s.GetTitle(),
	})
}

func (s *TvShow) ClearCache() {
	s.Poster = ""
	s.Background = ""
}
