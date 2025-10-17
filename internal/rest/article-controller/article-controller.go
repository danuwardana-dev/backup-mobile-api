package articleController

import (
	_ "backend-mobile-api/docs"
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	_ "backend-mobile-api/model/dto/swagger"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	articleSvc "backend-mobile-api/service/article-svc"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type articleController struct {
	articleService articleSvc.ArticleService
}

func NewArticleController(articleService articleSvc.ArticleService) ArticleController {
	return &articleController{articleService: articleService}
}

type ArticleController interface {
	InsertNewArticleController(e echo.Context) error
	UpdateArticleController(e echo.Context) error
	DeleteArticleController(e echo.Context) error
	GetArticleController(e echo.Context) error
	InternalGetArticleController(e echo.Context) error
}

// @Tags Articles
// @Summary New Articles
// @Description Add New Articles
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.NewArticleRequest true "Insert Articles Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 500 {object} swagger.CommonError
// @Router /api/internal/v1/articles/insert [post]
func (ctr *articleController) InsertNewArticleController(e echo.Context) error {
	var (
		req      request.NewArticleRequest
		validate = validator.New()
	)
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "Insert New Article"

	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{StatusCode: "", Message: pkgErr.INVALID_REQUEST_PAYLOAD_MSG})
	}
	err = validate.Struct(req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{
			StatusCode: "",
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.articleService.InsertArticleService(e.Request().Context(), &req, nil, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Articles
// @Summary Update Articles
// @Description Update Existing Articles
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.UpdateArticleRequest true "Update Articles Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.RecordNotFoundArticles
// @Failure 500 {object} swagger.CommonError
// @Router /api/internal/v1/articles/update [patch]
func (ctr *articleController) UpdateArticleController(e echo.Context) error {
	var (
		req      request.UpdateArticleRequest
		validate = validator.New()
	)
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "Article-update"

	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{StatusCode: "", Message: pkgErr.INVALID_REQUEST_PAYLOAD_MSG})
	}
	err = validate.Struct(req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{
			StatusCode: "",
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.articleService.UpdateArticleService(e.Request().Context(), &req, nil, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.ARTICLE_RECORD_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Articles
// @Summary Remove Articles
// @Description Remove New Articles
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.DeleteArticleRequest true "Remove Articles Request"
// @Success 200 {object} swagger.BasicSuccess
// @Failure 404 {object} swagger.RecordNotFoundArticles
// @Failure 500 {object} swagger.CommonError
// @Router /api/internal/v1/articles/delete [post]
func (ctr *articleController) DeleteArticleController(e echo.Context) error {
	var (
		req      request.DeleteArticleRequest
		validate = validator.New()
	)
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "Article-delete"

	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{StatusCode: "", Message: pkgErr.INVALID_REQUEST_PAYLOAD_MSG})
	}
	err = validate.Struct(req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{
			StatusCode: "",
			Message:    pkgErr.INVALID_REQUEST_PAYLOAD_MSG,
			Error:      err.Error(),
		})
	}
	res := ctr.articleService.DeleteArticleService(e.Request().Context(), &req, nil, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.ARTICLE_RECORD_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Articles
// @Summary Lit Articles
// @Description Inquiry List Articles
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param Authorization header string true "Authorization"
// @Param data body request.SelectListArticleRequest true "List Articles Request"
// @Success 200 {object} swagger.SuccessListArticles
// @Failure 403 {object} swagger.DeferenceDeviceArtice
// @Failure 404 {object} swagger.RecordNotFoundArticles
// @Failure 500 {object} swagger.CommonError
// @Router /api/v1/articles/list [post]
func (ctr articleController) GetArticleController(e echo.Context) error {
	var (
		req      request.SelectListArticleRequest
		validate = validator.New()
		userUUID string
	)
	customResource, ok := e.Request().Context().Value(enum.CUSTOM_CONTEXT_VALUE).(*dto.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")

		return e.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "failed to get custom resource",
			Data:       nil,
		})
	}
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "Article-list"

	userUUID = customResource.AuthUUID
	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{StatusCode: "", Message: pkgErr.INVALID_REQUEST_PAYLOAD_MSG})
	}
	err = validate.Struct(req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
	}
	res := ctr.articleService.SelectArticleListService(e.Request().Context(), &req, &userUUID, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.ARTICLE_DEFERENCE_DEVICE_CODE:
		return e.JSON(http.StatusForbidden, res)
	case pkgErr.ARTICLE_RECORD_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}

// @Tags Articles
// @Summary Lit Articles
// @Description Inquiry List Articles
// @Accept json
// @Produce json
// @Param X-NONCE header string true "X-NONCE"
// @Param X-SIGNATURE header string true "X-SIGNATURE"
// @Param X-DEVICE-ID header string true "X-DEVICE-ID"
// @Param X-TIMESTAMP header string true "X-TIMESTAMP"
// @Param X-LATITUDE header string true "X-LATITUDE"
// @Param X-LONGITUDE header string true "X-LONGITUDE"
// @Param data body request.SelectListArticleRequest true "Select List Articles Request"
// @Success 200 {object} swagger.SuccessListArticles
// @Failure 404 {object} swagger.RecordNotFoundArticles
// @Failure 500 {object} swagger.CommonError
// @Router /api/internal/v1/articles/list [post]
func (ctr articleController) InternalGetArticleController(e echo.Context) error {
	var (
		req      request.SelectListArticleRequest
		validate = validator.New()
	)
	logData, okData := e.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	logData.Remarks = "Article-internal-list"

	err := e.Bind(&req)
	if err != nil {
		logData.Error = err.Error()
		return e.JSON(http.StatusBadRequest, &dto.BaseResponse{StatusCode: "", Message: pkgErr.INVALID_REQUEST_PAYLOAD_MSG})
	}
	err = validate.Struct(req)
	if err != nil {
		err = helpers.CustomValidatePayload(err, req)
	}
	res := ctr.articleService.InternalSelectArticleListService(e.Request().Context(), &req, nil, logData)
	switch res.StatusCode {
	case pkgErr.SUCCESS_CODE:
		return e.JSON(http.StatusOK, res)
	case pkgErr.ARTICLE_RECORD_NOT_FOUND_CODE:
		return e.JSON(http.StatusNotFound, res)
	}
	return e.JSON(http.StatusInternalServerError, res)
}
