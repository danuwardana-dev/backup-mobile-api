package entity

import "time"

type Notification struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:user_id;type:bigint" json:"user_id"`
	Title     string    `gorm:"column:title;type:varchar;size:255" json:"title"`
	Message   string    `gorm:"column:message;type:text" json:"message"`
	Type      string    `gorm:"column:type;type:varchar;size:10" json:"type"`
	IsRead    bool      `gorm:"column:is_read;type:tinyint" json:"is_read"`
	CreatedAt time.Time `gorm:"type:timestamptz" json:"created_at"`
}

func (n Notification) TableName() string { return "notifications" }
