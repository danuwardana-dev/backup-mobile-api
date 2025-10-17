package entity

import (
	"gorm.io/gorm"
	"time"
)

type TokenBlacklist struct {
	gorm.Model
	Token       string    `gorm:"column:token;type:varchar" json:"token"`
	BlacklistAt time.Time `gorm:"type:timestamptz" json:"blacklist_at"`
	ExpiredAt   time.Time `gorm:"type:timestamptz" json:"expired_at"`
	Description string    `gorm:"column:description;type:varchar" json:"description"`
}

func (i TokenBlacklist) TableName() string { return "token_blacklists" }
