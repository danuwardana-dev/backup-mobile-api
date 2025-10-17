package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"

	"gorm.io/gorm"
)

type kycKtpRepository struct {
	masterDb *gorm.DB
	clog     *helpers.CustomLogger
}

// SaveKYCKtp implements KycKtpRepository.
func (k *kycKtpRepository) SaveKYCKtp(ctx context.Context, orm *gorm.DB, ktp *entity.IdentityKtp) error {
	err := orm.WithContext(ctx).Create(&ktp).Error
	if err != nil {
		k.clog.ErrorLogger(ctx, "SaveKYCKtp.repo.masterDb.WithContext(ctx).Create", err)
		return err
	}
	return err
}

// Tx implements KycKtpRepository.
func (k *kycKtpRepository) Tx(ctx context.Context) *gorm.DB {
	return k.masterDb.Begin()
}

type KycKtpRepository interface {
	Tx(ctx context.Context) *gorm.DB
	SaveKYCKtp(ctx context.Context, orm *gorm.DB, ktp *entity.IdentityKtp) error
}

func NewKycKtpRepository(masterDb *gorm.DB, clog *helpers.CustomLogger) KycKtpRepository {
	return &kycKtpRepository{
		masterDb: masterDb,
		clog:     clog,
	}
}
