package request

import "time"

type NewArticleRequest struct {
	DeviceID   string     `json:"device_id" validate:"required"`
	NewArticle NewArticle `json:"new_article" validate:"required"`
}
type UpdateArticleRequest struct {
	DeviceID string  `json:"device_id" validate:"required"`
	Article  Article `json:"article" validate:"required"`
}
type DeleteArticleRequest struct {
	DeviceID string  `json:"device_id" validate:"required"`
	Article  Article `json:"article" validate:"required"`
}
type SelectListArticleRequest struct {
	DeviceID string        `json:"device_id" validate:"required"`
	Article  ArticleByList `json:"article" validate:"required"`
}
type ArticleByList struct {
	Category       *string    `json:"category"`
	ActiveAfterDay *time.Time `json:"active_after_day"`
	Limit          int        `json:"limit"`
	Offset         int        `json:"offset"`
}
type Article struct {
	Id             uint       `json:"id" validate:"required,numeric"`
	Name           *string    `json:"name"`
	Url            *string    `json:"url"`
	Category       *string    `json:"category"`
	ActiveAfterDay *time.Time `json:"active_after_day"`
}
type NewArticle struct {
	Name           string     `json:"name" validate:"required"`
	Url            string     `json:"url" validate:"required,url"`
	Category       *string    `json:"category"`
	ActiveAfterDay *time.Time `json:"active_after_day"`
}
type CountArticleRequest struct {
	Category       *string    `json:"category"`
	ActiveAfterDay *time.Time `json:"active_after_day"`
}
