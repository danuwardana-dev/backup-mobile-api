package middleware

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
)

type ExcludeURLValidation struct {
	Authorization      AuthorizationMiddlewarePath
	MandatoryHeader    []string
	ValidationSignaure []string
	ValidationXNonce   []string
}
type AuthorizationMiddlewarePath struct {
	ExcludeURL   []string
	AccessByRole map[string][]enum.RolesEnum
}

type ListRouth map[string]bool

func (svc *customMiddleware) AccessMiddleware(excludeUrl *ExcludeURLValidation, routh *ListRouth) echo.MiddlewareFunc {
	return func(Next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				SourceData = &dto.ContextValue{
					HeaderContentType:   c.Request().Header.Get(echo.HeaderContentType),
					HeaderUserAgent:     c.Request().Header.Get(string(enum.HEADER_USER_AGENT)),
					HeaderXTimestamp:    c.Request().Header.Get(string(enum.HEADER_X_TIMESTAMP)),
					HeaderXSignature:    c.Request().Header.Get(string(enum.HEADER_X_SIGNATURE)),
					HeaderXRealIp:       c.RealIP(),
					HeaderXNonce:        c.Request().Header.Get(string(enum.HEADER_X_NONCE)),
					HeaderXDeviceID:     c.Request().Header.Get(string(enum.HEADER_X_DEVICE_ID)),
					HeaderXLatitude:     c.Request().Header.Get(string(enum.HEADER_X_LATITUDE)),
					HeaderXLongitude:    c.Request().Header.Get(string(enum.HEADER_X_LONGITUDE)),
					HeaderXApiKey:       c.Request().Header.Get(string(enum.HEADER_X_API_KEY)),
					HeaderAuthorization: c.Request().Header.Get(echo.HeaderAuthorization),
					HeaderRequestId:     c.Response().Header().Get(echo.HeaderXRequestID),
					HeaderHost:          c.Request().URL.Host,
					HeaderPath:          c.Request().URL.Path,
					RequestPath:         c.Path(),
					HeaderMethod:        c.Request().Method,
				}
				errRes  *dto.BaseResponse
				LogData = &dto.CustomLoggerRequest{Remarks: SourceData.HeaderPath}
			)
			ctx := context.WithValue(c.Request().Context(), enum.CUSTOM_CONTEXT_VALUE, SourceData)
			ctx = context.WithValue(ctx, enum.CUSTOM_LOG_DATA, LogData)
			c.SetRequest(c.Request().WithContext(ctx))

			defer func() {
				svc.AccessLogger(c.Request().Context())
			}()
			if routh != nil {
				if !(*routh)[SourceData.HeaderMethod+SourceData.RequestPath] {
					return Next(c)
				}
			}

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				LogData.Error = err.Error()
				svc.logger.ErrorLogger(c.Request().Context(), "CommonCustomHeaderMiddleware2", err)
				return c.JSON(http.StatusBadRequest, dto.BaseResponse{
					StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
					Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
					Error:      err.Error(),
					Data:       nil,
				})
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

			//validation

			g := errgroup.Group{}
			g.Go(func() error {
				resTmp, errTMP := svc.ValidateMandatoryHeader(
					c.Request().Context(),
					SourceData, excludeUrl,
				)
				if errTMP != nil {
					errRes = resTmp
				}
				return errTMP
			})
			g.Go(func() error {
				resTmp, errTmp := svc.ValidateSignature(
					c.Request().Context(),
					SourceData,
					body,
					excludeUrl,
				)
				if errTmp != nil {
					errRes = resTmp
				}
				return errTmp
			})
			g.Go(func() error {
				resTmp, errTmp := svc.ValidateAuthorization(
					c.Request().Context(),
					SourceData,
					excludeUrl,
				)
				if errTmp != nil {
					errRes = resTmp
				}
				return errTmp
			})

			//finish validate
			err = g.Wait()
			if err != nil {
				if errRes == nil {
					LogData.Error = err.Error()
					switch errRes.StatusCode {
					case pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE:
						return c.JSON(http.StatusBadRequest, *errRes)
					case pkgErr.AUTH_INVALID_SIGNATURE_CODE:
						return c.JSON(http.StatusBadRequest, *errRes)
					}

				}
				svc.logger.ErrorLogger(c.Request().Context(), "CommonCustomHeaderMiddleware2.err", err)
				return c.JSON(http.StatusBadRequest, *errRes)
			}

			return Next(c)
		}
	}
}

func (svc *customMiddleware) AccessLogger(ctx context.Context) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	var errString string
	slog.SetDefault(logger)
	if ctx.Err() != nil {
		errString = ctx.Err().Error()
	}

	customResource, ok := ctx.Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Warn("failed to get custom resource")
	}
	customlogger, ok := ctx.Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !ok {
		log.Warn("failed to get custom logger")
	}
	logData := dto.HandlerLog{
		RequestId: customResource.HeaderRequestId,
		Timestamp: customResource.HeaderXTimestamp,
		Url:       customResource.HeaderHost + customResource.HeaderPath,
		Method:    customResource.HeaderMethod,
		Device: dto.Device{
			DeviceID:  customResource.HeaderXDeviceID,
			Longitude: customResource.HeaderXLongitude,
			Latitude:  customResource.HeaderXLatitude,
			Ip:        customResource.HeaderXRealIp,
		},
		Success:  false,
		Error:    errString,
		UserUUID: customResource.AuthUUID,
		Email:    customResource.AuthEmail,
	}
	if customlogger != nil {
		if logData.UserUUID == "" {
			logData.UserUUID = customlogger.UserUUID
		}
		if logData.Email == "" {
			logData.Email = customlogger.Email
		}
		if customlogger.Error != "" {
			logData.Error = customlogger.Error
		}
		if customlogger.Remarks != "" {
			logData.Remarks = customlogger.Remarks
		}
		logData.Data = customlogger.Data
		if customlogger.Success {
		}
		logData.Success = customlogger.Success
	}
	jsonLog, _ := json.Marshal(logData)
	slog.Info(string(jsonLog))

}

func (svc *customMiddleware) ValidateSignature(ctx context.Context, req *dto.ContextValue, reader []byte, excUrl *ExcludeURLValidation) (*dto.BaseResponse, error) {
	if excUrl.ValidationSignaure != nil {
		for _, value := range excUrl.ValidationSignaure {
			if value == req.RequestPath {
				return nil, nil
			}
		}
	}
	var (
		hashedRequestBody string
		err               error
	)
	switch req.HeaderContentType {
	case "application/json":
		if reader == nil {
			err = errors.New("application/json but payload is empty")
			svc.logger.ErrorLogger(ctx, "SignatureValidate.json.Compact", err)
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      err.Error(),
			}, err
		}
		if string(reader) == "" {
			err = errors.New("application/json but payload is empty")
			svc.logger.ErrorLogger(ctx, "SignatureValidate.json.Compact", err)
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
				Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
				Error:      "application/json but payload is empty",
			}, err
		}
		if len(reader) > 0 {
			dst := bytes.Buffer{}
			if err := json.Compact(&dst, reader); err != nil {
				log.Errorf("SignatureValidate.json.Compact: %v", err)
				return &dto.BaseResponse{
					StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
					Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
					Error:      err.Error(),
				}, err
			}
			hashedRequestBody = strings.ToLower(fmt.Sprintf("%x", sha256.Sum256(dst.Bytes())))
		}
	default:
		strData := fmt.Sprintf("%s:%s:%s", req.HeaderXTimestamp, req.HeaderXNonce, req.HeaderXDeviceID)
		hashedRequestBody = strings.ToLower(fmt.Sprintf("%x", sha256.Sum256([]byte(strData))))
	}

	if hashedRequestBody != req.HeaderXSignature {
		err = errors.New("signature verification failed")
		svc.logger.ErrorLogger(ctx, "SignatureValidate", err)
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_SIGNATURE_CODE,
			Message:    pkgErr.INVALID_SIGNATURE_MSG,
			Error:      err.Error(),
			Data:       nil,
		}, err
	}
	return nil, nil
}

func (svc *customMiddleware) ValidateXNonce(ctx context.Context, req *dto.ContextValue, excUrl ExcludeURLValidation) (*dto.BaseResponse, error) {
	for _, value := range excUrl.ValidationXNonce {
		if value == req.RequestPath {
			return nil, nil
		}
	}
	var xNonce = req.HeaderXNonce

	value, err := svc.Redis.GetXNonce(ctx, xNonce)
	if err != nil {
		log.Errorf("XNonceValidate.s.Redis.GetXSessionI err: %v", err)
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      "X-NONCE must be unix",
		}, err
	}
	if value == "active" {
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      "X-NONCE must be unix",
		}, errors.New("X-NONCE must be unix")
	}
	err = svc.Redis.SetXNONCE(ctx, xNonce)
	if err != nil {
		log.Errorf("XNonceValidate.s.Redis.SetExSessionId err: %v ", err)
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		}, err
	}
	return nil, nil
}

func (svc *customMiddleware) ValidateMandatoryHeader(ctx context.Context, req *dto.ContextValue, excUrl *ExcludeURLValidation) (*dto.BaseResponse, error) {
	for _, value := range excUrl.MandatoryHeader {
		if value == req.RequestPath {
			return nil, nil
		}
	}
	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, *req)
		log.Errorf("ValidateMandatoryHeader.err: %v", err)
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		}, err
	}
	return nil, nil
}
