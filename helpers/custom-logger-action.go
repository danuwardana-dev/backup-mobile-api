package helpers

import (
	"backend-mobile-api/model/dto"
	customMiddleware "backend-mobile-api/model/dto"
	"backend-mobile-api/model/enum"
	"context"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"log/slog"
	"os"
)

func CustomeLogger(ctx context.Context, req *dto.CustomLoggerRequest) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	var errString string
	slog.SetDefault(logger)
	if ctx.Err() != nil {
		errString = ctx.Err().Error()
	}

	customResource, ok := ctx.Value(enum.CUSTOM_CONTEXT_VALUE).(*customMiddleware.ContextValue)
	if !ok {
		log.Error("failed to get custom resource")
		return
	}
	logData := dto.HandlerLog{
		RequestId: customResource.HeaderRequestId,
		Timestamp: customResource.HeaderXTimestamp,
		Url:       customResource.HeaderHost + customResource.HeaderPath,
		Method:    customResource.HeaderMethod,
		Device: dto.Device{
			DeviceID:  customResource.HeaderXDeviceID,
			Longitude: customResource.HeaderXLongitude,
			Latitude:  customResource.HeaderXLatitude,
			Ip:        customResource.HeaderXRealIp,
		},
		Success:  false,
		Error:    errString,
		UserUUID: customResource.AuthUUID,
		Email:    customResource.AuthEmail,
	}

	if req != nil {
		if logData.UserUUID == "" {
			logData.UserUUID = req.UserUUID
		}
		if logData.Email == "" {
			logData.Email = req.Email
		}
		if req.Error != "" {
			logData.Error = req.Error
		}
		logData.Remarks = req.Remarks
		logData.Data = req.Data
		logData.Success = req.Success
	}
	jsonLog, _ := json.Marshal(logData)
	slog.Info(string(jsonLog))

}
