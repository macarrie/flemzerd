package objects

import (
	"time"

	"github.com/jinzhu/gorm"
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
	FirstAired       time.Time
	Overview         string
	Name             string
	OriginalName     string
	CustomName       string
	Status           int
	NumberOfEpisodes int
	NumberOfSeasons  int
	Seasons          []TvSeason
	UseDefaultTitle  bool
}

type TvSeason struct {
	gorm.Model
	AirDate      *time.Time
	EpisodeCount int
	SeasonNumber int
	PosterPath   string
	TvShowID     int
}
