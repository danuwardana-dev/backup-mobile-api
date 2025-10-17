package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"gorm.io/gorm"
)

type deviceRepository struct {
	masterDb *gorm.DB
	clog     *helpers.CustomLogger
}
type DeviceRepository interface {
	Tx(ctx context.Context) *gorm.DB
	SelectDeviceByStruct(ctx context.Context, device *entity.Device) ([]entity.Device, error)
	InsertDevice(ctx context.Context, tx *gorm.DB, device *entity.Device) error
}

func NewDeviceRepository(db *gorm.DB, clog *helpers.CustomLogger) DeviceRepository {
	return &deviceRepository{
		masterDb: db,
		clog:     clog,
	}
}

func (repo *deviceRepository) Tx(ctx context.Context) *gorm.DB {
	return repo.masterDb.Begin()
}
func (repo *deviceRepository) SelectDeviceByStruct(ctx context.Context, device *entity.Device) ([]entity.Device, error) {
	var devicesData []entity.Device
	err := repo.masterDb.WithContext(ctx).Where(device).Order("created_at DESC").Find(&devicesData).Error
	if err != nil {
		repo.clog.ErrorLogger(ctx, "SelectDeviceByStruct.repo.masterDb.WithContext(ctx).Where(device).Find", err)
		return nil, err
	}
	return devicesData, nil
}
func (repo *deviceRepository) InsertDevice(ctx context.Context, tx *gorm.DB, device *entity.Device) error {
	err := tx.WithContext(ctx).Create(device).Error
	if err != nil {
		repo.clog.ErrorLogger(ctx, "InsertDevice.repo.masterDb.WithContext(ctx).Create", err)
	}
	return err
}
