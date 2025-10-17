package verihubs_invoker_controller

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/enum/pkgErr"
	verihubsInvokerService "backend-mobile-api/service/verihubs-invoker-service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type verihubsInvokerController struct {
	verihubsInvokerService verihubsInvokerService.VerihubsInvokerService
}
type VerihubsInvokerController interface {
	InvokerOtpController(e echo.Context) error
}

func NewVerihubsInvokerController(verihubsInvokerService verihubsInvokerService.VerihubsInvokerService) VerihubsInvokerController {
	return &verihubsInvokerController{
		verihubsInvokerService: verihubsInvokerService,
	}
}
func (ctr *verihubsInvokerController) InvokerOtpController(e echo.Context) error {
	var (
		req      = request.VerihubsOtpInvoker{}
		validate = validator.New()
		err      error
	)

	err = e.Bind(&req)
	if err != nil {
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	err = validate.Struct(&req)
	if err != nil {
		log.Info(err.Error())
		err = helpers.CustomValidatePayload(err, req)
		return e.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.AUTH_INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	resp := ctr.verihubsInvokerService.OtpInvokerService(e.Request().Context(), &req)
	switch resp.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, resp)
	case pkgErr.OUTBOUND_INVALID_PAYLOAD:
		return e.JSON(http.StatusBadRequest, resp)
	case pkgErr.OUTBOUND_RECORD_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, resp)
	default:
		return e.JSON(http.StatusInternalServerError, resp)
	}
}
