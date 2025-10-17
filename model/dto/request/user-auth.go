package request

import (
	"backend-mobile-api/model/enum"
	_ "github.com/go-playground/validator/v10"
)

type BaseLogin struct {
	Type enum.LoginType `json:"type" validate:"required,oneof=EMAIL PHONE_NUMBER BIOMETRIC"`
}

type LoginEmailRequest struct {
	Type      enum.LoginType `json:"type" validate:"required,oneof=EMAIL PHONE_NUMBER BIOMETRIC" example:"EMAIL"`
	Email     string         `json:"email" validate:"required,email"`
	Pin       string         `json:"pin" validate:"required,numeric"`
	DeviceID  string         `json:"device_id" validate:"required"`
	Longitude string         `json:"longitude" validate:"required"`
	Latitude  string         `json:"latitude" validate:"required"`
	Place     string         `json:"place" validate:"required"`
}
type LoginPhoneNumberRequest struct {
	Type        enum.LoginType `json:"type" validate:"required,oneof=EMAIL PHONE_NUMBER" example:"PHONE_NUMBER"`
	PhoneNumber string         `json:"phone_number" validate:"required"`
	Pin         string         `json:"pin" validate:"required,numeric"`
	DeviceID    string         `json:"device_id" validate:"required"`
	Longitude   string         `json:"longitude" validate:"required"`
	Latitude    string         `json:"latitude" validate:"required"`
	Place       string         `json:"place" validate:"required"`
}
type LoginBiometrcRequest struct {
	Type           enum.LoginType   `json:"type" validate:"required,oneof=EMAIL PHONE_NUMBER BIOMETRIC"`
	BiometricToken BiometricRequest `json:"biometric" validate:"required"`
}
type RegisterRequest struct {
	FullName    string `json:"full_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	//DeviceID    string `json:"device_id" validate:"required"`
}
type VerifyOtpRequest struct {
	Otp        string     `json:"otp" validate:"required,numeric"`
	VerifyID   string     `json:"verify_id" validate:"required"`
	DeviceID   string     `json:"device_id" validate:"required"`
	DeviceInfo DeviceData `json:"device_info" validate:"required"`
}
type DeviceData struct {
	AppVersionCode string `json:"app_version_code" validate:"required"`
	AppVersionName string `json:"app_version_name" validate:"required"`
	Manufacturer   string `json:"manufacturer" validate:"required"`
	Brand          string `json:"brand" validate:"required"`
	Model          string `json:"model" validate:"required"`
	Product        string `json:"product" validate:"required"`
	VersionSdk     string `json:"version_sdk" validate:"required"`
	VersionRelease string `json:"version_release" validate:"required"`
}

type SetPinRequest struct {
	Pin                string `json:"pin" validate:"required,numeric"`
	ConfirmedPin       string `json:"confirmed_pin" validate:"required,numeric"`
	AccountToken       string `json:"account_token" validate:"required"`
	DeviceID           string `json:"device_id" validate:"required"`
	EmailOrPhoneNumber string `json:"email_or_phone_number" validate:"required"`
}
type SendOtpRequest struct {
	ServiceType enum.OtpService `json:"service_type" validate:"required,oneof=OTP_VERIFY_ACCOUNT OTP_RESET_PIN OTP_FORGOT_PIN"`
	OtpMethod   enum.OtpType    `json:"otp_method" validate:"required,oneof=EMAIL SMS WHATSAPP"`
	VerifyTo    string          `json:"verify_to" validate:"required"`
	DeviceID    string          `json:"device_id" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	UUID         string `json:"uuid" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}
type ForgotPinRequest struct {
	Email       string     `json:"email" validate:"required,email"`
	DeviceID    string     `json:"device_id" validate:"required"`
	CallbackUrl string     `json:"callback_url" validate:"required,url"`
	DeviceInfo  DeviceData `json:"device_info" validate:"required"`
}

type AccessTokenRequest struct {
	AccessType enum.AccessTokenType `json:"access_type" validate:"required,oneof=SET_PIN RESET_PIN RESET_EMAIL  RESET_PHONE_NUMBER DELETE_ACCOUNT"`
	Pin        string               `json:"pin" validate:"required,numeric"`
	DeviceID   string               `json:"device_id" validate:"required"`
}
