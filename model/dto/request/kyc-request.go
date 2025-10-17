package request

type (
	KycRequest struct {
		Image string `json:"image,omitempty"`
	}
	VerifyKycSelfie struct {
		Nik         string `json:"nik,omitempty"`
		Name        string `json:"name,omitempty"`
		BirthDate   string `json:"birth_date,omitempty"`
		Email       string `json:"email,omitempty"`
		Phone       string `json:"phone,omitempty"`
		SelfiePhoto string `json:"selfie_photo,omitempty"`
		KtpPhoto    string `json:"ktp_photo,omitempty"`
	}

	KTPrequest struct {
		Nik            string               `json:"nik,omitempty"`
		FullName       string               `json:"full_name,omitempty"`
		PlaceOfBirth   string               `json:"place_of_birth,omitempty"`
		Gender         string               `json:"gender,omitempty"`
		DateOfBirth    string               `json:"date_of_birth,omitempty"`
		Occupation     string               `json:"occupation,omitempty"`
		Nationality    string               `json:"nationality,omitempty"`
		MartialStatus  string               `json:"martial_status,omitempty"`
		Religion       string               `json:"religion,omitempty"`
		Country        string               `json:"country,omitempty"`
		State          string               `json:"state,omitempty"`
		City           string               `json:"city,omitempty"`
		District       string               `json:"district,omitempty"`
		FullAddress    string               `json:"address,omitempty"`
		CurrentAddress DetailCurrentAddress `json:"current_address,omitempty"`
		Image          string               `json:"image,omitempty"`
	}

	PassportRequest struct {
		PassportType   string               `json:"passport_type,omitempty"`
		PassportNo     string               `json:"passport_no,omitempty"`
		DateOfIssue    string               `json:"date_of_issue,omitempty"`
		DateOfExpired  string               `json:"date_of_expired,omitempty"`
		ResiNumber     string               `json:"resi_number,omitempty"`
		PlaceOfIssue   string               `json:"place_of_issue,omitempty"`
		Gender         string               `json:"gender,omitempty"`
		FullName       string               `json:"full_name,omitempty"`
		Nationality    string               `json:"nationality,omitempty"`
		PlaceOfBirth   string               `json:"place_of_birth,omitempty"`
		DateOfBirth    string               `json:"date_of_birth,omitempty"`
		CurrentAddress DetailCurrentAddress `json:"current_address,omitempty"`
		Image          string               `json:"image,omitempty"`
	}

	DetailCurrentAddress struct {
		Country     string `json:"country,omitempty"`
		City        string `json:"city,omitempty"`
		District    string `json:"district,omitempty"`
		FullAddress string `json:"full_address,omitempty"`
	}
)
