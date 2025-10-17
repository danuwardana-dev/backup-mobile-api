package swagger

import (
	"backend-mobile-api/internal/middleware"
	"backend-mobile-api/model/dto/response"
)

type RefreshTokenInvalidPayload struct {
	StatusCode string      `json:"status_code" example:"110"`
	Message    string      `json:"message" example:"invalid request payload"`
	Error      string      `json:"error" example:"refresh_token is required"`
	Data       interface{} `json:"data"`
}
type UserNotFound struct {
	StatusCode string      `json:"status_code" example:"117"`
	Message    string      `json:"message" example:"user not found"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type Unauthorized struct {
	StatusCode string      `json:"status_code" example:"121"`
	Message    string      `json:"message" example:"unauthorized"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type DeferenceDevice struct {
	StatusCode string      `json:"status_code" example:"129"`
	Message    string      `json:"message" example:"deference device"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}

type LoginSuccess struct {
	StatusCode string `json:"status_code" example:"00"`
	Message    string `json:"message" example:"success"`
	Error      string `json:"error" example:""`
	Data       struct {
		User  UserData             `json:"user"`
		Token middleware.TokenData `json:"token"`
	} `json:"data"`
}
type UserData struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type SendOtpSuccess struct {
	StatusCode string                   `json:"status_code" example:"00"`
	Message    string                   `json:"message" example:"success"`
	Error      string                   `json:"error" example:""`
	Data       response.SendOtpResponse `json:"data"`
}
type CommonError struct {
	StatusCode string      `json:"status_code" example:"119"`
	Message    string      `json:"message" example:"server is busy"`
	Error      string      `json:"error" example:"time out"`
	Data       interface{} `json:"data"`
}
type UnverifiedUser struct {
	StatusCode string      `json:"status_code" example:"122"`
	Message    string      `json:"message" example:"unverified"`
	Error      string      `json:"error"`
	Data       interface{} `json:"data"`
}
type SuccesVerifyOtpSuccess struct {
	StatusCode string             `json:"status_code" example:"00"`
	Message    string             `json:"message" example:"success"`
	Error      string             `json:"error" example:""`
	Data       response.VerifyOtp `json:"data"`
}
type BasicSuccess struct {
	StatusCode string      `json:"status_code" example:"00"`
	Message    string      `json:"message" example:"success"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type InvalidAccessCode struct {
	StatusCode string      `json:"status_code" example:"127"`
	Message    string      `json:"message" example:"invalid access key"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type SetPinMisMatch struct {
	StatusCode string      `json:"status_code" example:"110"`
	Message    string      `json:"message" example:"invalid request payload"`
	Error      string      `json:"error" example:"deference confirmed pin"`
	Data       interface{} `json:"data"`
}
type RegisterEmailAlreadyUsed struct {
	StatusCode string      `json:"status_code" example:"110"`
	Message    string      `json:"message" example:"invalid request payload"`
	Error      string      `json:"error" example:"deference confirmed pin"`
	Data       interface{} `json:"data"`
}
type InquiryProfileResponse struct {
	StatusCode string                              `json:"status_code" example:"00"`
	Message    string                              `json:"message" example:"success"`
	Error      string                              `json:"error" example:""`
	Data       response.UserProfileInquiryResponse `json:"data"`
}
type DeferenceDeviceProfileRequest struct {
	StatusCode string      `json:"status_code" example:"143"`
	Message    string      `json:"message" example:"deference device"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}

type ExistingPinFailureResponse struct {
	StatusCode string      `json:"status_code" example:"148"`
	Message    string      `json:"message" example:" is existing pin, no update"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type ExpiredAccessFailureResponse struct {
	StatusCode string      `json:"status_code" example:"149"`
	Message    string      `json:"message" example:"expired"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type ProfileEmailAlreadyRegistered struct {
	StatusCode string      `json:"status_code" example:"145"`
	Message    string      `json:"message" example:"email already registered"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}

type UserNotFoundProfileFailureResponse struct {
	StatusCode string      `json:"status_code" example:"141"`
	Message    string      `json:"message" example:"user not found"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type ProfilePhoneNumberAlreadyRegistered struct {
	StatusCode string      `json:"status_code" example:"146"`
	Message    string      `json:"message" example:"phone number already registered"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type ProfileUnauthorized struct {
	StatusCode string      `json:"status_code" example:"121"`
	Message    string      `json:"message" example:"invalid pin"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type ProfileAccessInvalidFailureResponse struct {
	StatusCode string      `json:"status_code" example:"142"`
	Message    string      `json:"message" example:"invalid access key"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
type ProfileInvalidPayloadResponse struct {
	StatusCode string      `json:"status_code" example:"140"`
	Message    string      `json:"message" example:"invalid request payload"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}

type ProfileAccessTokenResponse struct {
	StatusCode string                       `json:"status_code" example:"00"`
	Message    string                       `json:"message" example:"success"`
	Error      string                       `json:"error" example:""`
	Data       response.AccessTokenResponse `json:"data"`
}

type Failed struct {
	StatusCode string                       `json:"status_code" example:"00"`
	Message    string                       `json:"message" example:"success"`
	Error      string                       `json:"error" example:""`
	Data       response.AccessTokenResponse `json:"data"`
}
type SuccessListArticles struct {
	StatusCode string                         `json:"status_code" example:"00"`
	Message    string                         `json:"message" example:"success"`
	Error      string                         `json:"error" example:""`
	Data       response.SelectArticleResponse `json:"data"`
}
type RecordNotFoundArticles struct {
	StatusCode string      `json:"status_code" example:"160"`
	Message    string      `json:"message" example:"record not found"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data" `
}
type DeferenceDeviceArtice struct {
	StatusCode string      `json:"status_code" example:"171"`
	Message    string      `json:"message" example:"deference device"`
	Error      string      `json:"error" example:""`
	Data       interface{} `json:"data"`
}
