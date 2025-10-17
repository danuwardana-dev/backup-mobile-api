package config

type Verihubs struct {
	Domain          string `envconfig:"VERIHUBS_DOMAIN" required:"true"`
	AppID           string `envconfig:"VERIHUBS_APP_ID" required:"true"`
	VerihubsKey     string `envconfig:"VERIHUBS_KEY" required:"true"`
	VerihubsOTPCode bool   `envconfig:"VERIHUBS_OTP_CODE"`

	OTPCallBackUrl  *string `envconfig:"VERIHUBS_OTP_CALLBACK_URL"`
	OTPSMSChallenge *string `envconfig:"VERIHUBS_OTP_SMS_CHALLENGE"`
	OTPSMSTemplate  *string `envconfig:"VERIHUBS_OTP_SMS_TEMPLATE"`

	OTPWhatsappChallenge    *string `envconfig:"VERIHUBS_OTP_WHATSAPP_CHALLENGE"`
	OTPWhatsappLangCode     string  `envconfig:"VERIHUBS_OTP_WHATSAPP_LANG_CODE" required:"true"`
	OTPWhatsappTemplateName string  `envconfig:"VERIHUBS_OTP_WHATSAPP_TEMPLATE" required:"true"`
	OTPWhatsappOtpLength    *string `envconfig:"VERIHUBS_OTP_WHATSAPP_OTP_LENGTH"`
}
