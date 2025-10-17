package request

import (
	"mime/multipart"
)

type ResetPinRequest struct {
	Pin          string `json:"pin" validate:"required,numeric"`
	ConfirmedPin string `json:"confirmed_pin" validate:"required,numeric"`
	AccountToken string `json:"account_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}

type ResetEmailRequest struct {
	Email        string `json:"email" validate:"required,email"`
	AccountToken string `json:"account_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}
type ResetPhoneNumberRequest struct {
	PhoneNumber  string `json:"phone_number" validate:"required"`
	AccountToken string `json:"account_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}
type ResetFullNameRequest struct {
	FullName string `json:"full_name" validate:"required"`
	DeviceID string `json:"device_id" validate:"required"`
}
type DeletetAccountRequest struct {
	AccountToken string `json:"account_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}
type BiometrictStatusRequest struct {
	Active   bool   `json:"active"`
	DeviceID string `json:"device_id" validate:"required"`
}
type ResetProfilePictureRequest struct {
	ProfilePicture *multipart.FileHeader `json:"profile_picture" validate:"required"`
	DeviceID       string                `json:"device_id" validate:"required"`
}
