package response

import (
	"backend-mobile-api/model/enum"
	"time"
)

type UserProfileInquiryResponse struct {
	Fullname        string                `json:"fullname"`
	Email           string                `json:"email"`
	PhoneNumber     string                `json:"phone_number"`
	Status          enum.UserDetailStatus `json:"status"`
	BiometricStatus enum.BiometricStatus  `json:"biometric_status"`
}

type GetProfilePictureResponse struct {
	UserUUID string    `json:"user_uuid"`
	Url      string    `json:"url"`
	ExpireAt time.Time `json:"expire_at"`
}
