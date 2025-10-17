package checkaccountbankcontroller

import (
	"backend-mobile-api/model/entity"

	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type CheckAccountBankController struct{}

func NewCheckAccountBankController() CheckAccountBankController {
	return CheckAccountBankController{}
}

// GeneratePartnerReferenceNo bikin angka random sepanjang 12 digit
func GeneratePartnerReferenceNo() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Int63n(999999999999) // max 12 digit
	return leftPad(strconv.FormatInt(num, 10), "0", 12)
}

// helper: padding ke kiri biar tetap 12 digit
func leftPad(s, pad string, length int) string {
	for len(s) < length {
		s = pad + s
	}
	return s
}

// POST /check-account
func (CheckAccountBankController) CheckAccount(ctx echo.Context) error {
	var req entity.CheckAccountRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid payload",
		})
	}
	//struct status hanya dummy
	type BaseResponse struct {
		StatusCode string      `json:"status_code"`
		Message    string      `json:"message"`
		Error      string      `json:"error"`
		Data       interface{} `json:"data"`
	}
	// ✅ dummy response
	var err error
	// override partnerReferenceNo kalau kosong / ingin selalu generate
	req.PartnerReferenceNo = GeneratePartnerReferenceNo()
	resp := entity.CheckAccountResponseData{

		BeneficiaryAccountNo:   req.BeneficiaryAccountNo,
		BeneficiaryAccountName: "TIMO TEST", // dummy fixed

		BeneficiaryBankCode: req.BeneficiaryBankCode,
	}

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, BaseResponse{
			StatusCode: "99", // contoh error code
			Message:    "failed",
			Error:      err.Error(),
			Data:       nil,
		})
	}

	return ctx.JSON(http.StatusOK, BaseResponse{
		StatusCode: "00", // sukses selalu 00
		Message:    "success",
		Error:      "",
		Data:       resp, // langsung balikin response check account
	})
}
