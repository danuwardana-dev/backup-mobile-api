package dto

import "backend-mobile-api/model/enum/pkgErr"

type BaseResponse struct {
	StatusCode pkgErr.Code `json:"status_code"`
	Message    string      `json:"message"`
	Error      string      `json:"error"`
	Data       any         `json:"data"`
}
