package entity

import (
	"backend-mobile-api/model/enum"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        enum.RolesEnum
	Description string
	Users       []User `gorm:"many2many:user_roles;" json:"users"`
}

func (r Role) TableName() string {
	return "roles"
}
