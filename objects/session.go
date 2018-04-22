package objects

import "github.com/jinzhu/gorm"

type SessionData struct {
	gorm.Model
	TraktToken string
}
