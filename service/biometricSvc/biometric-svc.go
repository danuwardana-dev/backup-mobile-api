package biometricSvc

import (
	"backend-mobile-api/internal/middleware"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/dto/response"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type biometricService struct {
	userRepository       postgres.UserRepository
	userDetailRepository postgres.UserDetailRepository
	middleware           middleware.CustomMiddleware
	tokenBlacklist       postgres.TokenBlacklistTokenRepository
}

func NewBiometricService(
	userRepository postgres.UserRepository,
	userDetailRepository postgres.UserDetailRepository,
	middleware middleware.CustomMiddleware,
	tokenBlacklist postgres.TokenBlacklistTokenRepository,
) BiometricService {
	return &biometricService{
		userRepository:       userRepository,
		userDetailRepository: userDetailRepository,
		middleware:           middleware,
		tokenBlacklist:       tokenBlacklist,
	}
}

type BiometricService interface {
	VerifyBiometric(ctx context.Context, request *request.BiometricRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	login(ctx context.Context, request *request.BiometricRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse
}

func (svc *biometricService) login(ctx context.Context, request *request.BiometricRequest, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		user           *entity.User
		userDt         *entity.UserDetail
		token          *middleware.TokenData
		err            error
		ok             bool
		blacklistError = errors.New("token is blacklisted")
	)

	g := errgroup.Group{}
	g.Go(func() error {
		var errtmp error
		isExist, errtmp := svc.tokenBlacklist.IsBlaclistTokenActive(ctx, request.Token)
		if err != nil {
			return errtmp
		}
		if isExist {
			return blacklistError
		}
		return nil
	})
	g.Go(func() error {
		var errtmp error
		user, errtmp = svc.userRepository.SelectUserByUUID(ctx, request.UUID)
		if errtmp != nil {
			if errors.Is(errtmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errtmp
		}
		return nil
	})
	g.Go(func() error {
		var errtmp error
		userDt, errtmp = svc.userDetailRepository.SelectUserDetailByUserUUID(ctx, &request.UUID)
		if errtmp != nil {
			if errors.Is(errtmp, gorm.ErrRecordNotFound) {
				return nil
			}
			return errtmp
		}
		return nil
	})

	err = g.Wait()
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, blacklistError) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
				Message:    pkgErr.UNAUTHORIZED_MSG,
				Error:      err.Error(),
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
			Data:       nil,
		}
	}
	if user == nil || userDt == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_USER_NOT_FOUND_CODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	ok, token, err = svc.middleware.RefreshToken(ctx, request.Token, user)
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
	if userDt.Biometric != enum.BIOMETRIC_ACTIVE {
		logData.Error = "biometric inactive"
		return &dto.BaseResponse{
			StatusCode: pkgErr.AUTH_BIOMETRC_INACTIVE_CODE,
			Message:    pkgErr.BIOMETRIC_INACTIVE_MSG,
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
			Token: *token,
		},
	}
}

func (svc *biometricService) VerifyBiometric(ctx context.Context, request *request.BiometricRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	switch request.ServiceType {
	case enum.BIOMETRIC_LOGIN:
		return svc.login(ctx, request, logData)
	}
	return &dto.BaseResponse{
		StatusCode: pkgErr.BIOMETRIC_INVALID_REQUEST_CODE,
		Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
		Error:      "invalid service type biometric ",
		Data:       nil,
	}
}
