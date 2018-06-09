package objects

import "github.com/jinzhu/gorm"

type MediaInfo struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Container string `json:"container"`
	Year      int    `json:"year"`
	Quality   string `json:"quality"`
	Size      string `json:"size"`
	Season    int    `json:"season"`
	Episode   int    `json:"episode"`
}

type MediaIds struct {
	gorm.Model
	Name  string
	Trakt int
	Tmdb  int
	Imdb  string
	Tvdb  int
}
