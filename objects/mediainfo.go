package objects

import "github.com/jinzhu/gorm"

type MediaInfo struct {
	Type             string `json:"type"`
	Title            string `json:"title"`
	Container        string `json:"container"`
	Year             int    `json:"year"`
	ReleaseGroup     string `json:"release_group"`
	Format           string `json:"format"`
	ScreenSize       string `json:"screen_size"`
	AudioChannels    string `json:"audio_channels"`
	AudioCodecs      string `json:"audio_codecs"`
	Language         string `json:"language"`
	SubtitleLanguage string `json:"subtitle_language"`
	Size             string `json:"size"`
}

type EpisodeTorrentInfo struct {
	Type           string `json:"type"`
	Title          string `json:"title"`
	Year           int    `json:"year"`
	Season         int    `json:"season"`
	Episode        int    `json:"episode"`
	EpisodeCount   int    `json:"episode_count"`
	EpisodeDetails string `json:"episode_details"`
	EpisodeFormat  string `json:"episode_format"`
	Part           int    `json:"part"`
	Version        string `json:"version"`
}

type MediaIds struct {
	gorm.Model
	Name  string
	Trakt int
	Tmdb  int
	Imdb  string
	Tvdb  int
}
