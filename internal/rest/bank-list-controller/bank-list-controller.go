package bankListController

import (
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum/pkgErr"
	service "backend-mobile-api/service/bank-list-svc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BankListController struct {
	service service.BankService
}

func NewBankListController(service service.BankService) BankListController {
	return BankListController{service: service}
}

func (c BankListController) GetBankListController(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	bankType := ctx.QueryParam("type") // NEW → filter by type (BANK / EWALLET)

	var (
		banks []entity.Bank
		err   error
	)

	switch {
	case bankType != "" && search != "":
		banks, err = c.service.SearchBanksByType(ctx.Request().Context(), bankType, search)
	case bankType != "":
		banks, err = c.service.GetBanksByType(ctx.Request().Context(), bankType)
	case search != "":
		banks, err = c.service.SearchBanks(ctx.Request().Context(), search)
	default:
		banks, err = c.service.GetBanks(ctx.Request().Context())
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
		Data:       banks,
	})
}
