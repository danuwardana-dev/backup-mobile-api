package entity

import (
	"backend-mobile-api/model/enum"
	"gorm.io/gorm"
	"time"
)

type OTP struct {
	gorm.Model
	OtpCode        string          `gorm:"column:otp;type:varchar;size:7" json:"otp" json:"otp"`
	OtpPurpose     enum.OtpService `gorm:"column:otp_purpose" json:"otp_purpose"`
	OtpMethod      enum.OtpType    `gorm:"column:otp_method" json:"otp_method"`
	OtpDestination string          `gorm:"column:otp_destination" json:"otp_destination"`
	UserId         int64           `gorm:"column:user_id;type:bigint" json:"user_id"`
	UserUUID       string          `gorm:"column:user_uuid;type:varchar(36);not null" json:"user_uuid"`
	IdentityUser   string          `gorm:"column:identity_user;type:varchar;size:255" json:"identityUser"`
	VerifyKey      string          `gorm:"column:verify_key;type:varchar;size:255" json:"verify_key"`
	ExpiredAt      time.Time       `gorm:"type:timestamptz" json:"expired_at" json:"expired_at"`
	SessionId      string          `json:"session_id"`
	Status         enum.OTPStatus  `gorm:"column:status" json:"status"`
}

func (i OTP) TableName() string { return "otps" }
