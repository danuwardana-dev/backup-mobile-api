package cmd

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/docs"
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/middleware"
	"backend-mobile-api/internal/outbond/smtp"
	"backend-mobile-api/internal/outbond/verihubs"
	"backend-mobile-api/internal/repository/minio"
	"backend-mobile-api/internal/repository/postgres"
	redisRepos "backend-mobile-api/internal/repository/redis"
	"backend-mobile-api/internal/rest"
	articleController "backend-mobile-api/internal/rest/article-controller"
	bankListController "backend-mobile-api/internal/rest/bank-list-controller"
	checkAccountBankController "backend-mobile-api/internal/rest/check-account-bank-controller"
	ppobListController "backend-mobile-api/internal/rest/ppob-list-controller"
	recipientController "backend-mobile-api/internal/rest/recipient-controller"
	transactionController "backend-mobile-api/internal/rest/transactions-controller"
	userPaymentAccountController "backend-mobile-api/internal/rest/user-account-payments-controller"
	recipientSvc "backend-mobile-api/service/recipient-svc"
	transactionsvc "backend-mobile-api/service/transactions-svc"

	kyccontroller "backend-mobile-api/internal/rest/kyc-controller"
	userAuth "backend-mobile-api/internal/rest/user-auth-controller"
	userProfileController "backend-mobile-api/internal/rest/user-profile-controller"
	verihubsInvokerController "backend-mobile-api/internal/rest/verihubs-invoker-controller"
	articleSvc "backend-mobile-api/service/article-svc"
	banklistsvc "backend-mobile-api/service/bank-list-svc"
	"backend-mobile-api/service/biometricSvc"
	kycservice "backend-mobile-api/service/kyc-service"
	"backend-mobile-api/service/notification"
	"backend-mobile-api/service/otp"
	ppoblistsvc "backend-mobile-api/service/ppob-list-svc"
	userAccountPaymentSvc "backend-mobile-api/service/user-accounts-payment-svc"
	user_auth_svc "backend-mobile-api/service/user-auth-svc"
	userProfileService "backend-mobile-api/service/user-profile-svc"
	verihubsInvokerService "backend-mobile-api/service/verihubs-invoker-service"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	EnvFilePath string
	rootCmd     = &cobra.Command{
		Use:   "cobra-cli",
		Short: "backend-mobile-api",
	}
)
var (
	rootConfig              config.Root
	MasterDatabase          *gorm.DB
	RedisClient             *redis.Client
	err                     error
	logger                  = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	CLoger                  *helpers.CustomLogger
	controller              rest.Controller
	customMiddlewareService middleware.CustomMiddleware
	healtCheckController    rest.HealthCheckHandler
)

func init() {
	cobra.OnInitialize(func() {
		initCustomLoger()
		initConfigReader()
		initRedisClient()
		initPostgres()
		initApp()
		initSwagger()

	})
}
func initConfigReader() {
	logger.Info("Loading config from environment")
	rootConfig = config.Load(EnvFilePath)
}
func initPostgres() {
	logger.Info("Loading postgres")
	postgres := config.LoadPostgres(rootConfig.Postgres)
	MasterDatabase, err = postgres.OpenPostgresDatabaseConnection()
	if err != nil {
		log.Infof("Error loading postgres: %s", err.Error())
		os.Exit(1)
	}
}
func initRedisClient() {
	logger.Info("Loading redis client")
	redisSet := config.LoadRedis(rootConfig.Redis)
	RedisClient, err = redisSet.RedisClient()
	if err != nil {
		log.Infof("Error loading redis client : %s", err.Error())
		os.Exit(1)
	}
}
func initCustomLoger() {
	logger.Info("Loading custom loger")
	CLoger = helpers.NewLogger(logger)
}
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&EnvFilePath, "env", "e", ".env", ".env file to read from")
	if err := rootCmd.Execute(); err != nil {
		slog.Error("E err : %s", err.Error())
		os.Exit(1)
	}

}
func initSwagger() {
	docs.SwaggerInfo.Title = "BEYOND MOBILE-API"
	docs.SwaggerInfo.Description = "Dokumentasi API untuk sistem backend mobile BEYOND"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
func initApp() {
	time.Local, _ = time.LoadLocation(rootConfig.App.TimeZone)
	//repository
	userRepository := postgres.NewUserRepository(MasterDatabase, CLoger)
	redisRepository := redisRepos.NewRedis(RedisClient, &rootConfig)
	accessStateRepsoitory := postgres.NewResetPinRepository(MasterDatabase, CLoger)
	otpRepository := postgres.NewOtpRepository(MasterDatabase, CLoger)
	tokenBlacklistRepository := postgres.NewTokenBlacklistTokenRepository(MasterDatabase, CLoger)
	userDetilRepository := postgres.NewUserDetailRepository(MasterDatabase, CLoger)
	deviceRepository := postgres.NewDeviceRepository(MasterDatabase, CLoger)
	articleRepository := postgres.NewArticleRepository(MasterDatabase, CLoger)
	bankRepo := postgres.NewBankListRepository(MasterDatabase, CLoger) // kalau pakai *gorm.DB, ambil DB() biar dapat *sql.DB
	ppobRepo := postgres.NewPpobListRepository(MasterDatabase, CLoger)
	userPaymentAccountRepo := postgres.NewUserPaymentsAccountRepository(MasterDatabase, CLoger)
	// service
	ppobListService := ppoblistsvc.NewPpobListService(ppobRepo)
	bankService := banklistsvc.NewBankListService(bankRepo)
	userAccountsPaymentService := userAccountPaymentSvc.NewUserAccountsPaymentService(userPaymentAccountRepo)
	// === Recipient ===
	// repository
	recipientRepo := postgres.NewRecipientRepository(MasterDatabase, CLoger)
	// service
	recipientService := recipientSvc.NewRecipientService(recipientRepo)
	// controller
	controller.RecipientController = recipientController.NewRecipientController(recipientService)
	controller.CheckAccountBankController = checkAccountBankController.NewCheckAccountBankController()
	// controller
	controller.UserAccountPaymentsController = userPaymentAccountController.NewUserAccountPaymentsController(userAccountsPaymentService)
	controller.PpobListController = ppobListController.NewPpobListController(ppobListService)
	controller.BankListController = bankListController.NewBankListController(bankService)
	// === Transaction ===
	firebaseNotifier, err := notification.InitFirebaseNotifier(
		context.Background(),
		rootConfig.Firebase.CredentialsFile, // pastikan ada di config
	)
	if err != nil {
		panic(err)
	}
	smtp := smtp.NewSmtp(&rootConfig, CLoger)
	transactionRepo := postgres.NewTransactionRepository(MasterDatabase, CLoger)
	transactionService := transactionsvc.NewTransactionService(transactionRepo, firebaseNotifier, smtp, userRepository)
	controller.TransactionController = transactionController.NewTransactionController(transactionService)

	minioClient, err := rootConfig.Minio.MinioClientSet()
	if err != nil {
		panic(err)
	}
	minioRepository := minio.NewMinioRepository(minioClient, &rootConfig, rootConfig.Minio.Bucket, CLoger)
	ktpRepository := postgres.NewKycKtpRepository(MasterDatabase, CLoger)
	passportRepository := postgres.NewKycPassportRepository(MasterDatabase, CLoger)

	//middleware
	customMiddlewareService = middleware.NewCustomMiddleware(&rootConfig.Jwt, CLoger, *redisRepository, &rootConfig)
	//xsesionMiddleware = middleware.NewXsesionMiddleware(&rootConfig, CLoger, *redisRepository)

	//outbound

	outboundVeriHubsSvc := verihubs.NewOutboundVeriHubsService(&rootConfig.Verihubs, &rootConfig, CLoger)

	//controller
	healtCheckController = rest.NewHealtCheckHandler(CLoger, MasterDatabase, RedisClient, minioClient)
	controller.UserAuthController = userAuth.NewUserAuthController(
		user_auth_svc.NewUserAuthService(
			userRepository,
			customMiddlewareService,
			*redisRepository,
			rootConfig,
			smtp,
			CLoger,
			accessStateRepsoitory,
			otpRepository,
			tokenBlacklistRepository,
			outboundVeriHubsSvc,
			userDetilRepository,
			deviceRepository,
		),
		biometricSvc.NewBiometricService(
			userRepository,
			userDetilRepository,
			customMiddlewareService,
			tokenBlacklistRepository,
		),
	)
	controller.VerihubsInvoker = verihubsInvokerController.NewVerihubsInvokerController(
		verihubsInvokerService.NewVerihubsInvokerService(
			otpRepository, CLoger,
		),
	)
	controller.UserProfileController = userProfileController.NewUserProfileController(
		userProfileService.NewUserProfileService(
			userRepository,
			redisRepository,
			rootConfig,
			smtp,
			CLoger,
			accessStateRepsoitory,
			otpRepository,
			outboundVeriHubsSvc,
			userDetilRepository,
			otp.NewOtpService(
				otpRepository,
				&rootConfig,
				redisRepository,
				CLoger,
				outboundVeriHubsSvc,
				smtp,
			),
			minioRepository,
		),
	)
	controller.ArticleController = articleController.NewArticleController(
		articleSvc.NewArticleService(
			userRepository,
			articleRepository,
			CLoger,
		),
	)

	controller.KycController = kyccontroller.NewKycController(
		kycservice.NewKycService(
			CLoger, outboundVeriHubsSvc, passportRepository, ktpRepository, userRepository, userDetilRepository, minioRepository),
	)

}
