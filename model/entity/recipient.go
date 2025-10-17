package entity

type Recipient struct {
	RecipientID  uint   `gorm:"column:recipient_id;primaryKey;autoIncrement" json:"recipient_id"`
	Bank         int    `gorm:"column:bank_id" json:"bank_id"`
	UserUUID     string `gorm:"-" json:"user_uuid"`
	User         int64  `gorm:"column:user_id" json:"user_id"`
	NamaPenerima string `gorm:"column:nama_penerima" json:"nama_penerima"`
	NoRekening   string `gorm:"column:no_rekening" json:"no_rekening"`
}

func (Recipient) TableName() string {
	return "tb_recipient"
}

type RecipientWithBank struct {
	RecipientID uint  `json:"recipient_id"`
	BankID      int   `json:"bank_id"`
	UserID      int64 `json:"user_id"`

	NamaPenerima string `json:"nama_penerima"`
	NoRekening   string `json:"no_rekening"`
	BankImageURL string `json:"bank_image_url"`
	NamaBank     string `json:"nama_bank"`
}

func EncapsulateRequestRecipientToEntity(u *User) *Recipient {
	return &Recipient{

		User: u.ID, // <- ini ngisi kolom user_id di tabel tb_recipient

	}
}
