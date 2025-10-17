package entity

type UserPaymentsAccount struct {
	ID       uint   `gorm:"column:id;primaryKey" json:"id"`
	UserID   uint   `gorm:"column:user_id" json:"user_id"`
	NoRek    string `gorm:"column:no_rekening" json:"no_rekening"`
	NoVa     string `json:"no_va" gorm:"column:no_va"`
	BankName string `json:"bank_name" gorm:"column:bank_name"`
}

// TableName override nama tabel
func (UserPaymentsAccount) TableName() string {
	return "user_payment_accounts"
}
