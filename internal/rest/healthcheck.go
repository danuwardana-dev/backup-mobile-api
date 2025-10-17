package rest

import (
	_ "backend-mobile-api/docs"
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto"
	_ "backend-mobile-api/model/dto/swagger"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	"github.com/labstack/gommon/log"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type healthCheckHandler struct {
	CLog     *helpers.CustomLogger
	masterDb *gorm.DB
	redis    *redis.Client
	minio    *minio.Client
}

func NewHealtCheckHandler(CLog *helpers.CustomLogger,
	masterDb *gorm.DB,
	redis *redis.Client,
	minio *minio.Client,
) HealthCheckHandler {
	return &healthCheckHandler{CLog: CLog, masterDb: masterDb, redis: redis, minio: minio}
}

type HealthCheckHandler interface {
	checkLiveness(echo.Context) error
	checkReadiness(echo.Context) error
}

func InitHealthcheckHandler(e *echo.Group, h HealthCheckHandler) {

	healt := e.Group("/healthcheck")
	healt.GET("/liveness", h.checkLiveness)
	healt.GET("/readiness", h.checkReadiness)
}

// checkLiveness godoc
// @Summary Check server liveness
// @Description Returns server active status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse
// @Router /healthcheck/liveness [get]
func (h *healthCheckHandler) checkLiveness(c echo.Context) (err error) {
	customlogger, ok := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !ok {
		log.Warn("failed to get custom logger")
	}
	customlogger.Remarks = "check liveness"
	customlogger.Success = true
	return c.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: "00",
		Message:    "server is active",
		Error:      "",
		Data:       nil,
	})
}

// checkLiveness godoc
// @Summary Check server readiness
// @Description Returns database, Redis, and Minio connection status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse
// @Router /healthcheck/readiness [get]
func (h *healthCheckHandler) checkReadiness(c echo.Context) (err error) {

	var (
		status   = http.StatusOK
		failed   = "FAILED"
		ok       = "OK"
		msrDb    = "master-db"
		redisSVR = "redis-server"
		minio    = "minio"
	)

	customlogger, okData := c.Request().Context().Value(enum.CUSTOM_LOG_DATA).(*dto.CustomLoggerRequest)
	if !okData {
		log.Warn("failed to get custom logger")
	}
	customlogger.Remarks = "check readiness"
	customlogger.Success = true

	responseData := make(map[string]string)
	sqlDB, err := h.masterDb.DB()
	if err != nil {
		h.CLog.ErrorLogger(c.Request().Context(), "Error database connection", err)
		responseData[msrDb] = failed
		status = http.StatusServiceUnavailable
	}
	responseData[msrDb] = ok
	if err = sqlDB.Ping(); err != nil {
		h.CLog.ErrorLogger(c.Request().Context(), "Error database connection", err)
		responseData[msrDb] = failed
		status = http.StatusServiceUnavailable
	}
	responseData[minio] = ok
	if _, err = h.minio.ListBuckets(c.Request().Context()); err != nil {
		h.CLog.ErrorLogger(c.Request().Context(), "Error minio connection", err)
		responseData[minio] = failed
	}
	responseData[redisSVR] = ok
	if h.redis.Ping(c.Request().Context()).Err() != nil {
		h.CLog.ErrorLogger(c.Request().Context(), "Error redis connection", err)
		responseData[redisSVR] = failed
		status = http.StatusServiceUnavailable
	}

	switch status {
	case http.StatusServiceUnavailable:
		return c.JSON(http.StatusServiceUnavailable, dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "something went wrong",
			Data:       responseData,
		})
	case http.StatusOK:
		return c.JSON(http.StatusOK, dto.BaseResponse{
			StatusCode: pkgErr.SUCCESS_CODE,
			Message:    pkgErr.SUCCES_MSG,
			Data:       responseData,
		})
	}
	return c.JSON(http.StatusInternalServerError, dto.BaseResponse{
		StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
		Message:    pkgErr.SERVER_BUSY,
		Error:      "",
		Data:       nil,
	})

}
