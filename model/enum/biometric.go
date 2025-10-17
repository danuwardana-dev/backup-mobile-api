package enum

type BiometricType string

const (
	BIOMETRIC_LOGIN              BiometricType = "LOGIN"
	BIOMETRIC_TRANSACTION_VERIFY BiometricType = "SIGNUP"
)
