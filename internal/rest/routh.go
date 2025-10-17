package rest

import (
	customMiddleware "backend-mobile-api/internal/middleware"
	articleController "backend-mobile-api/internal/rest/article-controller"
	bankListController "backend-mobile-api/internal/rest/bank-list-controller"
	checkaccountbankcontroller "backend-mobile-api/internal/rest/check-account-bank-controller"
	kycCtr "backend-mobile-api/internal/rest/kyc-controller"
	ppobListController "backend-mobile-api/internal/rest/ppob-list-controller"
	recipientController "backend-mobile-api/internal/rest/recipient-controller"
	transactionController "backend-mobile-api/internal/rest/transactions-controller"
	userPaymentAccountController "backend-mobile-api/internal/rest/user-account-payments-controller"
	userAuthCtr "backend-mobile-api/internal/rest/user-auth-controller"
	userProfileController "backend-mobile-api/internal/rest/user-profile-controller"
	verihubsInvokerCtr "backend-mobile-api/internal/rest/verihubs-invoker-controller"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Controller struct {
	UserAuthController            userAuthCtr.UserAuthController
	KycController                 kycCtr.KYCController
	VerihubsInvoker               verihubsInvokerCtr.VerihubsInvokerController
	UserProfileController         userProfileController.UserProfileController
	ArticleController             articleController.ArticleController
	BankListController            bankListController.BankListController
	RecipientController           recipientController.RecipientController
	CheckAccountBankController    checkaccountbankcontroller.CheckAccountBankController
	PpobListController            ppobListController.PpobListController
	TransactionController         transactionController.TransactionController
	UserAccountPaymentsController userPaymentAccountController.UserAccountPaymentsController
}

func RouthInit(e *echo.Group, ctr *Controller, middlewareCustom customMiddleware.CustomMiddleware) {
	log.Info("RouthInit")
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(customMiddleware.RateLimitMiddleware())
	internalV1 := e.Group("/api/internal/v1")

	v1 := e.Group("/api/v1")
	//	v1.Use(middlewareCustom.CommonCustomHeaderMiddleware())

	users := v1.Group("/users")
	//user-auth-controller
	authRouth := users.Group("/auth")
	authRouth.POST("/register", ctr.UserAuthController.RegisterController)
	authRouth.POST("/login", ctr.UserAuthController.LoginController)
	authRouth.POST("/logout", ctr.UserAuthController.LogoutController)
	authRouth.POST("/refresh-token", ctr.UserAuthController.RefreshTokenController)
	authRouth.POST("/otp/send", ctr.UserAuthController.SendOtpController)
	authRouth.POST("/otp/verify", ctr.UserAuthController.VerifyOtpController)
	authRouth.POST("/set-pin", ctr.UserAuthController.SetPinController)
	authRouth.POST("/forgot-pin", ctr.UserAuthController.ForgotPinController)

	//verhubs
	verihubs := internalV1.Group("/verihubs")
	verihubs.GET("/otp-invoker", ctr.VerihubsInvoker.InvokerOtpController)
	verihubs.POST("/verify-ktp", ctr.KycController.VerifyKycKTP)
	verihubs.POST("/verify-passport", ctr.KycController.VerifyKycPassport)
	verihubs.POST("/verify-selfie", ctr.KycController.VerifySelfie)

	// kyc controller
	userKyc := users.Group("/kyc")
	userKyc.POST("/register-kyc-ktp", ctr.KycController.SaveKycKTP)
	userKyc.POST("/register-kyc-passport", ctr.KycController.SaveKycPassport)

	//profile
	profileRouth := users.Group("/profile")
	profileRouth.POST("/access-token", ctr.UserProfileController.AccessTokenController)
	profileRouth.GET("/inquiry", ctr.UserProfileController.InquiryUserProfileController)
	profileRouth.POST("/reset/pin", ctr.UserProfileController.ResetPinController)
	profileRouth.POST("/reset/email", ctr.UserProfileController.ResetEmailController)
	profileRouth.POST("/reset/phone-number", ctr.UserProfileController.ResetPhoneNumberController)
	profileRouth.POST("/reset/full-name", ctr.UserProfileController.ResetFullNameController)
	profileRouth.POST("/reset/profile-image", ctr.UserProfileController.ResetProfileImageController)
	profileRouth.POST("/otp/verify", ctr.UserProfileController.VerifyOtpController)
	profileRouth.POST("/biometric", ctr.UserProfileController.BiometricController)
	profileRouth.POST("/delete-account", ctr.UserProfileController.DeleteAccountController)
	profileRouth.GET("/profile-image", ctr.UserProfileController.GetProfilePictureController)

	//articles
	internalArticle := internalV1.Group("/articles")
	articles := v1.Group("/articles")
	internalArticle.POST("/insert", ctr.ArticleController.InsertNewArticleController)
	internalArticle.PATCH("/update", ctr.ArticleController.UpdateArticleController)
	internalArticle.POST("/delete", ctr.ArticleController.DeleteArticleController)
	internalArticle.POST("/list", ctr.ArticleController.InternalGetArticleController)
	articles.POST("/list", ctr.ArticleController.GetArticleController)

	//inquiry bank list
	banks := users.Group("/inquiry-bank-list")
	banks.GET("", ctr.BankListController.GetBankListController)
	//ppob list
	ppobList := users.Group("/inquiry-ppob-list")
	ppobList.GET("", ctr.PpobListController.GetPpobListController)
	// recipient
	recipient := users.Group("/recipient")
	recipient.POST("/save-recipient", ctr.RecipientController.CreateRecipient)
	recipient.GET("/inquiry-recipient", ctr.RecipientController.GetRecipients)

	//check account bank
	checkAccountBank := users.Group("/check-account")
	checkAccountBank.POST("/check", ctr.CheckAccountBankController.CheckAccount)

	// transactions
	transactions := users.Group("/transactions")
	transactions.POST("/generate", ctr.TransactionController.GenerateTransactionCode)
	transactions.POST("", ctr.TransactionController.CreateTransaction) // create transaksi
	transactions.GET("", ctr.TransactionController.GetAllTransactions)
	transactions.POST("/status-update", ctr.TransactionController.UpdateTransactionStatus)
	// get userAccountPayments

	userAccountPayment := users.Group("/user-account-payment")
	userAccountPayment.GET("", ctr.UserAccountPaymentsController.GetUserAccountPaymentsController)
}
