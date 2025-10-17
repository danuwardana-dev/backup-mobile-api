package helpers

import (
	"backend-mobile-api/model/enum"
	"context"
	"log/slog"
	"runtime"
)

type CustomLogger struct {
	logger *slog.Logger
}

func NewLogger(logger *slog.Logger) *CustomLogger {
	return &CustomLogger{logger: logger}
}

func (s *CustomLogger) runtimeCaller(call int) (pc uintptr, file string, line int, funcName string, ok bool) {
	pc, file, line, ok = runtime.Caller(call)
	funcName = "???"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
		}
	} else {
		file = "???"
		line = 0
		funcName = "???"
	}
	return pc, file, line, funcName, ok
}
func (s *CustomLogger) ErrorLogger(ctx context.Context, msg string, err error) {
	var (
		reqId string
		atrr  []any
	)
	if Uuid := ctx.Value(enum.HEADER_REQUEST_ID); Uuid == nil {
		reqId = "xxxxx"
	} else {
		reqId = Uuid.(string)
	}
	_, file, line, _, ok := s.runtimeCaller(2)
	if !ok {
		file = "???"
		line = 0
	}
	atrr = append(atrr,
		slog.String("requestID", reqId),
		slog.String("error", err.Error()),
		slog.String("file", file),
		slog.Int("line", line))

	s.logger.Error(msg, atrr...)

}
func (s *CustomLogger) InfoLogger(ctx context.Context, msg string) {
	var (
		reqId string
		atrr  []any
	)
	if Uuid := ctx.Value("requestID"); Uuid == nil {
		reqId = "xxxxx"
	} else {
		reqId = Uuid.(string)
	}
	_, file, line, _, ok := s.runtimeCaller(2)
	if !ok {
		file = "???"
		line = 0
	}
	atrr = append(atrr,
		slog.String("requestID", reqId),
		slog.String("file", file),
		slog.Int("line", line))

	s.logger.Info(msg, atrr...)
}
func (s *CustomLogger) WarnLogger(ctx context.Context, msg string) {
	var (
		reqId string
		atrr  []any
	)
	if Uuid := ctx.Value("requestID"); Uuid == nil {
		reqId = "xxxxx"
	} else {
		reqId = Uuid.(string)
	}
	_, file, line, _, ok := s.runtimeCaller(2)
	if !ok {
		file = "???"
		line = 0
	}
	atrr = append(atrr,
		slog.String("requestID", reqId),
		slog.String("file", file),
		slog.Int("line", line))

	s.logger.Warn(msg, atrr...)
}
