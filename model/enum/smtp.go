package enum

type EmailSubject string

var VERIFY_OTP_SUBJECT EmailSubject = "Verify Your Account – OTP Code"
var ACCESS_RESET_PIN_SUBJECT EmailSubject = "Forgot PIN – Request for Assistance"
