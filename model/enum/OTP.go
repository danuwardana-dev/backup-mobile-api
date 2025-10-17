package enum

type VerihubsEnum string

const (
	VERIHUBS_KEY    VerihubsEnum = "API-Key"
	VERIHUBS_APP_ID VerihubsEnum = "App-ID"
)

type OTPStatus string

const (
	OTP_REQUESTED           OTPStatus = "Requested"
	OTP_DELIVERED           OTPStatus = "Delivered"
	OTP_VERIFIED            OTPStatus = "Verified"
	OTP_NOT_VERIFIED        OTPStatus = "Not Verified"
	OTP_FAILED              OTPStatus = "Failed"
	OTP_REQUEST_ERROR       OTPStatus = "Request Error"
	OTP_REJECTED            OTPStatus = "Rejected"
	OTP_UNDELIVERED         OTPStatus = "Undelivered"
	OTP_NO_DELIVERY_REPORT  OTPStatus = "No Delivery Report"
	OTP_BLOCKED             OTPStatus = "Blocked"
	OTP_SENT                OTPStatus = "Sent"
	OTP_READ                OTPStatus = "Read"
	OTP_TIER_LIMIT_EXCEEDED OTPStatus = "Tier Limit Exceeded"
	OTP_UNVERIFIED          OTPStatus = "Unverified"
)

var OtpSMSMapStatusVerihubs = map[int]OTPStatus{
	0: OTP_REQUESTED,
	1: OTP_DELIVERED,
	2: OTP_VERIFIED,
	3: OTP_NOT_VERIFIED,
	4: OTP_FAILED,
	5: OTP_REQUEST_ERROR,
	6: OTP_REJECTED,
	7: OTP_UNDELIVERED,
	8: OTP_NO_DELIVERY_REPORT,
	9: OTP_BLOCKED,
}

// 0	Requested	SMS OTP has been requested.
// 1	Delivered	SMS OTP has been delivered to destination number. This status will be charged.
// 2	Verified	SMS OTP has been verified using verify OTP API. This status will be charged.
// 3	Not Verified	SMS OTP has not been verified after the time limit. This status will be charged.
// 4	Failed	Operator cannot send SMS OTP to destination number
// 5	Request Error	Verihubs cannot reach Operator
// 6	Rejected	Request rejected due to Too Many Requests
// 7	Undelivered	SMS OTP has not been delivered to destination number (response received from Operator). This status will be charged.
// 8	No Delivery Report	Delivered without Delivery Report from Operator after certain period. This status will be charged.
// 9	Blocked	Destination Number has been blocked by Verihubs's system
var OtpWhatsappMapStatusVerihubs = map[int]OTPStatus{
	0:  OTP_REQUESTED,
	1:  OTP_SENT,
	2:  OTP_DELIVERED,
	3:  OTP_READ,
	4:  OTP_VERIFIED,
	5:  OTP_UNVERIFIED,
	6:  OTP_FAILED,
	7:  OTP_REQUEST_ERROR,
	10: OTP_TIER_LIMIT_EXCEEDED,
}

//Status	Description	Condition
//0	Requested	OTP has been requested. This status will be charged.
//1	Sent	OTP has been sent to Whatsapp server. This status will be charged.
//2	Delivered	OTP has been delivered to destination number. This status will be charged.
//3	Read	OTP has been read by user. This status will be charged.
//4	Verified	OTP has been verified by user. This status will be charged.
//5	Unverified	OTP has reached time limit. This status will be charged.
//6	Failed	Whatsapp cannot send OTP to destination number. This status will be charged.
//7	Request Error	Verihubs cannot reach WhatsApp
//10	Tier Limit Exceeded	Exceeding channel tier limit

type OtpService string

const (
	OTP_VERIFY_ACCOUNT     OtpService = "OTP_VERIFY_ACCOUNT"
	OTP_RESET_PIN          OtpService = "OTP_RESET_PIN"
	OTP_FORGOT_PIN         OtpService = "OTP_FORGOT_PIN"
	OTP_RESET_EMAIL        OtpService = "OTP_RESET_EMAIL"
	OTP_RESET_PHONE_NUMBER OtpService = "OTP_RESET_PHONE_NUMBER"
)

type RedisOtpTag string

const (
	OTP_ACCESS RedisOtpTag = "OTP_ACCESS"
)
