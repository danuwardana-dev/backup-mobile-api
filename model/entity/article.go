package entity

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	gorm.Model
	Name           string     `json:"name"`
	Url            string     `json:"url"`
	Category       *string    `json:"category"`
	ActiveAfterDay *time.Time `json:"active_after_day"`
}
