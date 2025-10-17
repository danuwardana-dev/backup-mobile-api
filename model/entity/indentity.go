package entity

import (
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/types"
	"time"
)

type (
	IdentityKtp struct {
		ID            int64     `gorm:"column:id;primaryKey" json:"id"`
		UserId        int64     `gorm:"column:user_id" json:"user_id"`
		Nik           string    `gorm:"column:nik;type:varchar" json:"nik"`
		FullName      string    `gorm:"column:full_name;type:varchar" json:"full_name"`
		PlaceOfBirth  string    `gorm:"column:place_of_birth;type:varchar" json:"place_of_birth"`
		Gender        string    `gorm:"column:gender;type:varchar;size:20" json:"gender"`
		DateOfBirth   time.Time `gorm:"column:date_of_birth;type:date" json:"date_of_birth"`
		Occupation    string    `gorm:"column:occupation;type:varchar" json:"occupation"`
		Nationality   string    `gorm:"column:nationality;type:varchar" json:"nationality"`
		MartialStatus string    `gorm:"column:martial_status;type:varchar" json:"martial_status"`
		Religion      string    `gorm:"column:religion;type:varchar" json:"religion"`
		Country       string    `gorm:"column:country;type:varchar" json:"country"`
		Province      string    `gorm:"column:province;type:varchar" json:"province"`
		City          string    `gorm:"column:city;type:varchar" json:"city"`
		District      string    `gorm:"column:district;type:varchar" json:"district"`
		FullAddress   string    `gorm:"column:full_address;type:text" json:"full_address"`
		IdentityImage string    `gorm:"column:identity_image;type:varchar" json:"identity_image"`
		SubmitAt      time.Time `gorm:"column:submit_at;type:timestamp" json:"submit_at"`
		VerifyAt      time.Time `gorm:"column:verify_at;type:datetime" json:"verify_at"`
		UpdatedAt     time.Time `gorm:"type:timestamptz" json:"updated_at" json:"updated_at"`

		User User `gorm:"foreignKey:UserId" json:"user"`
	}

	IdentityPassport struct {
		ID             int64      `gorm:"column:id;primaryKey" json:"id"`
		UserId         int64      `gorm:"column:user_id" json:"user_id"`
		PassportNumber string     `gorm:"column:passport_number;type:varchar;size:255" json:"passport_number"`
		PassportType   string     `gorm:"column:passport_type;type:varchar;size:50"`
		Gender         string     `gorm:"column:gender;type:varchar;size:20"`
		FullName       string     `gorm:"column:full_name;type:varchar" json:"full_name"`
		Nationality    string     `gorm:"column:nationality;type:varchar;size:100"`
		PlaceOfBirth   string     `gorm:"column:place_of_birth;type:varchar;size:50" json:"place_of_birth"`
		DateOfBirth    types.Date `gorm:"column:date_of_birth;type:date" json:"date_of_birth"`
		DateOfIssue    types.Date `gorm:"column:date_of_issue;type:date" json:"date_of_issue"`
		DateOfExpired  types.Date `gorm:"column:date_of_expired;type:date" json:"date_of_expired"`
		ResiNumber     string     `gorm:"column:resi_number;type:varchar;size:255" json:"resi_number"`
		PlaceOfIssue   string     `gorm:"column:place_of_issue;type:varchar;size:100" json:"place_of_issue"`
		IdentityImage  string     `gorm:"column:identity_image;type:varchar" json:"identity_image"`
		SubmitAt       time.Time  `gorm:"column:submit_at;type:timestamp" json:"submit_at"`
		VerifyAt       time.Time  `gorm:"column:verify_at;type:datetime" json:"verify_at"`
		UpdatedAt      time.Time  `gorm:"type:timestamptz" json:"updated_at" json:"updated_at"`

		User User `gorm:"foreignKey:UserId" json:"user"`
	}
)

func EncapsulateRequestKtpToEntity(request request.KTPrequest, u *User, imagePath string) *IdentityKtp {
	birthDate, err := time.Parse("02-01-2006", request.DateOfBirth)
	if err != nil {
		return nil // atau log error, validasi gagal, dll
	}
	return &IdentityKtp{
		UserId:        u.ID,
		Nik:           request.Nik,
		FullName:      request.FullName,
		PlaceOfBirth:  request.PlaceOfBirth,
		Gender:        request.Gender,
		DateOfBirth:   birthDate,
		Occupation:    request.Occupation,
		Nationality:   request.Nationality,
		MartialStatus: request.MartialStatus,
		Religion:      request.Religion,
		Country:       "INDONESIA",
		Province:      request.State,
		City:          request.City,
		District:      request.District,
		FullAddress:   request.FullAddress,
		IdentityImage: imagePath,
	}
}

func EncapsulateRequestPassportToEntity(request request.PassportRequest, u *User, pathImage string) *IdentityPassport {
	dateOfBirth, err := types.NewDate(request.DateOfBirth)
	if err != nil {
		return nil
	}

	dateOfIssue, err := types.NewDate(request.DateOfIssue)
	if err != nil {
		return nil
	}

	dateOfExpiry, err := types.NewDate(request.DateOfExpired)
	if err != nil {
		return nil
	}

	return &IdentityPassport{
		UserId:         u.ID,
		PassportNumber: request.PassportNo,
		PassportType:   request.PassportType,
		Gender:         request.Gender,
		FullName:       request.FullName,
		Nationality:    request.Nationality,
		PlaceOfBirth:   request.PlaceOfBirth,
		DateOfBirth:    dateOfBirth,
		DateOfIssue:    dateOfIssue,
		DateOfExpired:  dateOfExpiry,
		ResiNumber:     request.ResiNumber,
		IdentityImage:  pathImage,
		PlaceOfIssue:   request.PlaceOfIssue,
	}
}

func (i IdentityKtp) TableName() string      { return "identity_ktps" }
func (i IdentityPassport) TableName() string { return "identity_passports" }
