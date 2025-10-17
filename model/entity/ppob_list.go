package entity

type PPOB struct {
	ID           uint   `gorm:"column:id_list_ppob_internet;primaryKey" json:"id"`
	NameProvider string `gorm:"column:name_provider" json:"name_provider"`
	UrlImage     string `gorm:"column:image_url" json:"url_image"`
}

func (PPOB) TableName() string {
	return "list_ppob_internet"
}
