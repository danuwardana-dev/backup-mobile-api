package kyccontroller

import (
	_ "backend-mobile-api/docs"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	_ "backend-mobile-api/model/dto/swagger"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	kycservice "backend-mobile-api/service/kyc-service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type KYCController interface {
	VerifyKycPassport(e echo.Context) error
	VerifyKycKTP(e echo.Context) error
	SaveKycPassport(e echo.Context) error
	SaveKycKTP(e echo.Context) error
	VerifySelfie(e echo.Context) error
}

type kycController struct {
	kysService kycservice.KycService
}

func (k *kycController) VerifySelfie(e echo.Context) error {
	// TODO FOR VEIRFY SELFIE USER AFTER KTP KYC
	var (
		req request.VerifyKycSelfie
	)

	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	// SAVE KYC PASSPOPRT
	resp, _ := k.kysService.VerifyPhotoSelfie(e.Request().Context(), req)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusCreated, resp)
	}
	return e.JSON(http.StatusInternalServerError, resp)
}

// @Tags KYC
// @Summary KYC Save KTP User
// @Description Add new user KTP for KYC
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.KTPrequest true "KTP Save Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 400 {object} swagger.RegisterEmailAlreadyUsed
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/kyc/register-kyc-ktp [post]
func (k *kycController) SaveKycKTP(e echo.Context) error {
	var (
		req request.KTPrequest
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID := customResource.AuthUUID
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	// SAVE KYC PASSPOPRT
	resp, _ := k.kysService.SaveKycKTP(e.Request().Context(), req, userUUID)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusCreated, resp)
	}
	return e.JSON(http.StatusInternalServerError, resp)
}

// @Tags KYC
// @Summary KYC Save Passport User
// @Description Add new user Passport for KYC
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.PassportRequest true "Passport Save Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 400 {object} swagger.RegisterEmailAlreadyUsed
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/users/kyc/register-kyc-passport [post]
func (k *kycController) SaveKycPassport(e echo.Context) error {
	var (
		req request.PassportRequest
	)

	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	userUUID := customResource.AuthUUID

	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	// SAVE KYC PASSPOPRT
	resp, _ := k.kysService.SaveKycPassport(e.Request().Context(), req, userUUID)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusCreated, resp)
	}
	return e.JSON(http.StatusInternalServerError, resp)
}

// @Tags KYC
// @Summary KYC Verify KTP User
// @Description Endpoint for verify user KTP for KYC
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.KycRequest true "KTP Verify Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 400 {object} swagger.RegisterEmailAlreadyUsed
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/internal/verihubs/verify-ktp [post]
func (k *kycController) VerifyKycKTP(e echo.Context) error {
	var (
		err error
		req request.KycRequest
	)

	err = e.Bind(&req)
	if err != nil {
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	resp, _ := k.kysService.VerifyKycKTP(e.Request().Context(), req)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusCreated, resp)
	}
	return e.JSON(http.StatusInternalServerError, resp)
}

// @Tags KYC
// @Summary KYC Verify Passport User
// @Description Endpoint for verify user Passport for KYC
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.KycRequest true "KTP Verify Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.UserNotFound
// @Failure 400 {object} swagger.RegisterEmailAlreadyUsed
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/internal/verihubs/verify-passport [post]
func (k *kycController) VerifyKycPassport(e echo.Context) error {
	var (
		err error
		req request.KycRequest
	)

	err = e.Bind(&req)
	if err != nil {
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	resp, _ := k.kysService.VerifyKycPassport(e.Request().Context(), req)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusCreated, resp)
	}
	return e.JSON(http.StatusInternalServerError, resp)
}

func NewKycController(kycService kycservice.KycService) KYCController {
	return &kycController{
		kysService: kycService,
	}
}
