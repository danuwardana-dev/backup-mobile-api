package response

import (
	"time"
)

type SelectArticleResponse struct {
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Articles []ArticleResponse `json:"articles"`
}
type ArticleResponse struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	Name           string     `json:"name"`
	Url            string     `json:"url"`
	Category       *string    `json:"category"`
	ActiveAfterDay *time.Time `json:"active_after_day"`
}
