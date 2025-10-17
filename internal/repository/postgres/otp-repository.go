package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"gorm.io/gorm"
	"time"
)

type otpRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewOtpRepository(posgres *gorm.DB, clogger *helpers.CustomLogger) OtpRepository {
	return &otpRepository{masterDb: posgres, clogger: clogger}
}

type OtpRepository interface {
	InsertOtpDataRepository(ctx context.Context, tx *gorm.DB, otpData *entity.OTP) error
	Tx(ctx context.Context) *gorm.DB
	SelectOtpBySessionId(ctx context.Context, sessionId string) (*entity.OTP, error)
	SelectOtpByVerifyKey(ctx context.Context, verifyKey string) (*entity.OTP, error)
	SelectOtpByVerifyKeyBeforeExpire(ctx context.Context, verifyKey string) (*entity.OTP, error)
	UpdateOtpDataRepository(ctx context.Context, tx *gorm.DB, otpData *entity.OTP, updater *entity.OTP) error
}

func (repo *otpRepository) Tx(ctx context.Context) *gorm.DB {
	return repo.masterDb.Begin()
}
func (r *otpRepository) InsertOtpDataRepository(ctx context.Context, tx *gorm.DB, otpData *entity.OTP) error {
	err := r.masterDb.Create(otpData).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "InsertOtpData.masterDb.Create", err)
	}
	return err
}
func (repo *otpRepository) SelectOtpBySessionId(ctx context.Context, sessionId string) (*entity.OTP, error) {
	var otp entity.OTP
	err := repo.masterDb.Table("otps").Where("session_id = ?", sessionId).Order("created_at DESC").First(&otp).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectOtpBySessionId.gorm.DB", err)
		return nil, err
	}
	return &otp, nil
}
func (repo *otpRepository) SelectOtpByVerifyKey(ctx context.Context, verifyKey string) (*entity.OTP, error) {
	var otp entity.OTP
	err := repo.masterDb.Table("otps").Where("verify_key = ?", verifyKey).Order("created_at DESC").First(&otp).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectOtpByVerifyKey.gorm.DB", err)
		return nil, err
	}
	return &otp, nil
}
func (repo *otpRepository) UpdateOtpDataRepository(ctx context.Context, tx *gorm.DB, otpData *entity.OTP, updater *entity.OTP) error {
	if err := tx.Model(*otpData).Updates(*updater).Error; err != nil {
		repo.clogger.ErrorLogger(ctx, "UpdateOtpData.gorm.DB", err)
		return err
	}
	return nil
}
func (repo *otpRepository) SelectOtpByVerifyKeyBeforeExpire(ctx context.Context, verifyKey string) (*entity.OTP, error) {
	var otp entity.OTP
	err := repo.masterDb.Table("otps").Where("verify_key = ? AND expired_at > ?", verifyKey, time.Now()).Order("created_at DESC").First(&otp).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectOtpByVerifyKeyBeforeExpire.gorm.DB", err)
		return nil, err
	}
	return &otp, nil
}
