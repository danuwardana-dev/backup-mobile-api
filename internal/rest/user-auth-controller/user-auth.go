package userAuthController

import (
	_ "backend-mobile-api/docs"
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	_ "backend-mobile-api/model/dto/swagger"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	"backend-mobile-api/service/biometricSvc"
	userAuthSvc "backend-mobile-api/service/user-auth-svc"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type userAuthController struct {
	userAuthService  userAuthSvc.UserAuthService
	biometricService biometricSvc.BiometricService
}

func NewUserAuthController(
	userAuthService userAuthSvc.UserAuthService,
	biometricService biometricSvc.BiometricService,
) UserAuthController {
	a := &userAuthController{
		userAuthService:  userAuthService,
		biometricService: biometricService,
	}
	return a
}

type UserAuthController interface {
	RegisterController(c echo.Context) error
	LoginController(c echo.Context) error
	LogoutController(c echo.Context) error
	RefreshTokenController(c echo.Context) error
	VerifyOtpController(c echo.Context) error
	SendOtpController(c echo.Context) error
	SetPinController(c echo.Context) error
	ForgotPinController(c echo.Context) error
}

// @Tags Auth
// @Summary Register User
// @Description Register a new user
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.RegisterRequest true "Register Request"
// @Success 201 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 400 {object} swagger.RegisterEmailAlreadyUsed
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/auth/register [post]
func (ctr *userAuthController) RegisterController(e echo.Context) error {
	var (
		req      request.RegisterRequest
		validate = validator.New()
	)
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "register"
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err = validate.Struct(&req)
	if err != nil {

		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	resp := ctr.userAuthService.RegisterService(e.Request().Context(), &req, logData)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusCreated, resp)
	case pkgErr.AUHT_EMAIL_ALREADY_REGISTERED_CODE:
		return e.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_PHONE_NUMBER_ALREADY_REGISTERED_CODE:
		return e.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_USER_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, resp)
	}
	return e.JSON(http.StatusInternalServerError, resp)
}

// Auth godoc
// @Tags Auth
// @Summary Login
// @Description Login using email or username
// @Accept  json
// @Produce  json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param   loginEmailRequest body request.LoginEmailRequest false "Email Login Request"
// @Param   loginPhoneNumberRequest body request.LoginPhoneNumberRequest false "Phone number Login Request"
// @Success 200 {object} swagger.LoginSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 403 {object} swagger.DeferenceDevice
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/auth/login [post]
func (ctr *userAuthController) LoginController(c echo.Context) error {
	var (
		resp     *dto.BaseResponse
		body     = make(map[string]interface{})
		jsonBody []byte
		validate = validator.New()
	)
	logData, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = c.Request().URL.Path

	err := c.Bind(&body)
	if err != nil {
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	jsonBody, _ = json.Marshal(body)
	switch body["type"] {
	case string(enum.LOGIN_EMAIL):
		logData.Remarks = "login_email"
		var req request.LoginEmailRequest
		err := json.Unmarshal(jsonBody, &req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      err.Error(),
			})
		}
		err = validate.Struct(&req)
		if err != nil {
			err = helpers.CustomValidatePayload(err, req)

			return c.JSON(http.StatusBadRequest, dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      err.Error(),
			})
		}
		resp = ctr.userAuthService.LoginByEmailService(c.Request().Context(), &req, logData)
	case string(enum.LOGIN_PHONE_NUMBER):
		logData.Remarks = "login_phone_number"
		var req request.LoginPhoneNumberRequest
		err := json.Unmarshal(jsonBody, &req)
		if err != nil {
			logData.Error = err.Error()
			return c.JSON(http.StatusBadRequest, dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      err.Error(),
			})
		}
		err = validate.Struct(&req)
		if err != nil {
			err = helpers.CustomValidatePayload(err, req)
			logData.Error = err.Error()
			return c.JSON(http.StatusBadRequest, &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      err.Error(),
			})
		}
		resp = ctr.userAuthService.LoginByPhoneNumberService(c.Request().Context(), &req, logData)
	case string(enum.LOGIN_BIOMETRIC):
		logData.Remarks = "login_biometric"
		var req request.LoginBiometrcRequest
		err := json.Unmarshal(jsonBody, &req)
		if err != nil {
			logData.Error = err.Error()
			return c.JSON(http.StatusBadRequest, dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      err.Error(),
			})
		}
		err = validate.Struct(&req)
		if err != nil {
			err = helpers.CustomValidatePayload(err, req)
			logData.Error = err.Error()
			return c.JSON(http.StatusBadRequest, &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      err.Error(),
			})
		}
		resp = ctr.biometricService.VerifyBiometric(c.Request().Context(), &req.BiometricToken, nil, logData)

	default:
		return c.JSON(http.StatusNotFound, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
		})
	}
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return c.JSON(http.StatusOK, resp)
	case pkgErr.AUTH_USER_NOT_FOUND_CODE:
		return c.JSON(http.StatusNotFound, resp)
	case pkgErr.AUTH_UNAUTHORIZED_CODE:
		return c.JSON(http.StatusUnauthorized, resp)
	case pkgErr.AUTH_UNVERIFIED_CODE:
		return c.JSON(http.StatusUnauthorized, resp)
	case pkgErr.AUTH_DEFERENCE_DEVICE_CODE:
		return c.JSON(http.StatusForbidden, resp)
	case pkgErr.BIOMETRIC_INVALID_REQUEST_CODE:
		return c.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_BIOMETRC_INACTIVE_CODE:
		return c.JSON(http.StatusForbidden, resp)
	default:
		return c.JSON(http.StatusInternalServerError, resp)
	}

}

// @Tags Auth
// @Summary Log-out User
// @Description user log-out from apps
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.LogoutRequest true "Logout Request"
// @Success 200 {object} swagger.BasicSuccess
// @Router /api/v1/users/auth/logout [post]
func (ctr *userAuthController) LogoutController(c echo.Context) error {
	var (
		req      request.LogoutRequest
		validate = validator.New()
	)
	logData, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "logout"

	err := c.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err = validate.Struct(&req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	resp := ctr.userAuthService.LogoutService(c.Request().Context(), &req, logData)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return c.JSON(http.StatusOK, resp)
	default:
		return c.JSON(http.StatusInternalServerError, resp)
	}
}

// @Tags Auth
// @Summary Refresh token access
// @Description user reconecting after acces expired
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.RefreshTokenRequest true "Refresh Request"
// @Success 200 {object} swagger.LoginSuccess
// @Failure 400 {object} swagger.RefreshTokenInvalidPayload
// @Failure 401 {object} swagger.DeferenceDevice
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/auth/refresh [post]
func (ctr *userAuthController) RefreshTokenController(c echo.Context) error {
	var (
		req      request.RefreshTokenRequest
		validate = validator.New()
	)
	logData, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "refresh-token"

	err := c.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err = validate.Struct(&req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	resp := ctr.userAuthService.RefreshService(c.Request().Context(), req, logData)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return c.JSON(http.StatusOK, resp)
	case pkgErr.AUTH_USER_NOT_FOUND_CODE:
		return c.JSON(http.StatusUnauthorized, resp)
	case pkgErr.AUTH_UNAUTHORIZED_CODE:
		return c.JSON(http.StatusUnauthorized, resp)
	case pkgErr.AUTH_DEFERENCE_DEVICE_CODE:
		return c.JSON(http.StatusForbidden, resp)
	default:
		return c.JSON(http.StatusInternalServerError, resp)
	}
}

// @Tags Auth
// @Summary verfy OTP To get AccessTokenRequest
// @Description verify OTP to cridential access
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.VerifyOtpRequest true "Register Request"
// @Success 200 {object} swagger.SuccesVerifyOtpSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 400 {object} swagger.UnverifiedUser
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/auth/otp/verify [post]
func (ctr *userAuthController) VerifyOtpController(c echo.Context) error {
	var (
		req      request.VerifyOtpRequest
		validate = validator.New()
	)
	logData, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "verify-otp"
	err := c.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err = validate.Struct(&req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	resp := ctr.userAuthService.VerifyOtpService(c.Request().Context(), &req, logData)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return c.JSON(http.StatusOK, resp)
	case pkgErr.AUTH_INVALID_OTP_CODE:
		return c.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_UNVERIFIED_CODE:
		return c.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_ALREADY_VERIFIED_CODE:
		return c.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_USER_NOT_FOUND_CODE:
		return c.JSON(http.StatusNotFound, resp)
	default:
		return c.JSON(http.StatusInternalServerError, resp)

	}
}

// @Tags Auth
// @Summary Send Otp To get AccessTokenRequest
// @Description OTP to cridential access
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.SendOtpRequest true "Register Request"
// @Success 200 {object} swagger.SendOtpSuccess
// @Failure 401 {object} swagger.UserNotFound
// @Failure 403 {object} swagger.UnverifiedUser
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/auth/otp/send [post]
func (ctr *userAuthController) SendOtpController(c echo.Context) error {
	var (
		req      request.SendOtpRequest
		validate = validator.New()
	)
	logData, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "send-otp"

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err = validate.Struct(&req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)

		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	logData.Remarks = fmt.Sprintf("send-otp-%s-%s", req.OtpMethod, req.ServiceType)
	resp := ctr.userAuthService.SendOtpService(c.Request().Context(), &req, logData)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return c.JSON(http.StatusOK, resp)
	case pkgErr.AUTH_USER_NOT_FOUND_CODE:
		return c.JSON(http.StatusNotFound, resp)
	case pkgErr.AUTH_ALREADY_VERIFIED_CODE:
		return c.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_UNVERIFIED_CODE:
		return c.JSON(http.StatusForbidden, resp)
	case pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE:
		return c.JSON(http.StatusBadRequest, resp)
	default:
		return c.JSON(http.StatusInternalServerError, resp)

	}
}

// @Tags Auth
// @Summary Set Pin new pin
// @Description Set Pin with Req
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.SetPinRequest true "Register Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 403 {object} swagger.InvalidAccessCode
// @Failure 400 {object} swagger.SetPinMisMatch
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/auth/set-pin [post]
func (ctr *userAuthController) SetPinController(c echo.Context) error {
	var (
		req      request.SetPinRequest
		validate = validator.New()
	)
	logData, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "set-pin"
	if err := c.Bind(&req); err != nil {
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err := validate.Struct(&req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	resp := ctr.userAuthService.SetPinService(c.Request().Context(), &req, logData)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return c.JSON(http.StatusOK, resp)
	case pkgErr.AUTH_USER_NOT_FOUND_CODE:
		return c.JSON(http.StatusNotFound, resp)
	case pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE:
		return c.JSON(http.StatusBadRequest, resp)
	case pkgErr.AUTH_INVALID_ACCESS_CODE:
		return c.JSON(http.StatusForbidden, resp)
	default:
		return c.JSON(http.StatusInternalServerError, resp)

	}
}

// @Tags Auth
// @Summary Forgot Pin
// @Description request access to renew pin on forgot moment
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.ForgotPinRequest true "Register Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/auth/forgot-pin [post]
func (ctr *userAuthController) ForgotPinController(c echo.Context) error {
	var (
		req      request.ForgotPinRequest
		validate = validator.New()
	)
	logData, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "forgot-pin"

	if err := c.Bind(&req); err != nil {
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err := validate.Struct(&req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return c.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	rest := ctr.userAuthService.ForgotPinService(c.Request().Context(), &req, logData)
	switch rest.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return c.JSON(http.StatusOK, rest)
	case pkgErr.AUTH_USER_NOT_FOUND_CODE:
		return c.JSON(http.StatusNotFound, rest)
	default:
		return c.JSON(http.StatusInternalServerError, rest)

	}

}
