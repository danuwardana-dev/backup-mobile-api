package recipientController

import (
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum/pkgErr"
	service "backend-mobile-api/service/recipient-svc"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RecipientController struct {
	service service.RecipientService
}

func NewRecipientController(service service.RecipientService) RecipientController {
	if service == nil {
		log.Println("[ERROR] service nil saat init controller")
	}
	return RecipientController{service: service}
}

// ✅ GET /recipients?search=xxx
func (c RecipientController) GetRecipients(ctx echo.Context) error {
	keyword := ctx.QueryParam("search")

	var (
		recipients []entity.RecipientWithBank
		err        error
	)

	if keyword != "" {
		recipients, err = c.service.SearchRecipients(ctx.Request().Context(), keyword)
	} else {
		recipients, err = c.service.GetRecipients(ctx.Request().Context())
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
		Data:       recipients,
	})
}

// ✅ POST /recipients
func (c RecipientController) CreateRecipient(ctx echo.Context) error {
	var req entity.Recipient

	// bind payload
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	// pastikan UUID ada di payload
	if req.UserUUID == "" {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    "user uuid is required",
		})
	}

	// 🔑 ambil user dari service (yg nanti ke repo)
	user, err := c.service.GetUserByUUID(ctx.Request().Context(), req.UserUUID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    "user not found",
			Error:      err.Error(),
		})
	}

	// isi user_id dari DB
	req.User = user.ID

	// simpan recipient
	if err := c.service.AddRecipient(ctx.Request().Context(), &req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       req,
	})
}
