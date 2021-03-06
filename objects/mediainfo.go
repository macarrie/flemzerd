package objects

import "github.com/jinzhu/gorm"

type MediaInfo struct {
	Id           string
	AudioCodec   string `json:"audio_codec"`
	AudioQuality string `json:"audio_quality"`
	Container    string `json:"container"`
	Episode      int    `json:"episode"`
	Quality      string `json:"quality"`
	Raw          string `json:"raw"`
	ReleaseType  string `json:"release_type"`
	ReleaseGroup string `json:"release_group"`
	Season       int    `json:"season"`
	Size         string `json:"size"`
	Title        string `json:"title"`
	Type         string `json:"media_type"`
	Year         int    `json:"year"`
}

type MediaInfoEpisodes map[int]MediaInfo
type MediaInfoSeasons map[int]MediaInfoEpisodes
type MediaInfoGroupedByShow map[string]MediaInfoSeasons

type MediaIds struct {
	gorm.Model
	Title string
	Trakt int
	Tmdb  int
	Imdb  string
	Tvdb  int
}
