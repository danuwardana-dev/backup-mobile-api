package response

import (
	"backend-mobile-api/internal/middleware"
	"time"
)

type VerifyOtp struct {
	AccessKey string    `json:"access_key"`
	ExpireAt  time.Time `json:"expire_at"`
}
type SendOtpResponse struct {
	VerifyId string    `json:"verify_id"`
	ExpireAt time.Time `json:"expire_at"`
}
type LoginResponse struct {
	User  UserData             `json:"user"`
	Token middleware.TokenData `json:"token"`
}
type UserData struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AccessTokenResponse struct {
	AccessKey string    `json:"access_key"`
	ExpireAt  time.Time `json:"expire_at"`
}
