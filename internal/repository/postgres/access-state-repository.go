package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"gorm.io/gorm"
)

type accessStateRepository struct {
	masterDb *gorm.DB
	logger   *helpers.CustomLogger
}
type AccessStateRepository interface {
	Tx(ctx context.Context) *gorm.DB
	InsertAccessStateRepository(ctx context.Context, tx *gorm.DB, resetPin *entity.AccessState) error
	SelectAccessByAccessTokenRepository(ctx context.Context, resetToken string) (*entity.AccessState, error)
	SelectActiveAccessByUserId(ctx context.Context, userId int64) (*entity.AccessState, error)
	UpdateAccessByStruct(ctx context.Context, tx *gorm.DB, curentPin *entity.AccessState, newPin *entity.AccessState) error
}

func NewResetPinRepository(posgres *gorm.DB,
	logger *helpers.CustomLogger) AccessStateRepository {
	return &accessStateRepository{
		masterDb: posgres,
		logger:   logger,
	}
}
func (repo *accessStateRepository) Tx(ctx context.Context) *gorm.DB {
	return repo.masterDb.Begin()
}
func (repo *accessStateRepository) InsertAccessStateRepository(ctx context.Context, tx *gorm.DB, resetPin *entity.AccessState) error {
	err := tx.Create(&resetPin).Error
	if err != nil {
		repo.logger.ErrorLogger(ctx, "InsertAccessStateRepository", err)
	}
	return err
}
func (repo *accessStateRepository) SelectAccessByAccessTokenRepository(ctx context.Context, resetToken string) (*entity.AccessState, error) {
	var resetPin entity.AccessState
	err := repo.masterDb.Table("access_states").Where("access_token = ? and used = false", resetToken).Order("created_at DESC").First(&resetPin).Error
	if err != nil {
		repo.logger.ErrorLogger(ctx, "SelectAccessByAccessTokenRepository", err)
		return nil, err
	}
	return &resetPin, nil
}
func (repo *accessStateRepository) SelectActiveAccessByUserId(ctx context.Context, userId int64) (*entity.AccessState, error) {
	var pinData entity.AccessState
	err := repo.masterDb.Table("access_states").Where("user_id = ? and used = true", userId).First(&pinData).Error
	if err != nil {
		repo.logger.ErrorLogger(ctx, "SelectActiveAccessByUserId", err)
		return nil, err
	}
	return &pinData, nil
}
func (repo *accessStateRepository) UpdateAccessByStruct(ctx context.Context, tx *gorm.DB, curentPin *entity.AccessState, newPin *entity.AccessState) error {
	err := tx.Model(curentPin).Updates(newPin).Error
	if err != nil {
		repo.logger.ErrorLogger(ctx, "UpdateAccessByStruct", err)
	}
	return err
}
