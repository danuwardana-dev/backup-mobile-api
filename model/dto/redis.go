package dto

import (
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
)

type SendOtp struct {
	OtpCode        string          `json:"otp" json:"otp"`
	OtpPurpose     enum.OtpService `json:"otp_purpose"`
	OtpMethod      enum.OtpType    `json:"otp_method"`
	OtpDestination string          `gorm:"column:otp_destination" json:"otp_destination"`
	UserId         int64           `json:"user_id"`
	UserUUID       string          `json:"user_uuid"`
	VerifyKey      string          `json:"verify_key"`
}

type ProfileUpdateRequest struct {
	Field enum.ProfileFieldUpdate
	Value string
	User  *entity.User
}
