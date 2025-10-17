package entity

type CheckAccountRequest struct {
	PartnerReferenceNo   string `json:"partnerReferenceNo"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode  string `json:"beneficiaryBankCode"`
	Type                 string `json:"type"`
}

type CheckAccountResponseData struct {
	BeneficiaryAccountNo   string `json:"beneficiaryAccountNo"`
	BeneficiaryAccountName string `json:"beneficiaryAccountName"`

	BeneficiaryBankCode string `json:"beneficiaryBankCode"`
}
