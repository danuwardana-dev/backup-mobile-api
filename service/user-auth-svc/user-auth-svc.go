package userAuthSvc

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/middleware"
	"backend-mobile-api/internal/outbond/smtp"
	"backend-mobile-api/internal/outbond/verihubs"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/internal/repository/redis"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/dto/response"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	verihubsDto "backend-mobile-api/model/outbond/verihubs-dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	mathRand "math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type userAuthService struct {
	userRespository       postgres.UserRepository
	authService           middleware.CustomMiddleware
	redis                 redis.Redis
	rootConfig            config.Root
	smtp                  *smtp.Smtp
	clogger               *helpers.CustomLogger
	AccessStateRepository postgres.AccessStateRepository
	otpRepository         postgres.OtpRepository
	tokenBlacklist        postgres.TokenBlacklistTokenRepository
	outboundVeriHubsSvc   verihubs.OutboundVeriHubsService
	userDetailRepository  postgres.UserDetailRepository
	deviceRepository      postgres.DeviceRepository
}

func NewUserAuthService(
	userRespository postgres.UserRepository,
	authService middleware.CustomMiddleware,
	redis redis.Redis,
	rootConfig config.Root,
	smtp *smtp.Smtp,
	clogger *helpers.CustomLogger,
	AccessStateRepository postgres.AccessStateRepository,
	otpRepository postgres.OtpRepository,
	tokenBlacklist postgres.TokenBlacklistTokenRepository,
	//	loginLogRepository postgres.LoginLogRepository,
	outboundVeriHubsSvc verihubs.OutboundVeriHubsService,
	userDetilRepository postgres.UserDetailRepository,
	deviceRepository postgres.DeviceRepository,
) UserAuthService {
	return &userAuthService{
		userRespository:       userRespository,
		authService:           authService,
		redis:                 redis,
		rootConfig:            rootConfig,
		smtp:                  smtp,
		clogger:               clogger,
		AccessStateRepository: AccessStateRepository,
		otpRepository:         otpRepository,
		tokenBlacklist:        tokenBlacklist,
		//		loginLogRepository:   loginLogRepository,
		outboundVeriHubsSvc:  outboundVeriHubsSvc,
		userDetailRepository: userDetilRepository,
		deviceRepository:     deviceRepository,
	}
}

type UserAuthService interface {
	RegisterService(c context.Context, req *request.RegisterRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	sendOtp(c context.Context, req *dto.SendOtp) (*entity.OTP, error)
	LoginByEmailService(C context.Context, req *request.LoginEmailRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	LoginByPhoneNumberService(C context.Context, req *request.LoginPhoneNumberRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	LogoutService(c context.Context, req *request.LogoutRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	RefreshService(c context.Context, req request.RefreshTokenRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	VerifyOtpService(c context.Context, otpRequest *request.VerifyOtpRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	GenerateOtp(c context.Context, value string, uuidKey string) (string, error)
	SendOtpService(c context.Context, req *request.SendOtpRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	SetPinService(c context.Context, req *request.SetPinRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	ForgotPinService(ctx context.Context, req *request.ForgotPinRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
}

func (svc *userAuthService) RegisterService(c context.Context, req *request.RegisterRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		userByemail *[]entity.User
		userByphone *[]entity.User
	)
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		if userByemail, err = svc.userRespository.SelectUserByStruct(c, &entity.User{Email: req.Email}); err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		if userByphone, err = svc.userRespository.SelectUserByStruct(c, &entity.User{PhoneNumber: req.PhoneNumber}); err != nil {
			return err
		}
		return nil
	})
	err := eg.Wait()
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
				Message:    pkgErr.USER_NOT_FOUND_MSG,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if userByemail != nil {
		if len(*userByemail) > 0 {
			logData.Error = pkgErr.EMAIL_ALREADY_REGISTERED_MSG
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUHT_EMAIL_ALREADY_REGISTERED_CODE,
				Message:    pkgErr.EMAIL_ALREADY_REGISTERED_MSG,
			}
		}
	}
	if userByphone != nil {
		if len(*userByphone) > 0 {
			logData.Error = pkgErr.PHONE_NUBMBER_ALREADY_REGISTERED_MSG
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_PHONE_NUMBER_ALREADY_REGISTERED_CODE,
				Message:    pkgErr.PHONE_NUBMBER_ALREADY_REGISTERED_MSG,
			}
		}
	}
	tx := svc.userRespository.Tx(c)
	err = svc.userRespository.InsertUser(
		c,
		tx,
		&entity.User{
			UUID:        uuid.New().String(),
			FullName:    req.FullName,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Status:      enum.VERIFICATION_STATUS_UNVERIFIED,
		})

	if err != nil {
		tx.Rollback()
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	tx.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}

func (svc *userAuthService) LoginByEmailService(c context.Context, req *request.LoginEmailRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		user *entity.User
		err  error
	)
	user, err = svc.userRespository.SelectUserByEmailOrPhoneNumber(c, req.Email)
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
				Message:    pkgErr.USER_NOT_FOUND_MSG,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	if user.Status == enum.VERIFICATION_STATUS_UNVERIFIED {
		logData.Error = "user is not verified"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNVERIFIED_CODE,
			Message:    pkgErr.UNVERIFIED_MSG,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin))
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.WRONG_EMAIL_OR_PIN_MSG,
		}
	}
	token, err := svc.authService.CreateTokens(c, &middleware.Claims{
		Uuid:     user.UUID,
		Username: user.Email,
	})
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data: response.LoginResponse{
			User: response.UserData{
				UUID:  user.UUID,
				Name:  user.FullName,
				Email: user.Email,
				Phone: user.PhoneNumber,
			},
			Token: token,
		},
	}
}
func (svc *userAuthService) LoginByPhoneNumberService(c context.Context, req *request.LoginPhoneNumberRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		user *entity.User
		err  error
	)
	user, err = svc.userRespository.SelectUserByEmailOrPhoneNumber(c, req.PhoneNumber)
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
				Message:    pkgErr.USER_NOT_FOUND_MSG,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.Email = user.Email
	logData.UserUUID = user.UUID
	if user.DeviceID != req.DeviceID {
		logData.Error = "device id mismatch"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.UNAUTHORIZED_MSG,
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin))
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.WRONG_PHONE_NUMBER_OR_PIN_MSG,
		}
	}
	token, err := svc.authService.CreateTokens(c, &middleware.Claims{
		Uuid:     user.UUID,
		Username: user.Email,
	})
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data: response.LoginResponse{
			User: response.UserData{
				UUID:  user.UUID,
				Name:  user.FullName,
				Email: user.Email,
				Phone: user.PhoneNumber,
			},
			Token: token,
		},
	}
}

func (svc *userAuthService) RefreshService(c context.Context, req request.RefreshTokenRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	isExist, err := svc.tokenBlacklist.IsBlaclistTokenActive(c, req.RefreshToken)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if isExist {
		logData.Error = "unauthorized"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.UNAUTHORIZED_MSG,
		}
	}

	users, err := svc.userRespository.SelectUserByStruct(c, &entity.User{UUID: req.UUID})
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
				Message:    pkgErr.USER_NOT_FOUND_MSG,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if users == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	if len(*users) < 1 {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	userData := (*users)[0]
	logData.UserUUID = userData.UUID
	logData.Email = userData.Email

	if userData.DeviceID != req.DeviceID {
		logData.Error = "deference device"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}

	ok, tokenData, err := svc.authService.RefreshToken(c, req.RefreshToken, &userData)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if !ok {
		logData.Error = "unauthorized"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.UNAUTHORIZED_MSG,
		}
	}
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data: response.LoginResponse{
			User: response.UserData{
				UUID:  userData.UUID,
				Name:  userData.FullName,
				Email: userData.Email,
				Phone: userData.PhoneNumber,
			},
			Token: *tokenData,
		},
	}
}
func (svc *userAuthService) LogoutService(c context.Context, req *request.LogoutRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	customResource, ok := c.Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		logData.Error = "failed to get custom resource"
		log.Error("failed to get custom resource")
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		}
	}
	tokenString := customResource.HeaderAuthorization
	if tokenString == "" {
		err := errors.New("token is empty")
		logData.Error = err.Error()
		svc.clogger.ErrorLogger(c, "Logout.tokenString.empty", err)
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.UNAUTHORIZED_MSG,
		}
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(tokenString, prefix) {
		logData.Error = "no Bearer prefix on authorization header"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.UNAUTHORIZED_MSG,
		}
	}
	tokenString = strings.TrimPrefix(tokenString, prefix)

	if tokenString == "" {
		logData.Error = "unauthorized"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.UNAUTHORIZED_MSG}
	}
	var (
		g             = errgroup.Group{}
		accessExpire  int64
		refreshExpire int64
	)
	g.Go(func() error {
		rsa, err := svc.authService.EncodePublicKeyRSA(c, svc.rootConfig.Jwt.PublicKey)
		if err != nil {
			return err
		}
		token, err := svc.authService.ParseJwtToken(c, fmt.Sprint(tokenString), rsa)
		if err != nil {
			return err
		}
		mapAccessClaim, err := svc.authService.ClaimJWT(c, token)
		if err != nil {
			return err
		}
		accessExpire = int64((*mapAccessClaim)["exp"].(float64))
		return nil
	})
	g.Go(func() error {
		rsa, err := svc.authService.EncodePublicKeyRSA(c, svc.rootConfig.Jwt.RefreshPublicKey)
		if err != nil {
			return err
		}
		token, err := svc.authService.ParseJwtToken(c, req.RefreshToken, rsa)
		if err != nil {
			return err
		}
		mapRefreshClaim, err := svc.authService.ClaimJWT(c, token)
		if err != nil {
			return err
		}
		refreshExpire = int64((*mapRefreshClaim)["exp"].(float64))
		return nil
	})
	if err := g.Wait(); err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}

	//set-up blacklist token
	now := time.Now()
	nowUnix := now.Unix()
	if accessExpire-nowUnix > 0 {
		_ = svc.redis.SetBlaclistJwt(c, fmt.Sprint(tokenString), time.Duration(accessExpire-nowUnix)*time.Second)
	}
	tx := svc.tokenBlacklist.Tx(c)
	err := svc.tokenBlacklist.InsertBlaclistToken(c, tx, &entity.TokenBlacklist{
		Token:       req.RefreshToken,
		BlacklistAt: now,
		ExpiredAt:   now.Add(time.Duration(refreshExpire)),
		Description: "refresh-token log-out",
	})
	if err != nil {
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	tx.Commit()

	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}
func (svc *userAuthService) VerifyOtpService(c context.Context, req *request.VerifyOtpRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		otpData       *entity.OTP
		otpUpdate     *entity.OTP
		userUpdate    *entity.User
		userData      entity.User
		newUserDetail *entity.UserDetail
		newDevice     *entity.Device
		resetPin      *entity.AccessState
		newUUID       = uuid.New().String()
	)
	strJson, err := svc.redis.GetOtp(c, req.Otp, req.VerifyID)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if strJson != "" {
		var otpTmp entity.OTP
		err = json.Unmarshal([]byte(strJson), &otpTmp)
		if err != nil {
			logData.Error = err.Error()
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}
		otpData = &otpTmp
		_ = svc.redis.DeleteOtp(c, req.Otp, req.VerifyID)
	} else {
		otpData, err = svc.otpRepository.SelectOtpByVerifyKey(c, req.VerifyID)
		if err != nil {
			logData.Error = err.Error()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &dto.BaseResponse{
					StatusCode: pkgErr.AUTH_INVALID_OTP_CODE,
					Message:    pkgErr.INVALID_OTP_MSG,
				}
			}
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}
		if otpData.Status == enum.OTP_BLOCKED || otpData.ExpiredAt.Before(time.Now()) {
			logData.Error = "blocked or expired"
			return &dto.BaseResponse{
				StatusCode: pkgErr.RECORD_NOT_FOUND_MSG,
				Message:    pkgErr.RECORD_NOT_FOUND_MSG,
			}
		}
	}
	if otpData.OtpCode != req.Otp {
		logData.Error = "invalid otp"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_OTP_CODE,
			Message:    pkgErr.INVALID_OTP_MSG,
		}
	}
	users, err := svc.userRespository.SelectUserByStruct(c, &entity.User{ID: otpData.UserId})
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNVERIFIED_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if users == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	if len(*users) < 1 {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	userData = (*users)[0]
	logData.UserUUID = userData.UUID
	logData.Email = userData.Email
	switch otpData.OtpMethod {
	case enum.TYPE_EMAIL:
		otpUpdate = &entity.OTP{Status: enum.OTP_VERIFIED}
	case enum.TYPE_SMS:
		if _, err := svc.outboundVeriHubsSvc.VerifySMSOtpService(c, &verihubsDto.VerifyOtpBaseRequest{
			MSISDN:    otpData.OtpDestination,
			Otp:       otpData.OtpCode,
			Challenge: nil,
		}); err != nil {
			logData.Error = err.Error()
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}
		otpUpdate = &entity.OTP{Status: enum.OTP_VERIFIED}
	case enum.TYPE_WHATSAPP:
		if _, err = svc.outboundVeriHubsSvc.VerifyWhatsappsOtpService(c, &verihubsDto.VerifyOtpBaseRequest{
			MSISDN:    otpData.OtpDestination,
			Otp:       otpData.OtpCode,
			Challenge: nil,
		}); err != nil {
			logData.Error = err.Error()
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}
	default:

	}
	if otpData.OtpPurpose == enum.OTP_VERIFY_ACCOUNT {

		userUpdate = &entity.User{Status: enum.VERIFICATION_STATUS_VERIFIED}
		newUserDetail = &entity.UserDetail{
			UserId:    userData.ID,
			UserUUID:  userData.UUID,
			Biometric: enum.BIOMETRIC_IN_ACTIVE,
			Status:    enum.USER_WAITING_KYC_PROCESS,
			KycStatus: enum.USER_KYC_STATUS_UNKNOWN,
		}
	}
	if otpData.OtpPurpose == enum.OTP_FORGOT_PIN || otpData.OtpPurpose == enum.OTP_VERIFY_ACCOUNT {
		resetPin = &entity.AccessState{
			AccessType:  enum.ACCESS_SET_PIN,
			UserId:      userData.ID,
			UserUUID:    userData.UUID,
			DeviceId:    req.DeviceID,
			AccessToken: newUUID,
			ExpiredAt:   time.Now().Add(svc.rootConfig.App.AccessKeyExpire),
			Used:        false,
		}
	}

	if req.DeviceID != userData.DeviceID {
		if userUpdate == nil {
			userUpdate = &entity.User{}
		}
		userUpdate.DeviceID = req.DeviceID

		newDevice = &entity.Device{
			UserID:         uint(userData.ID),
			UserUUID:       userData.UUID,
			DeviceID:       req.DeviceID,
			AppVersionCode: req.DeviceInfo.AppVersionCode,
			AppVersionName: req.DeviceInfo.AppVersionName,
			Manufacturer:   req.DeviceInfo.Manufacturer,
			Brand:          req.DeviceInfo.Brand,
			DeviceModel:    req.DeviceInfo.Model,
			Product:        req.DeviceInfo.Product,
			VersionSdk:     req.DeviceInfo.VersionSdk,
			VersionRelease: req.DeviceInfo.VersionRelease,
		}
	}

	//tx update
	txOtp := svc.otpRepository.Tx(c)
	txUser := svc.userRespository.Tx(c)
	txUserDt := svc.userDetailRepository.Tx(c)
	txDevice := svc.deviceRepository.Tx(c)
	txResetPin := svc.AccessStateRepository.Tx(c)

	g := errgroup.Group{}
	//otp
	g.Go(func() error {
		var errData error
		if otpUpdate != nil {
			errData = svc.otpRepository.UpdateOtpDataRepository(c, txOtp, otpData, otpUpdate)
		}
		return errData
	})
	//user
	g.Go(func() error {
		var errData error
		if userUpdate != nil {
			errData = svc.userRespository.UpdateUser(c, txUser, &userData, userUpdate)
		}
		return errData
	})
	//userDetail
	g.Go(func() error {
		var errData error
		if newUserDetail != nil {
			errData = svc.userDetailRepository.InsertUserDetail(c, txUserDt, newUserDetail)
		}
		return errData
	})
	g.Go(func() error {
		var errData error
		if newDevice != nil {
			errData = svc.deviceRepository.InsertDevice(c, txDevice, newDevice)
		}
		return errData
	})
	g.Go(func() error {
		var errData error
		if resetPin != nil {
			errData = svc.AccessStateRepository.InsertAccessStateRepository(c, txResetPin, resetPin)
		}
		return errData
	})

	err = g.Wait()
	if err != nil {
		txOtp.Rollback()
		txUser.Rollback()
		txUserDt.Rollback()
		txDevice.Rollback()
		txResetPin.Rollback()

		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}

	}
	txOtp.Commit()
	txUser.Commit()
	txUserDt.Commit()
	txDevice.Commit()
	txResetPin.Commit()

	jsonUser, _ := json.Marshal(userData)
	expireAcc := time.Now().Add(svc.rootConfig.App.AccessKeyExpire)
	err = svc.redis.SetAccessKey(c, newUUID, string(jsonUser), svc.rootConfig.App.AccessKeyExpire)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       response.VerifyOtp{AccessKey: newUUID, ExpireAt: expireAcc},
	}
}
func (svc *userAuthService) GenerateOtp(c context.Context, value string, uuidKey string) (string, error) {
	var (
		exist = false
		err   error
	)
	mathRand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", mathRand.Intn(1000000))

	if exist, err = svc.redis.OtpIsExist(c, otp, uuidKey); err != nil {
		svc.clogger.ErrorLogger(c, "GenerateOtp.svc.redis.OtpIsExist", err)
		return "", err
	}
	if exist {
		return svc.GenerateOtp(c, value, uuidKey)
	}
	return otp, err
}
func (svc *userAuthService) SendOtpService(ctx context.Context, req *request.SendOtpRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	logData.Remarks = fmt.Sprintf("send-otp:%s:%s", req.OtpMethod, req.ServiceType)
	user, err := svc.userRespository.SelectUserByEmailOrPhoneNumber(ctx, req.VerifyTo)
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
				Message:    pkgErr.USER_NOT_FOUND_MSG,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	aksesKey := uuid.New().String()
	wrapContext := helpers.WrapContext(ctx)
	switch req.ServiceType {
	case enum.OTP_VERIFY_ACCOUNT:
		if user.Status != enum.VERIFICATION_STATUS_UNVERIFIED {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_ALREADY_VERIFIED_CODE,
				Message:    pkgErr.ALREADY_VERIFIED_MSG,
			}
		}

		otpData, err := svc.sendOtp(wrapContext, &dto.SendOtp{
			OtpPurpose:     req.ServiceType,
			OtpMethod:      req.OtpMethod,
			OtpDestination: req.VerifyTo,
			UserId:         user.ID,
			UserUUID:       user.UUID,
			VerifyKey:      aksesKey,
		})
		if err != nil {
			logData.Error = err.Error()
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}
		jsonOtp, _ := json.Marshal(otpData)
		if err = svc.redis.SetOtp(ctx, otpData.OtpCode, aksesKey, string(jsonOtp), svc.rootConfig.App.OtpExpire); err != nil {
			logData.Error = err.Error()
			svc.clogger.ErrorLogger(ctx, "GenerateOtp.redis.SetOtp", err)
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}
	case enum.OTP_FORGOT_PIN:
		if user.Status != enum.VERIFICATION_STATUS_VERIFIED {
			logData.Error = "user not verified"
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_UNVERIFIED_CODE,
				Message:    pkgErr.UNVERIFIED_MSG,
			}
		}

		otpData, err := svc.sendOtp(ctx, &dto.SendOtp{
			OtpPurpose:     req.ServiceType,
			OtpMethod:      req.OtpMethod,
			OtpDestination: req.VerifyTo,
			UserId:         user.ID,
			UserUUID:       user.UUID,
			VerifyKey:      aksesKey,
		})
		if err != nil {
			logData.Error = err.Error()
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}
		jsonOtp, _ := json.Marshal(otpData)
		if err = svc.redis.SetOtp(ctx, otpData.OtpCode, aksesKey, string(jsonOtp), svc.rootConfig.App.OtpExpire); err != nil {
			svc.clogger.ErrorLogger(ctx, "SendOtpService.ForgotPin.redis.SetOtp", err)
			logData.Error = err.Error()
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}
		}

	default:
		logData.Error = "invalid otp service_type"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      "invalid otp service_type",
		}
	}
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data: response.SendOtpResponse{
			VerifyId: aksesKey,
			ExpireAt: time.Now().Add(svc.rootConfig.App.OtpExpire),
		},
	}
}
func (svc *userAuthService) sendOtp(ctx context.Context, req *dto.SendOtp) (*entity.OTP, error) {
	var (
		jsonData, _ = json.Marshal(req)
		otpData     = entity.OTP{
			OtpCode:        "",
			OtpPurpose:     req.OtpPurpose,
			OtpMethod:      req.OtpMethod,
			OtpDestination: req.OtpDestination,
			UserId:         req.UserId,
			UserUUID:       req.UserUUID,
			VerifyKey:      req.VerifyKey,
			ExpiredAt:      time.Now().Add(svc.rootConfig.App.OtpExpire),
			SessionId:      "",
			Status:         "",
		}
	)

	switch req.OtpMethod {
	case enum.TYPE_EMAIL:
		otpData.Status = enum.OTP_DELIVERED
		otp, err := svc.GenerateOtp(ctx, string(jsonData), req.VerifyKey)
		if err != nil {
			return nil, err
		}
		if err := svc.smtp.SendMail(ctx, []string{req.OtpDestination}, enum.VERIFY_OTP_SUBJECT, svc.smtp.RegisterOtpMsg(otp)); err != nil {
			svc.clogger.ErrorLogger(ctx, "sendOtp.smtp.SendMail", err)
			otpData.Status = enum.OTP_FAILED
		}
		otpData.OtpCode = otp
		otpData.Status = enum.OTP_DELIVERED
		otpData.SessionId = req.VerifyKey
	case enum.TYPE_SMS:
		var otp *string
		if !svc.rootConfig.Verihubs.VerihubsOTPCode {
			otpStr, err := svc.GenerateOtp(ctx, string(jsonData), req.VerifyKey)
			if err != nil {
				return nil, err
			}
			otp = &otpStr
		}
		resp, err := svc.outboundVeriHubsSvc.SendSMSOtpService(ctx, &verihubsDto.SendOtpBaseRequest{
			MSISDN:      req.OtpDestination,
			Otp:         otp,
			Template:    svc.rootConfig.Verihubs.OTPSMSTemplate,
			TimeLimit:   int64(svc.rootConfig.App.OtpExpire.Seconds()),
			Challenge:   svc.rootConfig.Verihubs.OTPSMSChallenge,
			CallbackUrl: svc.rootConfig.Verihubs.OTPCallBackUrl,
		})
		if err != nil {
			otpData.Status = enum.OTP_FAILED
			svc.clogger.ErrorLogger(ctx, "sendOtp.outboundVeriHubsSvc.SendSMSOtpService", err)
			return nil, err
		}
		otpData.OtpCode = resp.Otp
		otpData.SessionId = resp.SessionID
		otpData.Status = enum.OTP_DELIVERED
	case enum.TYPE_WHATSAPP:
		var otp *string
		if !svc.rootConfig.Verihubs.VerihubsOTPCode {
			otpStr, err := svc.GenerateOtp(ctx, string(jsonData), req.VerifyKey)
			if err != nil {
				return nil, err
			}
			otp = &otpStr
		}
		resp, err := svc.outboundVeriHubsSvc.SendWhatsappsService(ctx, &verihubsDto.SendWhatsappOtpBaseRequest{
			MSISDN:       req.OtpDestination,
			Otp:          otp,
			Challenge:    nil,
			TimeLimit:    int64(svc.rootConfig.App.OtpExpire.Seconds()),
			LangCode:     svc.rootConfig.Verihubs.OTPWhatsappLangCode,
			TemplateName: svc.rootConfig.Verihubs.OTPWhatsappTemplateName,
			OtpLength:    svc.rootConfig.Verihubs.OTPWhatsappOtpLength,
			CallbackUrl:  svc.rootConfig.Verihubs.OTPCallBackUrl,
		})
		if err != nil {
			otpData.Status = enum.OTP_FAILED
		}
		otpData.OtpCode = resp.Otp
		otpData.SessionId = resp.SessionID
		otpData.Status = enum.OTP_DELIVERED
		otpData.UserUUID = req.UserUUID
	default:

	}
	tx := svc.otpRepository.Tx(ctx)
	err := svc.otpRepository.InsertOtpDataRepository(ctx, tx, &otpData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return svc.otpRepository.SelectOtpByVerifyKeyBeforeExpire(ctx, req.VerifyKey)
}
func (svc *userAuthService) SetPinService(c context.Context, req *request.SetPinRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		pinData  *entity.AccessState
		userData *entity.User
		err      error
	)
	if req.Pin != req.ConfirmedPin {
		logData.Error = "miss match confirmed pin"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      "deference confirmed pin",
		}
	}
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		pinData, errTmp = svc.AccessStateRepository.SelectAccessByAccessTokenRepository(c, req.AccountToken)
		return errTmp
	})
	g.Go(func() error {
		var errTmp error
		userJson, errTmp := svc.redis.GetAccessKey(c, req.AccountToken)
		if errTmp != nil {
			return errTmp
		}
		if userJson == "" {
			userData, errTmp = svc.userRespository.SelectUserByEmailOrPhoneNumber(c, req.EmailOrPhoneNumber)
			return errTmp
		}
		_ = svc.redis.DeleteAccessKey(c, req.AccountToken)
		var user entity.User
		errTmp = json.Unmarshal([]byte(userJson), &user)
		if errTmp != nil {
			return errTmp
		}
		userData = &user
		return nil
	})
	if err := g.Wait(); err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if userData == nil {
		logData.Error = "user not exist"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = userData.UUID
	logData.Email = userData.Email
	if pinData == nil {
		svc.clogger.ErrorLogger(c, "SetPinService.SelectResetPinByToken", errors.New("record not found"))
		logData.Error = "record pin reset not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}

	tx := svc.userRespository.Tx(c)
	hashPin, _ := bcrypt.GenerateFromPassword([]byte(req.Pin), bcrypt.DefaultCost)
	err = svc.userRespository.UpdateUser(c, tx, userData, &entity.User{Pin: string(hashPin)})
	if err != nil {
		tx.Rollback()
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	tx.Commit()
	tx = svc.AccessStateRepository.Tx(c)
	_ = svc.AccessStateRepository.UpdateAccessByStruct(c, tx, pinData, &entity.AccessState{
		Used: true,
	})
	tx.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}

func (svc *userAuthService) ForgotPinService(ctx context.Context, req *request.ForgotPinRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		user        *entity.User
		updaterUser *entity.User
		newDevice   *entity.Device
	)
	user, err := svc.userRespository.SelectUserByEmailOrPhoneNumber(ctx, req.Email)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	if user.DeviceID != req.DeviceID {
		updaterUser = &entity.User{DeviceID: req.DeviceID}
		newDevice = &entity.Device{
			UserID:         uint(user.ID),
			DeviceID:       req.DeviceID,
			AppVersionCode: req.DeviceInfo.AppVersionCode,
			AppVersionName: req.DeviceInfo.AppVersionName,
			Manufacturer:   req.DeviceInfo.Manufacturer,
			Brand:          req.DeviceInfo.Brand,
			DeviceModel:    req.DeviceInfo.Model,
			Product:        req.DeviceInfo.Product,
			VersionSdk:     req.DeviceInfo.VersionSdk,
			VersionRelease: req.DeviceInfo.VersionRelease,
		}
	}
	user.DeviceID = req.DeviceID
	jsonUser, _ := json.Marshal(user)

	txUser := svc.userRespository.Tx(ctx)
	txDevice := svc.userRespository.Tx(ctx)
	txPin := svc.AccessStateRepository.Tx(ctx)
	g, _ := errgroup.WithContext(ctx)
	accessKey := uuid.New().String()
	g.Go(func() error {
		if updaterUser != nil {
			return svc.userRespository.UpdateUser(ctx, txUser, user, updaterUser)
		}
		return nil
	})
	g.Go(func() error {
		if newDevice != nil {
			return svc.deviceRepository.InsertDevice(ctx, txDevice, newDevice)
		}
		return nil
	})
	g.Go(func() error {
		return svc.redis.SetAccessKey(ctx, accessKey, string(jsonUser), svc.rootConfig.App.AccessKeyExpire)
	})
	g.Go(func() error {
		return svc.AccessStateRepository.InsertAccessStateRepository(ctx, txPin, &entity.AccessState{
			AccessType:  enum.ACCESS_FORGOT_PIN,
			UserId:      user.ID,
			UserUUID:    user.UUID,
			DeviceId:    req.DeviceID,
			AccessToken: accessKey,
			ExpiredAt:   time.Now().Add(svc.rootConfig.App.AccessKeyExpire),
			Used:        false,
		})
	})
	err = g.Wait()
	if err != nil {
		txUser.Rollback()
		txDevice.Rollback()
		txPin.Rollback()
		logData.Error = err.Error()

		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	txUser.Commit()
	txDevice.Commit()
	txPin.Commit()
	url := fmt.Sprintf("%s%s", req.CallbackUrl, fmt.Sprintf("?acces_key=%s", accessKey))
	svc.smtp.SendMail(ctx, []string{user.Email}, enum.ACCESS_RESET_PIN_SUBJECT, url)
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}
