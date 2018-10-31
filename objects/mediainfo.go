package objects

import "github.com/jinzhu/gorm"

type MediaInfo struct {
	AudioCodec   string `json:"audio_codec"`
	AudioQuality string `json:"audio_quality"`
	Container    string `json:"container"`
	Episode      int    `json:"episode"`
	Quality      string `json:"quality"`
	ReleaseType  string `json:"release_type"`
	ReleaseGroup string `json:"release_group"`
	Season       int    `json:"season"`
	Size         string `json:"size"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	Year         int    `json:"year"`
}

type MediaIds struct {
	gorm.Model
	Title string
	Trakt int
	Tmdb  int
	Imdb  string
	Tvdb  int
}
