package otp

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/outbond/smtp"
	"backend-mobile-api/internal/outbond/verihubs"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/internal/repository/redis"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	verihubsDto "backend-mobile-api/model/outbond/verihubs-dto"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"time"
)

type otpService struct {
	otpRepository           postgres.OtpRepository
	rootConfig              *config.Root
	redis                   redis.Redis
	clogger                 helpers.CustomLogger
	outboundVerihubsService verihubs.OutboundVeriHubsService
	smtp                    smtp.Smtp
}

func NewOtpService(
	otpRepository postgres.OtpRepository,
	rootConfig *config.Root,
	redis *redis.Redis,
	clogger *helpers.CustomLogger,
	outboundVerihubsService verihubs.OutboundVeriHubsService,
	smtp *smtp.Smtp,
) OtpService {
	return &otpService{
		otpRepository:           otpRepository,
		rootConfig:              rootConfig,
		redis:                   *redis,
		clogger:                 *clogger,
		outboundVerihubsService: outboundVerihubsService,
		smtp:                    *smtp,
	}
}

type OtpService interface {
	generateOtpCode(ctx context.Context, verifyKey string) (string, error)
	SendOtp(ctx context.Context, req *dto.SendOtp) (*entity.OTP, error)
	VerifyOtpCode(ctx context.Context, req *request.VerifyOtpRequest) (*entity.OTP, error)
}

func (svc *otpService) generateOtpCode(ctx context.Context, verifyKey string) (string, error) {
	var (
		langth = 6
		table  = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
		exist  bool
	)
	b := make([]byte, langth)
	n, err := io.ReadAtLeast(rand.Reader, b, langth)
	if n != langth {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	if exist, err = svc.redis.OtpIsExist(ctx, string(b), verifyKey); err != nil {
		svc.clogger.ErrorLogger(ctx, "GenerateOtp.svc.redis.OtpIsExist", err)
		return "", err
	}
	if exist {
		return svc.generateOtpCode(ctx, verifyKey)
	}
	return string(b), err
}
func (svc *otpService) SendOtp(ctx context.Context, req *dto.SendOtp) (*entity.OTP, error) {
	var (
		otpData = entity.OTP{
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
	CustomOtp, err := svc.generateOtpCode(ctx, req.VerifyKey)
	if err != nil {
		return nil, err
	}
	otpData.Status = enum.OTP_DELIVERED
	switch req.OtpMethod {
	case enum.TYPE_EMAIL:
		if err := svc.smtp.SendMail(ctx, []string{req.OtpDestination}, enum.VERIFY_OTP_SUBJECT, svc.smtp.RegisterOtpMsg(CustomOtp)); err != nil {
			svc.clogger.ErrorLogger(ctx, "sendOtp.smtp.SendMail", err)
			otpData.Status = enum.OTP_FAILED
		}
		otpData.OtpCode = CustomOtp
		otpData.SessionId = req.VerifyKey
	case enum.TYPE_SMS:
		var otp *string
		if !svc.rootConfig.Verihubs.VerihubsOTPCode {
			otp = &CustomOtp
		}
		resp, err := svc.outboundVerihubsService.SendSMSOtpService(ctx, &verihubsDto.SendOtpBaseRequest{
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
		if resp != nil {
			otpData.OtpCode = resp.Otp
			otpData.SessionId = resp.SessionID
		}

	case enum.TYPE_WHATSAPP:
		var otp *string
		if !svc.rootConfig.Verihubs.VerihubsOTPCode {
			otp = &CustomOtp
		}
		resp, err := svc.outboundVerihubsService.SendWhatsappsService(ctx, &verihubsDto.SendWhatsappOtpBaseRequest{
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
		if resp != nil {
			otpData.OtpCode = resp.Otp
			otpData.SessionId = resp.SessionID
			otpData.UserUUID = req.UserUUID
		}
	}
	tx := svc.otpRepository.Tx(ctx)
	err = svc.otpRepository.InsertOtpDataRepository(ctx, tx, &otpData)
	if err != nil {
		svc.clogger.ErrorLogger(ctx, "otpRepository.InsertOtpDataRepository", err)
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return svc.otpRepository.SelectOtpByVerifyKeyBeforeExpire(ctx, req.VerifyKey)
}

func (svc *otpService) VerifyOtpCode(ctx context.Context, req *request.VerifyOtpRequest) (*entity.OTP, error) {
	var (
		otpData   *entity.OTP
		otpUpdate *entity.OTP
	)
	strJson, err := svc.redis.GetOtp(ctx, req.Otp, req.VerifyID)
	if err != nil {
		svc.clogger.ErrorLogger(ctx, "VerifyOtpCode", err)
		return nil, err
	}
	if strJson != "" {
		var otpTmp entity.OTP
		err = json.Unmarshal([]byte(strJson), &otpTmp)
		if err != nil {
			svc.clogger.ErrorLogger(ctx, "VerifyOtpCode", err)
			return nil, err
		}
		otpData = &otpTmp
		_ = svc.redis.DeleteOtp(ctx, req.Otp, req.VerifyID)
	} else {
		otpData, err = svc.otpRepository.SelectOtpByVerifyKey(ctx, req.VerifyID)
		if err != nil {
			svc.clogger.ErrorLogger(ctx, "VerifyOtpCode.SelectOtpByVerifyKey", err)
			return nil, err
		}

	}
	if otpData.OtpCode != req.Otp {
		return nil, errors.New("otp code is invalid")
	}
	if otpData.Status == enum.OTP_BLOCKED {
		return nil, errors.New("otp code is blocked")
	}
	if otpData.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("otp is expired")
	}
	if otpData.Status == enum.OTP_VERIFIED {
		return nil, errors.New("otp already verified")
	}

	switch otpData.OtpMethod {
	case enum.TYPE_EMAIL:
		otpUpdate = &entity.OTP{Status: enum.OTP_VERIFIED}
	case enum.TYPE_SMS:
		if _, err := svc.outboundVerihubsService.VerifySMSOtpService(ctx, &verihubsDto.VerifyOtpBaseRequest{
			MSISDN:    otpData.OtpDestination,
			Otp:       otpData.OtpCode,
			Challenge: nil,
		}); err != nil {
			svc.clogger.ErrorLogger(ctx, "VerifyOtpCode", err)
			return nil, err
		}
		otpUpdate = &entity.OTP{Status: enum.OTP_VERIFIED}
	case enum.TYPE_WHATSAPP:
		if _, err = svc.outboundVerihubsService.VerifyWhatsappsOtpService(ctx, &verihubsDto.VerifyOtpBaseRequest{
			MSISDN:    otpData.OtpDestination,
			Otp:       otpData.OtpCode,
			Challenge: nil,
		}); err != nil {
			svc.clogger.ErrorLogger(ctx, "VerifyOtpCode", err)
			return nil, err
		}
		otpUpdate = &entity.OTP{Status: enum.OTP_VERIFIED}
	default:
		return nil, errors.New("invalid otp method")

	}
	if otpUpdate != nil {
		txOtp := svc.otpRepository.Tx(ctx)
		err = svc.otpRepository.UpdateOtpDataRepository(ctx, txOtp, otpData, otpUpdate)
		if err != nil {
			svc.clogger.ErrorLogger(ctx, "otpRepository.UpdateOtpDataRepository", err)
			txOtp.Rollback()
			return nil, err
		}
		txOtp.Commit()
		return svc.otpRepository.SelectOtpByVerifyKey(ctx, req.VerifyID)
	}
	return otpData, nil

}
