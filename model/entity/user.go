package entity

import (
	"backend-mobile-api/model/enum"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	UUID        string          `gorm:"column:uuid;type:varchar(36);not null" json:"uuid"`
	FullName    string          `gorm:"column:full_name;type:varchar" json:"full_name"`
	Roles       string          `gorm:"column:roles;" json:"roles"`
	Email       string          `gorm:"column:email;type:varchar" json:"email"`
	PhoneNumber string          `gorm:"column:phone_number;type:varchar" json:"phone_number"`
	Password    string          `gorm:"column:password;type:varchar" json:"password"`
	Pin         string          `gorm:"column:pin;type:varchar;size:60" json:"pin"` //hash
	DeviceID    string          `gorm:"column:device_id;type:varchar" json:"device_id"`
	Status      enum.UserStatus `gorm:"column:status;type:varchar" json:"status"`
	CreatedAt   time.Time       `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"type:timestamptz" json:"deleted_at"`
}
type UserDetail struct {
	gorm.Model     `json:"gorm_._model"`
	UserId         int64                 `gorm:"column:user_id" json:"user_id" json:"user_id,omitempty"`
	UserUUID       string                `gorm:"column:user_uuid" json:"user_uuid" json:"user_uuid,omitempty"`
	Country        string                `gorm:"column:country" json:"country" json:"country,omitempty"`
	Province       string                `gorm:"column:province" json:"province" json:"province,omitempty"`
	Regency        string                `gorm:"column:regency" json:"regency" json:"regency,omitempty"`
	District       string                `gorm:"district" json:"district"`
	Address        string                `gorm:"address" json:"address"`
	Biometric      enum.BiometricStatus  `gorm:"column:biometric" json:"biometric" json:"biometric,omitempty"`
	Status         enum.UserDetailStatus `gorm:"column:status" json:"status" json:"status,omitempty"`
	KycStatus      enum.UserKYCStatus    `gorm:"column:kyc_status" json:"kyc_status" json:"kyc_status,omitempty"`
	KycType        enum.UserKYCType      `gorm:"column:kyc_type" json:"kyc_type" json:"kyc_type,omitempty"`
	ProfilePicture string                `gorm:"column:profile_picture"  json:"profile_picture,omitempty"`
}

func (b User) TableName() string       { return "users" }
func (b UserDetail) TableName() string { return "user_details" }
