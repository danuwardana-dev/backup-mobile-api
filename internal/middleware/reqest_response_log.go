package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"log/slog"
	"os"
	"time"
)

func CustomBodyLogger() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		slog.SetDefault(logger)
		var reqMap map[string]interface{}
		if err := json.Unmarshal(reqBody, &reqMap); err != nil {
			reqMap = map[string]interface{}{"raw": string(reqBody)}
		}

		// Filter field sensitif
		sensitiveFields := []string{"pin"}
		for _, field := range sensitiveFields {
			if _, ok := reqMap[field]; ok {
				reqMap[field] = "*****"
			}
		}

		logs := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"url":       fmt.Sprintf("%s %s%s", c.Request().Method, c.Request().Host, c.Request().URL.Path),
			"method":    c.Request().Method,
			"ip":        c.RealIP(),
			"request": map[string]interface{}{
				"headers": c.Request().Header,
				"body":    reqMap,
			},
			"response": map[string]interface{}{
				"headers": c.Response().Header(),
				"body":    string(resBody),
			},
			"status": c.Response().Status,
		}

		logDump, err := json.Marshal(logs)
		if err != nil {
			slog.Error("failed to marshal log: %v", err)
			return
		}

		slog.Info("Request Log: %s", string(logDump))
	})
}
