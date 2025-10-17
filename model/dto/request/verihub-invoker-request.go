package request

type VerihubsOtpInvoker struct {
	Method    string `param:"method" json:"method" validate:"required"`
	SessionId string `query:"session_id" validate:"required" json:"session_id"`
	Status    string `query:"status" validate:"required,numeric" json:"status"`
}
