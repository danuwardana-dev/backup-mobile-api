package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"

	"gorm.io/gorm"
)

type kycPassportRepository struct {
	masterDb *gorm.DB
	clog     *helpers.CustomLogger
}

// SaveKYCPassport implements KycPassportRepository.
func (k *kycPassportRepository) SaveKYCPassport(orm *gorm.DB, ctx context.Context, passport *entity.IdentityPassport) error {
	err := orm.WithContext(ctx).Create(&passport).Error
	if err != nil {
		k.clog.ErrorLogger(ctx, "SaveKYCPassport.repo.masterDb.WithContext(ctx).Create", err)
		return err
	}
	return err
}

// Tx implements KycPassportRepository.
func (k *kycPassportRepository) Tx(ctx context.Context) *gorm.DB {
	return k.masterDb.Begin()
}

type KycPassportRepository interface {
	Tx(ctx context.Context) *gorm.DB
	SaveKYCPassport(orm *gorm.DB, ctx context.Context, ktp *entity.IdentityPassport) error
}

func NewKycPassportRepository(masterDb *gorm.DB, clog *helpers.CustomLogger) KycPassportRepository {
	return &kycPassportRepository{
		masterDb: masterDb,
		clog:     clog,
	}
}
