package verihubsDto

type SendOtpBaseRequest struct {
	MSISDN      string  `json:"msisdn"`
	Otp         *string `json:"otp"`
	Template    *string `json:"template"`
	TimeLimit   int64   `json:"time_limit"`
	Challenge   *string `json:"challenge"`
	CallbackUrl *string `json:"callback_url"`
}
type SendWhatsappOtpBaseRequest struct {
	MSISDN       string  `json:"msisdn"`
	Otp          *string `json:"otp"`
	Challenge    *string `json:"challenge"`
	TimeLimit    int64   `json:"time_limit"`
	LangCode     string  `json:"lang_code"`
	TemplateName string  `json:"template_name"`
	OtpLength    *string `json:"otp_length"`
	CallbackUrl  *string `json:"callback_url"`
}

type VerifyOtpBaseRequest struct {
	MSISDN    string  `json:"msisdn"`
	Otp       string  `json:"otp"`
	Challenge *string `json:"challenge"`
}

type (
	VerihubIdentityRequest struct {
		ValidateQuality bool   `json:"validate_quality,omitempty"`
		Image           string `json:"image,omitempty"`
	}
	VerifyKycSelfie struct {
		Nik         string `json:"nik,omitempty"`
		Name        string `json:"name,omitempty"`
		BirthDate   string `json:"birth_date,omitempty"`
		Email       string `json:"email,omitempty"`
		Phone       string `json:"phone,omitempty"`
		SelfiePhoto string `json:"selfie_photo,omitempty"`
		KtpPhoto    string `json:"ktp_photo,omitempty"`
		IsLiveness  bool   `json:"is_liveness"`
	}
)
