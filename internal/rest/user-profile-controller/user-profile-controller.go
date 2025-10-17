package userProfileController

import (
	_ "backend-mobile-api/docs"
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	_ "backend-mobile-api/model/dto/swagger"
	"github.com/labstack/gommon/log"

	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	userProfileService "backend-mobile-api/service/user-profile-svc"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userProfileController struct {
	userProfileService userProfileService.UserProfileService
}

func NewUserProfileController(userProfileService userProfileService.UserProfileService) UserProfileController {
	return &userProfileController{
		userProfileService: userProfileService,
	}
}

type UserProfileController interface {
	InquiryUserProfileController(e echo.Context) error
	ResetPinController(e echo.Context) error
	ResetEmailController(e echo.Context) error
	ResetPhoneNumberController(e echo.Context) error
	AccessTokenController(e echo.Context) error
	VerifyOtpController(e echo.Context) error
	ResetFullNameController(e echo.Context) error
	DeleteAccountController(e echo.Context) error
	BiometricController(e echo.Context) error
	ResetProfileImageController(e echo.Context) error
	GetProfilePictureController(e echo.Context) error
}

// @Tags Profile
// @Summary inquiry user profile
// @Description inquiry profile data
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Success 200 {object} swagger.InquiryProfileResponse
// @Failure 401 {object} swagger.Unauthorized
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/inquiry [get]
func (ctr *userProfileController) InquiryUserProfileController(e echo.Context) error {
	var (
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID

	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "inquiry-profile"
	res := ctr.userProfileService.InquiryUserProfileService(e.Request().Context(), &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	default:
		return e.JSON(http.StatusInternalServerError, res)
	}
}

// @Tags Profile
// @Summary set new pin
// @Description user change pin when user is still login
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.ResetPinRequest true "Reset-Pin Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 400 {object} swagger.ExistingPinFailureResponse
// @Failure 401 {object} swagger.Unauthorized
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.ExpiredAccessFailureResponse
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/reset/pin [post]
func (ctr *userProfileController) ResetPinController(e echo.Context) error {
	var (
		req      request.ResetPinRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "reset-pin"
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.ResetPinService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_EXPIRED_ACCESS_CODE:
		return e.JSON(http.StatusForbidden, res)
	case pkgErr.PROFILE_IS_EXISTING_PIN:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_INVALID_ACCESS_CODE:
		e.JSON(http.StatusUnauthorized, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Profile
// @Summary set new email
// @Description user change email active
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.ResetEmailRequest true "Reset Email Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 400 {object} swagger.ProfileEmailAlreadyRegistered
// @Failure 401 {object} swagger.Unauthorized
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/reset/email [post]
func (ctr *userProfileController) ResetEmailController(e echo.Context) error {
	var (
		req      request.ResetEmailRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = ""
	userUUID = customResource.AuthUUID
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.ResetEmailService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_EXPIRED_ACCESS_CODE:
		return e.JSON(http.StatusForbidden, res)
	case pkgErr.PROFILE_EMAIL_ALREADY_REGISTERED_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_INVALID_ACCESS_CODE:
		return e.JSON(http.StatusUnauthorized, res)
	}
	return e.JSON(http.StatusInternalServerError, res)

}

// @Tags Profile
// @Summary set new phone number
// @Description user change phone number active
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.ResetPhoneNumberRequest true "Reset Phone Number Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 400 {object} swagger.ProfilePhoneNumberAlreadyRegistered
// @Failure 401 {object} swagger.Unauthorized
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/reset/phone-number [post]
func (ctr *userProfileController) ResetPhoneNumberController(e echo.Context) error {
	var (
		logData  = &dto.CustomLoggerRequest{Remarks: "reset-phone-number", Success: false}
		req      request.ResetPhoneNumberRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	defer func() {
		helpers.CustomeLogger(e.Request().Context(), logData)
	}()
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.ResetPhoneNumber(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_EXPIRED_ACCESS_CODE:
		return e.JSON(http.StatusForbidden, res)
	case pkgErr.PROFILE_PHONE_NUMBER_ALREADY_REGISTERED_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_INVALID_ACCESS_CODE:
		return e.JSON(http.StatusUnauthorized, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Profile
// @Summary request access token
// @Description access token to update profile data
// @Accept json
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.AccessTokenRequest true "Access Token Request"
// @Success 200 {object} swagger.ProfileAccessTokenResponse
// @Failure 400 {object} swagger.ProfilePhoneNumberAlreadyRegistered
// @Failure 401 {object} swagger.ProfileUnauthorized
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/access-token [post]
func (ctr *userProfileController) AccessTokenController(e echo.Context) error {
	var (
		req      request.AccessTokenRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "access-token-profile"
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.AccessTokenService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	case pkgErr.AUTH_UNAUTHORIZED_CODE:
		return e.JSON(http.StatusUnauthorized, res)
	case pkgErr.PROFILE_INVALID_ACCESS_CODE:
		e.JSON(http.StatusUnauthorized, res)
	}
	return e.JSON(http.StatusInternalServerError, res)

}

// @Tags Profile
// @Summary verify otp
// @Description verify otp to change update profile
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.VerifyOtpRequest true "Verify Otp Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 400 {object} swagger.ProfileInvalidPayloadResponse
// @Failure 401 {object} swagger.ProfileAccessInvalidFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/otp/verify [post]
func (ctr userProfileController) VerifyOtpController(e echo.Context) error {
	var (
		req      request.VerifyOtpRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "verify-otp-profile"
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.VerifyOtpService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_INVALID_ACCESS_CODE:
		return e.JSON(http.StatusUnauthorized, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_RECORD_NOT_FOUND_CODE:

		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Profile
// @Summary set new email
// @Description user change email active
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.ResetFullNameRequest true "Register Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 400 {object} swagger.ProfileInvalidPayloadResponse
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/reset/full-name [post]
func (ctr userProfileController) ResetFullNameController(e echo.Context) error {
	var (
		req      request.ResetFullNameRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "reset-full-name"
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.ResetFullNameService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	}
	return e.JSON(http.StatusInternalServerError, res)

}

// @Tags Profile
// @Summary delete account
// @Description user remove own account active
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.DeletetAccountRequest true "delete account Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 400 {object} swagger.ProfileInvalidPayloadResponse
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 401 {object} swagger.ProfileAccessInvalidFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/delete-account [post]
func (ctr userProfileController) DeleteAccountController(e echo.Context) error {
	var (
		req      request.DeletetAccountRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "delete-account"
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.DeleteAccountService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	case pkgErr.PROFILE_INVALID_ACCESS_CODE:
		return e.JSON(http.StatusUnauthorized, res)
	}
	return e.JSON(http.StatusInternalServerError, res)

}

// @Tags Profile
// @Summary set biometric active in_active
// @Description user change biometric active in_activeemail active
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.BiometrictStatusRequest true "biometric status Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/biometric [post]
func (ctr *userProfileController) BiometricController(e echo.Context) error {
	var (
		req      request.BiometrictStatusRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "biometric-status"
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.userProfileService.BiometricStatusService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Profile
// @Summary set new profile picture
// @Description user change profile picture
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param profile_picture formData file true "profile_picture"
// @Param device_id formData string true "device_id"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 401 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/profile-image [post]
func (ctr *userProfileController) ResetProfileImageController(e echo.Context) error {
	var (
		req      request.ResetProfilePictureRequest
		validate = validator.New()
		userUUID string
		err      error
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "reset-profile-picture"

	req.ProfilePicture, err = e.FormFile("profile_picture")
	if err != nil {
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	req.DeviceID = e.FormValue("device_id")

	if err = validate.Struct(req); err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return echo.NewHTTPError(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
		})
	}
	res := ctr.userProfileService.ResetProfilePictureService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_INVALID_PAYLOAD_CODE:
		return e.JSON(http.StatusBadRequest, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	}
	return e.JSON(http.StatusInternalServerError, res)

}

// @Tags Profile
// @Summary Get Profile picture url
// @Description user get url profie picture with login true
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFoundProfileFailureResponse
// @Failure 403 {object} swagger.DeferenceDeviceProfileRequest
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/profile/profile-image [get]
func (ctr *userProfileController) GetProfilePictureController(e echo.Context) error {
	var (
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID = customResource.AuthUUID
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "get-profile-picture"
	res := ctr.userProfileService.GetProfilePictureController(e.Request().Context(), &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.PROFILE_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	case pkgErr.PROFILE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}
