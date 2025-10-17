package userProfileService

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/outbond/smtp"
	"backend-mobile-api/internal/outbond/verihubs"
	"backend-mobile-api/internal/repository/minio"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/internal/repository/redis"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/dto/response"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	"backend-mobile-api/service/otp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type userProfileService struct {
	userRepository        postgres.UserRepository
	redis                 *redis.Redis
	rootConfig            config.Root
	smtp                  *smtp.Smtp
	clogger               *helpers.CustomLogger
	AccessStateRepository postgres.AccessStateRepository
	otpRepository         postgres.OtpRepository
	outboundVeriHubsSvc   verihubs.OutboundVeriHubsService
	userDetailRepository  postgres.UserDetailRepository
	otpService            otp.OtpService
	minioRepository       minio.MinioRepository
}

func NewUserProfileService(
	userRepository postgres.UserRepository,
	redis *redis.Redis,
	rootConfig config.Root,
	smtp *smtp.Smtp,
	clogger *helpers.CustomLogger,
	AccessStateRepository postgres.AccessStateRepository,
	otpRepository postgres.OtpRepository,
	outboundVeriHubsSvc verihubs.OutboundVeriHubsService,
	userDetailRepository postgres.UserDetailRepository,
	otpService otp.OtpService,
	minioRepository minio.MinioRepository) UserProfileService {
	return &userProfileService{
		userRepository:        userRepository,
		redis:                 redis,
		rootConfig:            rootConfig,
		smtp:                  smtp,
		clogger:               clogger,
		AccessStateRepository: AccessStateRepository,
		otpRepository:         otpRepository,
		outboundVeriHubsSvc:   outboundVeriHubsSvc,
		userDetailRepository:  userDetailRepository,
		otpService:            otpService,
		minioRepository:       minioRepository,
	}
}

type UserProfileService interface {
	AccessTokenService(ctx context.Context, req *request.AccessTokenRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	InquiryUserProfileService(ctx context.Context, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	ResetPinService(ctx context.Context, req *request.ResetPinRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	sendOtpUpdateProfile(ctx context.Context, reg *entity.OTP, profile *dto.ProfileUpdateRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	ResetEmailService(ctx context.Context, req *request.ResetEmailRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	ResetPhoneNumber(ctx context.Context, req *request.ResetPhoneNumberRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	VerifyOtpService(ctx context.Context, req *request.VerifyOtpRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	ResetFullNameService(ctx context.Context, req *request.ResetFullNameRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	DeleteAccountService(ctx context.Context, req *request.DeletetAccountRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	BiometricStatusService(ctx context.Context, req *request.BiometrictStatusRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	ResetProfilePictureService(ctx context.Context, req *request.ResetProfilePictureRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	GetProfilePictureController(ctx context.Context, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
}

func (svc *userProfileService) AccessTokenService(ctx context.Context, req *request.AccessTokenRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	customResource, ok := ctx.Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		err := errors.New("failed to get custom resource")
		svc.clogger.ErrorLogger(ctx, "AccessTokenService", err)
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		}
	}
	user, err := svc.userRepository.SelectUserByUUID(ctx, *userUUID)
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
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
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}

	logData.UserUUID = user.UUID
	logData.Email = user.Email

	if user.DeviceID != req.DeviceID || user.DeviceID != customResource.HeaderXDeviceID {
		logData.Error = "deference device found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin))
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
			Message:    pkgErr.INVALID_PIN,
		}
	}
	accessKey := uuid.New().String()
	expired := time.Now().Add(svc.rootConfig.App.AccessKeyExpire)
	txAccess := svc.AccessStateRepository.Tx(ctx)

	g := errgroup.Group{}
	g.Go(func() error {
		return svc.AccessStateRepository.InsertAccessStateRepository(ctx, txAccess, &entity.AccessState{
			AccessType:  req.AccessType,
			UserId:      user.ID,
			UserUUID:    user.UUID,
			DeviceId:    user.DeviceID,
			AccessToken: accessKey,
			ExpiredAt:   expired,
			Used:        false,
		})
	})
	g.Go(func() error {
		var jsonUser, _ = json.Marshal(user)

		return svc.redis.SetAccessKey(ctx, accessKey, string(jsonUser), svc.rootConfig.App.AccessKeyExpire)
	})

	if err = g.Wait(); err != nil {
		txAccess.Rollback()
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	txAccess.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data: response.AccessTokenResponse{
			AccessKey: accessKey,
			ExpireAt:  expired,
		},
	}
}
func (svc *userProfileService) sendOtpUpdateProfile(ctx context.Context, req *entity.OTP, newProfileRequest *dto.ProfileUpdateRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		otpKey  = uuid.New().String()
		err     error
		otpData *entity.OTP
	)
	otpData, err = svc.otpService.SendOtp(ctx, &dto.SendOtp{
		OtpPurpose:     req.OtpPurpose,
		OtpMethod:      req.OtpMethod,
		OtpDestination: req.OtpDestination,
		UserId:         req.UserId,
		UserUUID:       req.UserUUID,
		VerifyKey:      otpKey,
	})
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
		}
	}
	JsonOtp, _ := json.Marshal(otpData)
	err = svc.redis.SetOtp(ctx, otpData.OtpCode, otpData.VerifyKey, string(JsonOtp), svc.rootConfig.App.OtpExpire)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{StatusCode: pkgErr.UNDEFINED_ERROR_CODE, Message: pkgErr.SERVER_BUSY}
	}
	err = svc.redis.SetDataUpdateUserProfile(ctx, otpKey, newProfileRequest, svc.rootConfig.App.AccessKeyExpire)
	if err != nil {
		svc.clogger.ErrorLogger(ctx, "SetDataUpdateUserProfile", err)
		logData.Error = err.Error()
		return &dto.BaseResponse{StatusCode: pkgErr.UNDEFINED_ERROR_CODE, Message: pkgErr.SERVER_BUSY}
	}
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data: response.VerifyOtp{
			AccessKey: otpData.VerifyKey,
			ExpireAt:  otpData.ExpiredAt,
		},
	}
}

func (svc *userProfileService) InquiryUserProfileService(ctx context.Context, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		user   *entity.User
		userDt *entity.UserDetail
	)
	customResource, ok := ctx.Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		err := errors.New("failed to get custom resource")
		svc.clogger.ErrorLogger(ctx, "AccessTokenService", err)
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		}
	}
	logData.Remarks = "get-profile-picture"
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil
	})
	g.Go(func() error {
		var errTmp error
		userDt, errTmp = svc.userDetailRepository.SelectUserDetailByUserUUID(ctx, userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errTmp
	})
	if err := g.Wait(); err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}

	if userDt == nil || user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	if customResource.HeaderXDeviceID != user.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	logData.Success = true

	resData := response.UserProfileInquiryResponse{
		Fullname:        user.FullName,
		Email:           user.Email,
		PhoneNumber:     user.PhoneNumber,
		Status:          userDt.Status,
		BiometricStatus: userDt.Biometric,
	}
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       resData,
	}

}

func (svc *userProfileService) ResetPinService(ctx context.Context, req *request.ResetPinRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		accessData *entity.AccessState
		user       *entity.User
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		accessData, errTmp = svc.AccessStateRepository.SelectAccessByAccessTokenRepository(ctx, req.AccountToken)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil
	})
	g.Go(func() error {
		strUserJson, errTmp := svc.redis.GetAccessKey(ctx, req.AccountToken)
		if errTmp != nil {
			return errTmp
		}
		if strUserJson != "" {
			errTmp = json.Unmarshal([]byte(strUserJson), &user)
			if errTmp != nil {
				return errTmp
			}
			return nil
		}
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil

	})
	err := g.Wait()
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if user == nil {
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	if accessData == nil {
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessData.AccessType != enum.ACCESS_RESET_PIN {
		logData.Error = "invalid access type"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessData.ExpiredAt.Before(time.Now()) {
		logData.Error = "access expired"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_EXPIRED_ACCESS_CODE,
			Message:    pkgErr.EXPIRED_TIME_MSG,
		}
	}
	if req.DeviceID != user.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin))
	if err == nil {
		logData.Error = "is existing pin, no update"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_IS_EXISTING_PIN,
			Message:    pkgErr.IS_EXISTING_PIN,
		}
	}

	return svc.sendOtpUpdateProfile(ctx, &entity.OTP{
		OtpPurpose:     enum.OTP_RESET_PIN,
		OtpMethod:      enum.TYPE_SMS,
		OtpDestination: user.PhoneNumber,
		UserId:         user.ID,
		UserUUID:       user.UUID,
	}, &dto.ProfileUpdateRequest{
		Field: enum.PROFILE_UPDATE_PIN,
		Value: req.Pin,
		User:  user,
	}, logData)

}

func (svc *userProfileService) ResetEmailService(ctx context.Context, req *request.ResetEmailRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		err         error
		userIsExist *entity.User
		user        *entity.User
		accessData  *entity.AccessState
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		accessData, errTmp = svc.AccessStateRepository.SelectAccessByAccessTokenRepository(ctx, req.AccountToken)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil
	})

	g.Go(func() error {
		strUserJson, errTmp := svc.redis.GetAccessKey(ctx, req.AccountToken)
		if errTmp != nil {
			return errTmp
		}
		if strUserJson != "" {
			errTmp = json.Unmarshal([]byte(strUserJson), &user)
			if errTmp != nil {
				return errTmp
			}
			return nil
		}
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil

	})

	g.Go(func() error {
		var errTmp error
		userIsExist, errTmp = svc.userRepository.SelectUserByEmailOrPhoneNumber(ctx, req.Email)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errTmp
	})
	err = g.Wait()
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if accessData == nil {
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessData.AccessType != enum.ACCESS_RESET_EMAIL {
		logData.Error = "invalid access type"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessData.ExpiredAt.Before(time.Now()) {
		logData.Error = "access expired"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_EXPIRED_ACCESS_CODE,
			Message:    pkgErr.EXPIRED_TIME_MSG,
		}
	}
	if userIsExist != nil {
		logData.Error = "email already registered"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_EMAIL_ALREADY_REGISTERED_CODE,
			Message:    pkgErr.EMAIL_ALREADY_REGISTERED_MSG,
		}
	}
	return svc.sendOtpUpdateProfile(ctx, &entity.OTP{
		OtpPurpose:     enum.OTP_RESET_EMAIL,
		OtpMethod:      enum.TYPE_EMAIL,
		OtpDestination: req.Email,
		UserId:         user.ID,
		UserUUID:       user.UUID,
	}, &dto.ProfileUpdateRequest{
		Field: enum.PROFILE_UPDATE_EMAIL,
		Value: req.Email,
		User:  user,
	}, logData)

}
func (svc *userProfileService) ResetPhoneNumber(ctx context.Context, req *request.ResetPhoneNumberRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		err         error
		userIsExist *entity.User
		user        *entity.User
		accessData  *entity.AccessState
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		accessData, errTmp = svc.AccessStateRepository.SelectAccessByAccessTokenRepository(ctx, req.AccountToken)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil
	})

	g.Go(func() error {
		strUserJson, errTmp := svc.redis.GetAccessKey(ctx, req.AccountToken)
		if errTmp != nil {
			return errTmp
		}
		if strUserJson != "" {
			errTmp = json.Unmarshal([]byte(strUserJson), &user)
			if errTmp != nil {
				return errTmp
			}
			return nil
		}
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil

	})
	g.Go(func() error {
		var errTmp error
		userIsExist, errTmp = svc.userRepository.SelectUserByEmailOrPhoneNumber(ctx, req.PhoneNumber)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errTmp
	})
	err = g.Wait()
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if accessData == nil {
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessData.AccessType != enum.ACCESS_RESET_PHONE_NUMBER {
		logData.Error = "invalid access type"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessData.ExpiredAt.Before(time.Now()) {
		logData.Error = "access expired"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_EXPIRED_ACCESS_CODE,
			Message:    pkgErr.EXPIRED_TIME_MSG,
		}
	}
	if userIsExist != nil {
		logData.Error = "email already registered"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_PHONE_NUMBER_ALREADY_REGISTERED_CODE,
			Message:    pkgErr.PHONE_NUBMBER_ALREADY_REGISTERED_MSG,
		}
	}
	return svc.sendOtpUpdateProfile(ctx, &entity.OTP{
		OtpPurpose:     enum.OTP_RESET_PHONE_NUMBER,
		OtpMethod:      enum.TYPE_SMS,
		OtpDestination: req.PhoneNumber,
		UserId:         user.ID,
		UserUUID:       user.UUID,
	}, &dto.ProfileUpdateRequest{
		Field: enum.PROFILE_UPDATE_PHONE_NUMBER,
		Value: req.PhoneNumber,
		User:  user,
	}, logData)

}
func (svc *userProfileService) VerifyOtpService(ctx context.Context, req *request.VerifyOtpRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		otpData           *entity.OTP
		updateProfileData *dto.ProfileUpdateRequest
		user              *entity.User
		updateUser        entity.User
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		return errTmp
	})
	g.Go(func() error {
		var errTmp error
		updateProfileData, errTmp = svc.redis.GetDataUpdateProfile(ctx, req.VerifyID)
		return errTmp
	})

	err := g.Wait()
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
		}
	}

	otpData, err = svc.otpService.VerifyOtpCode(ctx, &request.VerifyOtpRequest{
		Otp:        req.Otp,
		VerifyID:   req.VerifyID,
		DeviceID:   req.DeviceID,
		DeviceInfo: req.DeviceInfo,
	})
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_OTP_CODE,
			Message:    pkgErr.INVALID_OTP_MSG,
			Error:      err.Error(),
		}
	}
	if updateProfileData == nil {
		logData.Error = "record Profile Not Found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_RECORD_NOT_FOUND_CODE,
			Message:    pkgErr.RECORD_NOT_FOUND_MSG,
		}
	}

	if user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	logData.Data = map[string]interface{}{
		"verify_by":   otpData.OtpMethod,
		"destination": otpData.OtpDestination,
		"update_data": updateProfileData.Field,
		"value":       updateProfileData.Value,
	}
	if user.DeviceID != req.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	switch updateProfileData.Field {
	case enum.PROFILE_UPDATE_PIN:
		hashPin, _ := bcrypt.GenerateFromPassword([]byte(updateProfileData.Value), bcrypt.DefaultCost)
		updateUser.Pin = string(hashPin)
	case enum.PROFILE_UPDATE_EMAIL:
		updateUser.Email = updateProfileData.Value
	case enum.PROFILE_UPDATE_PHONE_NUMBER:
		updateUser.PhoneNumber = updateProfileData.Value
	}
	userTx := svc.userRepository.Tx(ctx)
	err = svc.userRepository.UpdateUser(ctx, userTx, user, &updateUser)
	if err != nil {
		logData.Error = err.Error()
		userTx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY}
	}
	userTx.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}
func (svc *userProfileService) ResetFullNameService(ctx context.Context, req *request.ResetFullNameRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		user *entity.User
	)
	user, err := svc.userRepository.SelectUserByUUID(ctx, *userUUID)
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
				Message:    pkgErr.USER_NOT_FOUND_MSG,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}

	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email

	if req.DeviceID != user.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}

	}
	txuser := svc.userRepository.Tx(ctx)
	err = svc.userRepository.UpdateUser(ctx, txuser, user, &entity.User{FullName: req.FullName})
	if err != nil {
		txuser.Rollback()
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	txuser.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}
func (svc *userProfileService) DeleteAccountService(ctx context.Context, req *request.DeletetAccountRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		accessState *entity.AccessState
		user        *entity.User
	)
	g := errgroup.Group{}
	g.Go(func() error {
		strUserJson, errTmp := svc.redis.GetAccessKey(ctx, req.AccountToken)
		if errTmp != nil {
			return errTmp
		}
		if strUserJson != "" {
			errTmp = json.Unmarshal([]byte(strUserJson), &user)
			if errTmp != nil {
				return errTmp
			}
			return nil
		}
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil
	})
	g.Go(func() error {
		var errTmp error
		accessState, errTmp = svc.AccessStateRepository.SelectAccessByAccessTokenRepository(ctx, req.AccountToken)
		return errTmp
	})
	err := g.Wait()
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}

	if user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	if user.DeviceID != req.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	if accessState == nil {
		logData.Error = "invalid access state"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessState.AccessType != enum.ACCESS_DELETE_ACCOUNT {
		logData.Error = "invalid access type"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_ACCESS_CODE,
			Message:    pkgErr.INVALID_ACCESS_KEY_MSG,
		}
	}
	if accessState.ExpiredAt.Before(time.Now()) {
		logData.Error = "access expired"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_EXPIRED_ACCESS_CODE,
			Message:    pkgErr.EXPIRED_TIME_MSG,
		}
	}
	txuser := svc.userRepository.Tx(ctx)
	txAccess := svc.AccessStateRepository.Tx(ctx)
	g = errgroup.Group{}
	g.Go(func() error {
		return svc.AccessStateRepository.UpdateAccessByStruct(ctx, txAccess, accessState, &entity.AccessState{Used: true})
	})
	g.Go(func() error {
		return svc.userRepository.UpdateUser(ctx, txuser, user, &entity.User{Status: enum.USER_INACTIVE})
	})
	err = g.Wait()
	if err != nil {
		logData.Error = err.Error()
		txuser.Rollback()
		txAccess.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	err = svc.userRepository.DeleteUser(ctx, txuser, user)
	if err != nil {
		logData.Error = err.Error()
		txuser.Rollback()
		txAccess.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	txuser.Commit()
	txAccess.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}

}
func (svc *userProfileService) BiometricStatusService(ctx context.Context, req *request.BiometrictStatusRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		userDetail *entity.UserDetail
		user       *entity.User
		err        error
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		userDetail, err = svc.userDetailRepository.SelectUserDetailByUserUUID(ctx, userUUID)
		if errTmp != nil {
			errors.Is(errTmp, gorm.ErrRecordNotFound)
			return nil
		}
		return errTmp
	})
	g.Go(func() error {
		var errTmp error
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			errors.Is(errTmp, gorm.ErrRecordNotFound)
			return nil
		}
		return errTmp
	})

	if err = g.Wait(); err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
		}
	}
	if userDetail == nil || user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	if req.DeviceID != user.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	biometrictStatus := enum.BIOMETRIC_ACTIVE
	if !req.Active {
		biometrictStatus = enum.BIOMETRIC_IN_ACTIVE
	}
	txUserDetail := svc.userRepository.Tx(ctx)
	if err = svc.userDetailRepository.UpdateUserDetail(ctx, txUserDetail, userDetail, &entity.UserDetail{Biometric: biometrictStatus}); err != nil {
		logData.Error = err.Error()
		txUserDetail.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	txUserDetail.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}

func (svc *userProfileService) ResetProfilePictureService(ctx context.Context, req *request.ResetProfilePictureRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		userDt *entity.UserDetail
		user   *entity.User
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errTmp
		}
		return nil
	})
	g.Go(func() error {
		var errTmp error
		userDt, errTmp = svc.userDetailRepository.SelectUserDetailByUserUUID(ctx, userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errTmp
	})
	if err := g.Wait(); err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}

	if userDt == nil || user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	if req.DeviceID != user.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	pictureName := fmt.Sprintf("%s_%s_%s", *userUUID, fmt.Sprint(time.Now().Unix()), "profile")
	_, objectName, err := svc.minioRepository.PutObject(ctx, req.ProfilePicture, enum.MinioPathNameMap[enum.MINIO_USER_PROFILE_PICTURE], &pictureName)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	txUserDetail := svc.userDetailRepository.Tx(ctx)
	err = svc.userDetailRepository.UpdateUserDetail(ctx, txUserDetail, userDt, &entity.UserDetail{ProfilePicture: *objectName})
	if err != nil {
		logData.Error = err.Error()
		txUserDetail.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	txUserDetail.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       objectName,
	}
}

func (svc *userProfileService) GetProfilePictureController(ctx context.Context, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	customResource, ok := ctx.Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		err := errors.New("failed to get custom resource")
		svc.clogger.ErrorLogger(ctx, "AccessTokenService", err)
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		}
	}
	var (
		userDt *entity.UserDetail
		user   *entity.User
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errTmp error
		user, errTmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}

		}
		return errTmp
	})
	g.Go(func() error {
		var errTmp error
		userDt, errTmp = svc.userDetailRepository.SelectUserDetailByUserUUID(ctx, userUUID)
		if errTmp != nil {
			if errors.Is(errTmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errTmp
	})
	if err := g.Wait(); err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	if userDt == nil || user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	if customResource.HeaderXDeviceID != user.DeviceID {
		logData.Error = "invalid device id"
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	if userDt.ProfilePicture == "" {
		err := errors.New("profile picture not found")
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.PROFILE_INVALID_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		}

	}
	texpired := time.Now().Add(svc.rootConfig.Minio.MinioPresignedDuration)
	url, err := svc.minioRepository.GenerateMinioPresignedURL(ctx, &userDt.ProfilePicture, svc.rootConfig.Minio.MinioPresignedDuration)
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
		Data: response.GetProfilePictureResponse{
			Url:      url,
			ExpireAt: texpired,
			UserUUID: user.UUID,
		},
	}
}
