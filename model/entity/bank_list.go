package entity

type Bank struct {
	ID       uint    `gorm:"column:bank_id;primaryKey" json:"id"`
	NamaBank string  `gorm:"column:nama_bank" json:"nama_bank"`
	UrlImage string  `gorm:"column:url_image" json:"url_image"`
	VAName   string  `json:"va_name" gorm:"column:va_name"`
	Price    float64 `json:"price" gorm:"column:price"`
	AdminFee float64 `json:"admin_fee" gorm:"column:admin_fee"`
}

// TableName override nama tabel
func (Bank) TableName() string {
	return "tb_bank_list"
}
