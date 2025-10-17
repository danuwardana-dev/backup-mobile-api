package verihubsInvokerService

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type verihubsInvokerService struct {
	OTPRepository postgres.OtpRepository
	Clog          *helpers.CustomLogger
}
type VerihubsInvokerService interface {
	OtpInvokerService(ctx context.Context, req *request.VerihubsOtpInvoker) *dto.BaseResponse
}

func NewVerihubsInvokerService(OTPRepository postgres.OtpRepository, Clog *helpers.CustomLogger) VerihubsInvokerService {
	return &verihubsInvokerService{OTPRepository: OTPRepository, Clog: Clog}
}
func (svc *verihubsInvokerService) OtpInvokerService(ctx context.Context, req *request.VerihubsOtpInvoker) *dto.BaseResponse {
	otpData, err := svc.OTPRepository.SelectOtpBySessionId(ctx, req.SessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.OUTBOUND_RECORD_NOT_FOUND_CODE,
				Message:    pkgErr.RECORD_NOT_FOUND_MSG,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.OUTBOUND_UNDIFINED_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	tx := svc.OTPRepository.Tx(ctx)
	statusNum, _ := strconv.Atoi(req.Status)
	status := enum.OtpSMSMapStatusVerihubs[statusNum]
	if status == "" {
		err = fmt.Errorf("out of map status on status num: %s", req.Status)
		svc.Clog.ErrorLogger(ctx, "OtpInvoker.enum.OtpSMSMapStatusVerihubs", err)
		return &dto.BaseResponse{
			StatusCode: pkgErr.OUTBOUND_INVALID_PAYLOAD,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      "invalid status",
		}
	}
	err = svc.OTPRepository.UpdateOtpDataRepository(ctx, tx, otpData, &entity.OTP{
		Status: status,
	})
	if err != nil {
		tx.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.OUTBOUND_UNDIFINED_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
		}
	}
	tx.Commit()
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
	}
}
