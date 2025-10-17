package entity

import (
	"backend-mobile-api/model/enum"
	"gorm.io/gorm"
	"time"
)

type AccessState struct {
	gorm.Model
	AccessType  enum.AccessTokenType `json:"access_type" validate:"required"`
	UserId      int64                `gorm:"column:user_id;type:bigint" json:"user_id"`
	UserUUID    string               `gorm:"column:user_uuid;type:varchar(36);not null" json:"user_uuid"`
	DeviceId    string               `gorm:"column:device_id;type:bigint" json:"device_id"`
	AccessToken string               `gorm:"column:access_token;type:varchar;size:50" json:"access_token"`
	ExpiredAt   time.Time            `gorm:"column:expired_at;type:datetime" json:"expired_at"`
	Used        bool                 `gorm:"column:used;type:boolean" json:"used"`
}

func (rp AccessState) TableName() string {
	return "access_states"
}
