package enum

type OtpType string

const (
	TYPE_EMAIL    OtpType = "EMAIL"
	TYPE_SMS      OtpType = "SMS"
	TYPE_WHATSAPP OtpType = "WHATSAPP"
)

type UserStatus string

const (
	VERIFICATION_STATUS_UNVERIFIED UserStatus = "UNVERIFIED"
	VERIFICATION_STATUS_VERIFIED   UserStatus = "VERIFIED"
	USER_INACTIVE                  UserStatus = "IN_ACTIVE"
)

type UserDetailStatus string

const (
	USER_WAITING_VERIFICATION UserDetailStatus = "WAITING_VERIFICATION"
	USER_WAITING_KYC_PROCESS  UserDetailStatus = "WAITING_KYC_PROCESS"
)

type RolesEnum string

const (
	ROLE_USER RolesEnum = "USER"
)

type authContext string

const (
	AUTH_JWT        authContext = "AUTH_JWT"
	AUTH_UUID       authContext = "AUTH_UUID"
	AUTH_USER_ID    authContext = "AUTH_USER_ID"
	AUTH_USER_EMAIL authContext = "AUTH_USER_EMAIL"
	Auth_ROLE       authContext = "AUTH_ROLE"
	USER
)

type LoginType string

const (
	LOGIN_EMAIL        LoginType = "EMAIL"
	LOGIN_PHONE_NUMBER LoginType = "PHONE_NUMBER"
	LOGIN_BIOMETRIC    LoginType = "BIOMETRIC"
)

type LoginStatus string

const (
	LOGIN_SUCCESS LoginStatus = "SUCCESS"
	LOGIN_FAILED  LoginStatus = "FAILED"
)

type BiometricStatus string

const (
	BIOMETRIC_ACTIVE    BiometricStatus = "ACTIVE"
	BIOMETRIC_IN_ACTIVE BiometricStatus = "IN_ACTIVE"
)

type UserKYCStatus string

const (
	USER_KYC_STATUS_UNKNOWN UserKYCStatus = "UNKNOWN"
)

type UserKYCType string

const (
	USER_KYC_KTP      UserKYCType = "KTP"
	USER_KYC_PASSPORT UserKYCType = "PASSPORT"
)

type AccessTokenType string

const (
	ACCESS_SET_PIN            AccessTokenType = "SET_PIN"
	ACCESS_FORGOT_PIN         AccessTokenType = "FORGOT_PIN"
	ACCESS_RESET_PIN          AccessTokenType = "RESET_PIN"
	ACCESS_RESET_EMAIL        AccessTokenType = "RESET_EMAIL"
	ACCESS_RESET_PHONE_NUMBER AccessTokenType = "RESET_PHONE_NUMBER"
	ACCESS_DELETE_ACCOUNT     AccessTokenType = "DELETE_ACCOUNT"
)
