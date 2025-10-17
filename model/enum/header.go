package enum

// middleware
type HeaderEnum string

const (
	HEADER_X_NONCE        HeaderEnum = "X-NONCE"
	HEADER_REQUEST_ID     HeaderEnum = "Request-Id"
	HEADER_X_REAL_IP      HeaderEnum = "Client-Ip"
	HEADER_AUTHORIZED     HeaderEnum = "Authorization"
	HEADER_USER_AGENT     HeaderEnum = "User-Agent"
	HEADER_X_TIMESTAMP    HeaderEnum = "X-TIMESTAMP"
	HEADER_X_SIGNATURE    HeaderEnum = "X-SIGNATURE"
	HEADER_X_DEVICE_ID    HeaderEnum = "X-DEVICE-ID"
	HEADER_X_LATITUDE     HeaderEnum = "X-LATITUDE"
	HEADER_X_LONGITUDE    HeaderEnum = "X-LONGITUDE"
	HEADER_X_API_KEY      HeaderEnum = "X-API-KEY"
	HEADER_HOST           HeaderEnum = "host"
	HEADER_PATH           HeaderEnum = "path"
	HEADER_METHOD_REQUEST HeaderEnum = "Method"
)
