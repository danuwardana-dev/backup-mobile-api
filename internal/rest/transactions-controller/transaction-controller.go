package transactionController

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum/pkgErr"
	service "backend-mobile-api/service/transactions-svc"
	"encoding/json"

	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	service service.TransactionService
}

func NewTransactionController(service service.TransactionService) TransactionController {
	if service == nil {
		log.Println("[ERROR] service nil saat init controller")
	}
	return TransactionController{service: service}
}

// DTO khusus request
type TransactionRequest struct {
	UserUUID      string                           `json:"user_uuid"`
	Type          string                           `json:"type"`
	Description   string                           `json:"description"`
	Nominal       int64                            `json:"nominal"`
	AdminFee      int64                            `json:"admin_fee"`
	UniqueCode    int64                            `json:"unique_code"`
	PaymentMethod string                           `json:"payment_method"`
	Status        string                           `json:"status"`
	AccountNumber string                           `json:"account_number"`
	ImageURL      string                           `json:"image_url"`
	ExpiredAt     time.Time                        `json:"expired_at"`
	BankTransfer  *entity.TransactionBankTransfer  `json:"bank_transfer,omitempty"`
	Ewallet       *entity.TransactionEwallet       `json:"ewallet,omitempty"`
	PhoneCredit   *entity.TransactionPhoneCredit   `json:"phone_credit,omitempty"`
	InternetTV    *entity.TransactionInternetTV    `json:"internet_tv,omitempty"`
	International *entity.TransactionInternational `json:"international,omitempty"`
}
type UpdateStatusRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	Status        string `json:"status" validate:"required,oneof=pending success failed canceled expired"`
}

// object request
type GenerateCodeRequest struct {
	Type string `json:"type" validate:"required,oneof=bank_transfer va"`
}

// ✅ POST /transactions/generate
func (c TransactionController) GenerateTransactionCode(ctx echo.Context) error {
	var req GenerateCodeRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	// panggil service
	codeData, err := c.service.GenerateTransactionCode(ctx.Request().Context(), req.Type)
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
		Data:       codeData,
	})
}

// ✅ POST /transactions
func (c TransactionController) CreateTransaction(ctx echo.Context) error {
	var req TransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	// hitung total otomatis
	total := req.Nominal + req.AdminFee

	// mapping DTO → Entity
	tx := entity.Transaction{
		Type:          req.Type,
		Description:   req.Description,
		Nominal:       float64(req.Nominal),
		AdminFee:      float64(req.AdminFee),
		UniqueCode:    float64(req.UniqueCode),
		PaymentMethod: req.PaymentMethod,
		Total:         float64(total),
		Status:        req.Status,
		ExpiredAt:     req.ExpiredAt,
	}

	// create transaksi utama
	newtx, err := c.service.CreateTransaction(ctx.Request().Context(), &tx, req.UserUUID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		})
	}

	// sesuai type → insert detail
	switch req.Type {
	case "bank_transfer":
		if req.BankTransfer != nil {
			req.BankTransfer.TransactionID = newtx.TransactionID
			if err := c.service.AddTransactionBankTransfer(ctx.Request().Context(), req.BankTransfer); err != nil {
				return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
					StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
					Message:    pkgErr.INTERNAL_SERVER_MSG,
					Error:      err.Error(),
				})
			}
		}

	case "ewallet":
		if req.Ewallet != nil {
			req.Ewallet.TransactionID = newtx.TransactionID
			if err := c.service.AddTransactionEwallet(ctx.Request().Context(), req.Ewallet); err != nil {
				return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
					StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
					Message:    pkgErr.INTERNAL_SERVER_MSG,
					Error:      err.Error(),
				})
			}
		}

	case "phone_credit":
		if req.PhoneCredit != nil {
			req.PhoneCredit.TransactionID = newtx.TransactionID
			if err := c.service.AddTransactionPhoneCredit(ctx.Request().Context(), req.PhoneCredit); err != nil {
				return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
					StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
					Message:    pkgErr.INTERNAL_SERVER_MSG,
					Error:      err.Error(),
				})
			}
		}

	case "internet_tv":
		if req.InternetTV != nil {
			req.InternetTV.TransactionID = newtx.TransactionID
			if err := c.service.AddTransactionInternetTV(ctx.Request().Context(), req.InternetTV); err != nil {
				return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
					StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
					Message:    pkgErr.INTERNAL_SERVER_MSG,
					Error:      err.Error(),
				})
			}
		}

	case "international":
		if req.International != nil {
			req.International.TransactionID = newtx.TransactionID
			if err := c.service.AddTransactionInternational(ctx.Request().Context(), req.International); err != nil {
				return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
					StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
					Message:    pkgErr.INTERNAL_SERVER_MSG,
					Error:      err.Error(),
				})
			}
		}
	}

	return ctx.JSON(http.StatusCreated, dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       newtx, // ⚡ user_uuid tidak akan ikut kebawa lagi
	})
}

// ✅ GET /transactions/:id
func (c TransactionController) GetTransaction(ctx echo.Context) error {
	transactionID := ctx.Param("id")

	tx, err := c.service.GetTransactionByID(ctx.Request().Context(), transactionID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, dto.BaseResponse{
			StatusCode: pkgErr.SUCCESS_CODE,
			Message:    "transaction not found",
			Error:      err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       tx,
	})
}

// ✅ PUT /transactions/status
func (c TransactionController) UpdateTransactionStatus(ctx echo.Context) error {
	var req UpdateStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}

	// panggil service untuk update status + kirim notif
	if err := c.service.UpdateTransactionStatus(ctx.Request().Context(), req.TransactionID, req.Status); err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    "Transaction status updated successfully",
	})
}

// ✅ GET /transactions
func (c TransactionController) GetAllTransactions(ctx echo.Context) error {
	userUUID := ctx.QueryParam("user_uuid")
	if userUUID == "" {
		return ctx.JSON(http.StatusUnauthorized, dto.BaseResponse{
			StatusCode: pkgErr.INVALID_REQUEST_PAYLOAD_CODE,
			Message:    "user_uuid not found in query parameter",
		})
	}

	// Pagination
	pageParam := ctx.QueryParam("page")
	limitParam := ctx.QueryParam("limit")
	search := ctx.QueryParam("search")
	status := ctx.QueryParam("status")
	txType := ctx.QueryParam("type")
	transactionID := ctx.QueryParam("transaction_id")
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")

	page := 1
	limit := 10
	if pageParam != "" {
		_ = json.Unmarshal([]byte(pageParam), &page)
	}
	if limitParam != "" {
		_ = json.Unmarshal([]byte(limitParam), &limit)
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// --- Panggil service ---
	txs, total, err := c.service.GetAllTransactionsPaginated(
		ctx.Request().Context(),
		userUUID,
		limit,
		offset,
		search,
		status,
		txType,
		transactionID,
		startDate,
		endDate,
	)

	if err != nil {
		helpers.CustomeLogger(ctx.Request().Context(), &dto.CustomLoggerRequest{
			Error:   err.Error(),
			Remarks: "failed to get transactions",
			Success: false,
		})
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.INTERNAL_SERVER_ERROR_CODE,
			Message:    pkgErr.INTERNAL_SERVER_MSG,
			Error:      err.Error(),
		})
	}

	response := map[string]interface{}{
		"total":  total,
		"page":   page,
		"limit":  limit,
		"result": txs,
	}

	helpers.CustomeLogger(ctx.Request().Context(), &dto.CustomLoggerRequest{
		Success: true,
		Remarks: "get all transactions success",
		Data:    response,
	})

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Data:       response,
	})
}
