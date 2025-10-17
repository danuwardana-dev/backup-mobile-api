package kycservice

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/outbond/verihubs"
	"backend-mobile-api/internal/repository/minio"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	verihubsDto "backend-mobile-api/model/outbond/verihubs-dto"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type KycService interface {
	VerifyKycKTP(ctx context.Context, req request.KycRequest) (*dto.BaseResponse, *dto.CustomLoggerRequest)
	VerifyKycPassport(ctx context.Context, req request.KycRequest) (*dto.BaseResponse, *dto.CustomLoggerRequest)
	VerifyPhotoSelfie(ctx context.Context, req request.VerifyKycSelfie) (*dto.BaseResponse, *dto.CustomLoggerRequest)
	SaveKycKTP(ctx context.Context, req request.KTPrequest, userUUID string) (*dto.BaseResponse, *dto.CustomLoggerRequest)
	SaveKycPassport(ctx context.Context, req request.PassportRequest, userUUID string) (*dto.BaseResponse, *dto.CustomLoggerRequest)
}

type kycService struct {
	clogger               *helpers.CustomLogger
	outboundVeriHubsSvc   verihubs.OutboundVeriHubsService
	KycKtpRepository      postgres.KycKtpRepository
	KycPassportRepository postgres.KycPassportRepository
	UserRepository        postgres.UserRepository
	UserDetailsRepository postgres.UserDetailRepository
	MinioRepository       minio.MinioRepository
}

func (k *kycService) VerifyPhotoSelfie(ctx context.Context, req request.VerifyKycSelfie) (*dto.BaseResponse, *dto.CustomLoggerRequest) {

	var (
		err             error
		logData         = &dto.CustomLoggerRequest{Remarks: "kyc-verify-selfie", Success: false}
		verihubResponse *verihubsDto.VerifySelfieResponse
	)
	kycRequest := verihubsDto.VerifyKycSelfie{
		Nik:         req.Nik,
		Name:        req.Name,
		BirthDate:   req.BirthDate,
		Email:       req.Email,
		Phone:       req.Phone,
		SelfiePhoto: req.SelfiePhoto,
		KtpPhoto:    req.KtpPhoto,
		IsLiveness:  true,
	}

	verihubResponse, err = k.outboundVeriHubsSvc.SendVerifySelfie(ctx, &kycRequest)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       verihubResponse.Data,
	}, nil
}

// SaveKycPassport implements KycService.
func (k *kycService) SaveKycPassport(ctx context.Context, req request.PassportRequest, userUUID string) (*dto.BaseResponse, *dto.CustomLoggerRequest) {
	var (
		err     error
		tx      = k.KycKtpRepository.Tx(ctx)
		logData = &dto.CustomLoggerRequest{Remarks: "kyc-save-passport", Success: false}
	)

	userData, err := k.UserRepository.SelectUserByStructOne(ctx, &entity.User{
		UUID: userUUID,
	})

	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.GetUser.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	//Get user detail
	userDetails, err := k.UserDetailsRepository.SelectUserDetailByUserId(ctx, userData.ID)
	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.SelectUserDetailByUserId.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	//Update the detail address
	newUserDetails := userDetails
	newUserDetails.Address = req.CurrentAddress.FullAddress
	newUserDetails.District = req.CurrentAddress.District
	newUserDetails.Country = req.CurrentAddress.Country
	newUserDetails.Regency = req.CurrentAddress.City

	fileName := fmt.Sprintf("%s-%s", userData.UUID, "pasport")
	fileHeader, err := Base64ToMultipartFileHeader(req.Image, fileName, "image/jpeg")
	if err != nil {
		tx.Rollback()
		k.clogger.ErrorLogger(ctx, "KYCIdenityPassport.PutObject", err)
		return nil, logData
	}

	_, objectName, err := k.MinioRepository.PutObject(ctx, fileHeader, string(enum.MINIO_KYC_OCR), &fileName)
	if err != nil {
		k.clogger.ErrorLogger(ctx, "KYCIdenityPassport.PutObject", err)
		return nil, logData
	}

	err = k.KycPassportRepository.SaveKYCPassport(tx, ctx, entity.EncapsulateRequestPassportToEntity(req, userData, *objectName))
	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.SaveKYCKtp.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	err = k.UserDetailsRepository.UpdateUserDetail(ctx, tx, userDetails, newUserDetails)
	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.UpdateUserDetail.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	tx.Commit()
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       "Success to save KTP",
	}, nil
}

func Base64ToMultipartFileHeader(base64Str, fileName, contentType string) (*multipart.FileHeader, error) {
	// Decode base64
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	// Create a buffer with the decoded data
	buffer := bytes.NewBuffer(data)

	// Create a new multipart writer
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Create a form file
	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	// Copy the data to the form file
	if _, err = io.Copy(fw, buffer); err != nil {
		return nil, err
	}

	// Close the writer
	w.Close()

	// Create a new request with the form data
	req, err := http.NewRequest("POST", "http://dummy-url", &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	// Parse the form to get the file header
	err = req.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		return nil, err
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	file.Close()

	// Set the content type
	handler.Header.Set("Content-Type", contentType)

	return handler, nil
}

// SaveKycKTP implements KycService.
func (k *kycService) SaveKycKTP(ctx context.Context, req request.KTPrequest, userUUID string) (*dto.BaseResponse, *dto.CustomLoggerRequest) {
	var (
		err     error
		tx      = k.KycKtpRepository.Tx(ctx)
		logData = &dto.CustomLoggerRequest{Remarks: "kyc-save-ktp", Success: false}
	)

	userData, err := k.UserRepository.SelectUserByStructOne(ctx, &entity.User{
		UUID: userUUID,
	})

	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.GetUser.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	//Get user detail
	userDetails, err := k.UserDetailsRepository.SelectUserDetailByUserId(ctx, userData.ID)
	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.SelectUserDetailByUserId.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	//Update the detail address
	newUserDetails := userDetails
	newUserDetails.Address = req.CurrentAddress.FullAddress
	newUserDetails.District = req.CurrentAddress.District
	newUserDetails.Country = req.CurrentAddress.Country
	newUserDetails.Regency = req.CurrentAddress.City

	fileName := fmt.Sprintf("%s-%s", userData.UUID, "ktp")
	fileHeader, err := Base64ToMultipartFileHeader(req.Image, fileName, "image/jpeg")
	if err != nil {
		tx.Rollback()
		k.clogger.ErrorLogger(ctx, "KYCIdenityKTP.PutObject", err)
		return nil, logData
	}

	_, objectPath, err := k.MinioRepository.PutObject(ctx, fileHeader, "kyc/ktp", &fileName)
	if err != nil {
		tx.Rollback()
		k.clogger.ErrorLogger(ctx, "KYCIdenityKTP.PutObject", err)
		return nil, logData
	}

	err = k.KycKtpRepository.SaveKYCKtp(ctx, tx, entity.EncapsulateRequestKtpToEntity(req, userData, *objectPath))
	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.SaveKYCKtp.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	err = k.UserDetailsRepository.UpdateUserDetail(ctx, tx, userDetails, newUserDetails)
	if err != nil {
		k.clogger.ErrorLogger(ctx, "SaveKycKTP.UpdateUserDetail.FailOnDb", err)
		logData.Error = err.Error()
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	tx.Commit()
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       "Success to save KTP",
	}, nil
}

// VerifyKycKTP implements KycService.
func (k *kycService) VerifyKycKTP(ctx context.Context, req request.KycRequest) (*dto.BaseResponse, *dto.CustomLoggerRequest) {
	var (
		err             error
		logData         = &dto.CustomLoggerRequest{Remarks: "kyc-verify", Success: false}
		verihubResponse *verihubsDto.IdentityKTPResponse
	)

	kycRequest := verihubsDto.VerihubIdentityRequest{
		ValidateQuality: true,
		Image:           req.Image,
	}

	verihubResponse, err = k.outboundVeriHubsSvc.SendKYCIdenityKTP(ctx, &kycRequest)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       verihubResponse.Data,
	}, nil

}

// VerifyKycPasport implements KycService.
func (k *kycService) VerifyKycPassport(ctx context.Context, req request.KycRequest) (*dto.BaseResponse, *dto.CustomLoggerRequest) {
	var (
		err             error
		logData         = &dto.CustomLoggerRequest{Remarks: "kyc-verify", Success: false}
		verihubResponse *verihubsDto.IdentityPassportResponse
	)

	kycRequest := verihubsDto.VerihubIdentityRequest{
		ValidateQuality: true,
		Image:           req.Image,
	}

	verihubResponse, err = k.outboundVeriHubsSvc.SendKYCIdenityPassport(ctx, &kycRequest)
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		}, logData
	}

	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       verihubResponse,
	}, nil
}

func NewKycService(
	clogger *helpers.CustomLogger,
	outboundVeriHubsSvc verihubs.OutboundVeriHubsService,
	kycPassportRepository postgres.KycPassportRepository,
	kycKtpRepository postgres.KycKtpRepository,
	userRepository postgres.UserRepository,
	userDetailsRepository postgres.UserDetailRepository,
	minioRepository minio.MinioRepository) KycService {
	return &kycService{
		clogger:               clogger,
		outboundVeriHubsSvc:   outboundVeriHubsSvc,
		KycKtpRepository:      kycKtpRepository,
		KycPassportRepository: kycPassportRepository,
		UserRepository:        userRepository,
		UserDetailsRepository: userDetailsRepository,
		MinioRepository:       minioRepository,
	}
}
