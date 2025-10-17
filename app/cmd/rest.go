package cmd

import (
	"backend-mobile-api/app/config"
	_ "backend-mobile-api/docs"
	internalMiddleware "backend-mobile-api/internal/middleware"
	"backend-mobile-api/internal/rest"
	"backend-mobile-api/model/enum"
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Start REST Server",
	Run:   restServer,
}

func init() {
	rootCmd.AddCommand(restCommand)
}
func restServer(cmd *cobra.Command, args []string) {
	excludedPaths := internalMiddleware.ExcludeURLValidation{
		Authorization: internalMiddleware.AuthorizationMiddlewarePath{
			ExcludeURL: []string{
				"/swagger/*",
				"/api/v1/users/auth/login",
				"/api/v1/users/auth/register",
				"/api/v1/users/auth/forgot-pin",
				"/api/v1/users/auth/set-pin",
				"/api/v1/users/auth/refresh-token",
				"/api/v1/users/auth/otp/send",
				"/api/v1/users/auth/otp/verify",

				"/healthcheck/liveness",
				"/healthcheck/readiness",

				"/api/internal/v1/verihubs/otp-invoker",
				"/api/internal/v1/verihubs/verify-ktp",
				"/api/internal/v1/verihubs/verify-passport",
				"/api/internal/v1/verihubs/verify-selfie",

				"/api/internal/v1/articles/insert",
				"/api/internal/v1/articles/update",
				"/api/internal/v1/articles/delete",
				"/api/internal/v1/articles/list",

				"/api/internal/v1/verify-ktp",
				"/api/internal/v1/verify-passport",
				"/api/internal/v1/verify-selfie",
			},
			AccessByRole: map[string][]enum.RolesEnum{},
		},
		MandatoryHeader: []string{
			"/swagger/*",
			"/api/internal/v1/verihubs/otp-invoker",
			"/api/internal/v1/verihubs/verify-ktp",
			"/api/internal/v1/verihubs/verify-passport",
			"/api/internal/v1/verihubs/verify-selfie",

			"/api/internal/v1/articles/insert",
			"/api/internal/v1/articles/update",
			"/api/internal/v1/articles/delete",
			"/api/internal/v1/articles/list",

			"/healthcheck/liveness",
			"/healthcheck/readiness",
		},
		ValidationSignaure: []string{
			"/swagger/*",
			"/api/internal/v1/verihubs/otp-invoker",
			"/api/internal/v1/verihubs/verify-ktp",
			"/api/internal/v1/verihubs/verify-passport",
			"/api/internal/v1/verihubs/verify-selfie",

			"/api/internal/v1/articles/insert",
			"/api/internal/v1/articles/update",
			"/api/internal/v1/articles/delete",
			"/api/internal/v1/articles/list",

			"/healthcheck/liveness",
			"/healthcheck/readiness",
		},
		ValidationXNonce: []string{
			"/swagger/*",
			"/api/internal/v1/verihubs/otp-invoker",
			"/api/internal/v1/verihubs/verify-ktp",
			"/api/internal/v1/verihubs/verify-passport",
			"/api/internal/v1/verihubs/verify-selfie",

			"/api/internal/v1/articles/insert",
			"/api/internal/v1/articles/update",
			"/api/internal/v1/articles/delete",
			"/api/internal/v1/articles/list",

			"/healthcheck/liveness",
			"/healthcheck/readiness",
		},
	}

	props := config.LoadForServer(EnvFilePath)
	e := echo.New()
	e.Use(middleware.Recover())
	//cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Menu-Slug, X-Origin-Path, X-Request-Id,VerificationCode,XMLHttpRequest"},
		AllowMethods:     []string{"POST, HEAD, PATCH, OPTIONS, GET, PUT"},
		AllowCredentials: true,
	}))

	e.Use((middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
		RequestIDHandler: func(c echo.Context, requestID string) {
			c.Set(echo.HeaderXRequestID, requestID)
			ctx := context.WithValue(c.Request().Context(), enum.HEADER_REQUEST_ID, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
		},
	})))
	r := e.Group("")
	rest.RouthInit(r, &controller, customMiddlewareService)
	rest.InitHealthcheckHandler(r, healtCheckController)

	routeList := internalMiddleware.ListRouth{}
	for _, route := range e.Routes() {
		routeList[route.Method+route.Path] = true
	}
	e.Use(customMiddlewareService.AccessMiddleware(&excludedPaths, &routeList))
	log.Infof("Starting server on port %s", props.Port)
	err := e.Start(props.Port)
	if err != nil {
		if err != nil {
			log.Fatalf("Error starting server on %s: %v", props.Port, err)
		}
		os.Exit(1)
	}
}
