package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"gorm.io/gorm"
	"time"
)

type tokenBlacklistRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}
type TokenBlacklistTokenRepository interface {
	Tx(ctx context.Context) *gorm.DB
	InsertBlaclistToken(ctx context.Context, tx *gorm.DB, blaclistToken *entity.TokenBlacklist) error
	IsBlaclistTokenActive(ctx context.Context, token string) (bool, error)
}

func NewTokenBlacklistTokenRepository(masterDb *gorm.DB, clogger *helpers.CustomLogger) TokenBlacklistTokenRepository {
	return &tokenBlacklistRepository{masterDb: masterDb, clogger: clogger}
}
func (repo *tokenBlacklistRepository) Tx(ctx context.Context) *gorm.DB {
	return repo.masterDb.Begin()
}
func (repo *tokenBlacklistRepository) InsertBlaclistToken(ctx context.Context, tx *gorm.DB, blaclistToken *entity.TokenBlacklist) error {
	err := tx.Create(blaclistToken).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "InsertBlaclistToken.tx.Create", err)
	}
	return err
}
func (repo *tokenBlacklistRepository) IsBlaclistTokenActive(ctx context.Context, token string) (bool, error) {
	var count int64
	err := repo.masterDb.Table("token_blacklists").Where("token = ? AND expired_at <= ?", token, time.Now()).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
