package entity

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	UserID         uint   `gorm:"column:user_id" json:"user_id"`
	UserUUID       string `gorm:"column:user_uuid" json:"user_uuid"`
	DeviceID       string `gorm:"column:device_id" json:"device_id"`
	FCMToken       string `gorm:"column:fcm_token" json:"fcm_token"`
	AppVersionCode string `gorm:"column:app_version_code" json:"app_version_code"`
	AppVersionName string `gorm:"column:app_version_name" json:"app_version_name"`
	Manufacturer   string `gorm:"column:manufacturer" json:"manufacturer"`
	Brand          string `gorm:"column:brand" json:"brand"`
	DeviceModel    string `gorm:"column:device_model" json:"device_model"`
	Product        string `gorm:"column:product" json:"product"`
	VersionSdk     string `gorm:"column:version_sdk" json:"version_sdk"`
	VersionRelease string `gorm:"column:version_release" json:"version_release"`
}

func (i Device) TableName() string { return "devices" }
