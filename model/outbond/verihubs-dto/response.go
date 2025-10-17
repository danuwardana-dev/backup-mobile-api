package verihubsDto

type SendOtpSMSResponse struct {
	Message      string `json:"message"`
	Otp          string `json:"otp"`
	MSISDN       string `json:"msisdn"`
	SessionID    string `json:"session_id"`
	SegmentCount int    `json:"segment_count"`
}
type SendOtpWAResponse struct {
	Message   string `json:"message"`
	Otp       string `json:"otp"`
	MSISDN    string `json:"msisdn"`
	SessionID string `json:"session_id"`
	TryCount  int    `json:"try_count"`
}
type VerifyOtpBaseResponse struct {
	Message string `json:"message"`
}

type (
	IdentityKTPResponse struct {
		Message   string              `json:"message,omitempty"`
		ErrorCode string              `json:"error_code,omitempty"`
		Data      IdentityKTPVerihubs `json:"data,omitempty"`
	}

	IdentityKTPVerihubs struct {
		FullName               string `json:"full_name,omitempty"`
		Gender                 string `json:"gender,omitempty"`
		Address                string `json:"address,omitempty"`
		Administrative_village string `json:"administrative_village,omitempty"`
		BloodType              string `json:"blood_type,omitempty"`
		City                   string `json:"city,omitempty"`
		DateOfBirth            string `json:"date_of_birth,omitempty"`
		District               string `json:"district,omitempty"`
		MartialStatus          string `json:"martial_status,omitempty"`
		Nationality            string `json:"nationality,omitempty"`
		Nik                    string `json:"nik,omitempty"`
		Occupation             string `json:"occupation,omitempty"`
		PlaceOfBirth           string `json:"place_of_birth,omitempty"`
		Religion               string `json:"religion,omitempty"`
		RtRw                   string `json:"rt_rw,omitempty"`
		State                  string `json:"state,omitempty"`
		ImageQuality           any    `json:"image_quality,omitempty"`
	}

	IdentityPassportResponse struct {
		Message   string                   `json:"message,omitempty"`
		ErrorCode string                   `json:"error_code,omitempty"`
		Data      IdentityPassportVerihubs `json:"data,omitempty"`
	}

	IdentityPassportVerihubs struct {
		Id           string                `json:"id,omitempty"`
		Reference_id string                `json:"reference_id,omitempty"`
		ResultData   ResultDataKycPassport `json:"result_data,omitempty"`
	}

	//VerifySelfieResponse struct {
	//	Message string      `json:"message"`
	//	Data    interface{} `json:"data"`
	//}

	VerifySelfieResponse struct {
		Message string                       `json:"message"`
		Data    DetailDataSelfieVerification `json:"data"`
	}
	DetailDataSelfieVerification struct {
		ID          string   `json:"id"`
		Status      string   `json:"status"`
		RejectField []string `json:"reject_field"`
		ReferenceID string   `json:"reference_id"`
	}

	VerihubsErrorResponse struct {
		Message     string               `json:"message"`
		ErrorCode   string               `json:"error_code"`
		ErrorFields []VerihubsErrorField `json:"error_fields"`
	}

	VerihubsErrorField struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	DetailVerifySelfie struct {
		Id          string   `json:"id"`
		Status      string   `json:"status"`
		RejectField []string `json:"reject_field"`
		ReferenceId string   `json:"reference_id"`
	}

	ResultDataKycPassport struct {
		Authority    string `json:"authority,omitempty"`
		DateOfBirth  string `json:"date_of_birth,omitempty"`
		DateOfExpiry string `json:"date_of_expiry,omitempty"`
		DateOfIssue  string `json:"date_of_issue,omitempty"`
		FullName     string `json:"full_name,omitempty"`
		Gender       string `json:"gender,omitempty"`
		Nationality  string `json:"nationality,omitempty"`
		PassportNo   string `json:"passport_no,omitempty"`
		PlaceOfBirth string `json:"place_of_birth,omitempty"`
	}
)
