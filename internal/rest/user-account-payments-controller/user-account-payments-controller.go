package useraccountpaymentscontroller

import (
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum/pkgErr"
	service "backend-mobile-api/service/user-accounts-payment-svc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserAccountPaymentsController struct {
	service service.UserAccountsPaymentService
}

func NewUserAccountPaymentsController(service service.UserAccountsPaymentService) UserAccountPaymentsController {
	return UserAccountPaymentsController{service: service}
}
func (c UserAccountPaymentsController) GetUserAccountPaymentsController(ctx echo.Context) error {
	search := ctx.QueryParam("search")

	var (
		userPaymentsAccount []entity.UserPaymentsAccount
		err                 error
	)
	if search != "" {
		userPaymentsAccount, err = c.service.SearchUserAccountsPaymentService(ctx.Request().Context(), search)
	}
	if search == "" {
		userPaymentsAccount, err = c.service.GetUserAccountsPaymentService(ctx.Request().Context())
	}

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data:       userPaymentsAccount,
	})
}
