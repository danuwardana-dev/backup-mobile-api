package dto

type ContextValue struct {
	HeaderContentType string `header:"Content-Type" json:"Content-Type"`
	HeaderUserAgent   string `header:"User-Agent" json:"User-Agent" validate:"required"`
	HeaderXTimestamp  string `header:"X-TIMESTAMP" json:"X-TIMESTAMP" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	HeaderXSignature  string `header:"X-SIGNATURE" json:"X-SIGNATURE" validate:"required"`
	HeaderXRealIp     string `header:"X-Real-Ip" json:"X-Real-Ip" validate:"required"`
	HeaderXNonce      string `header:"X-NONCE" json:"X-NONCE" validate:"required"`
	HeaderXDeviceID   string `header:"X-DEVICE-ID" json:"X-DEVICE-ID" validate:"required"`
	HeaderXLatitude   string `header:"X-LATITUDE" json:"X-LATITUDE" validate:"required"`
	HeaderXLongitude  string `header:"X-LONGITUDE" json:"X-LONGITUDE" validate:"required"`
	HeaderXApiKey     string `header:"X-API-KEY" json:"X-API-KEY"`

	HeaderRequestId     string
	HeaderHost          string
	HeaderPath          string
	RequestPath         string
	HeaderMethod        string
	HeaderAuthorization string

	AuthUUID     string
	AuthEmail    string
	AuthDeviceID string
	AuthRole     string
}
