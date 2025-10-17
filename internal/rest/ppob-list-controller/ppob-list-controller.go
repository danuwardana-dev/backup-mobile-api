package ppoblistcontroller

import (
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum/pkgErr"
	service "backend-mobile-api/service/ppob-list-svc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PpobListController struct {
	service service.PpobListService
}

func NewPpobListController(service service.PpobListService) PpobListController {
	return PpobListController{service: service}
}

func (c PpobListController) GetPpobListController(ctx echo.Context) error {
	search := ctx.QueryParam("search")

	var (
		ppob []entity.PPOB
		err  error
	)
	if search != "" {
		ppob, err = c.service.SearchPpobList(ctx.Request().Context(), search)
	} else {
		ppob, err = c.service.GetPpobList(ctx.Request().Context())
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
		Data:       ppob,
	})

}
