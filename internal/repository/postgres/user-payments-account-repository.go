package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"

	"gorm.io/gorm"
)

type UserPaymentsAccountRepository interface {
	SearchBanksAccountUserPayment(ctx context.Context, keyword string) ([]entity.UserPaymentsAccount, error)
	GetAllBanksAccountUserPayment(ctx context.Context) ([]entity.UserPaymentsAccount, error)
}

type userPaymentsAccountRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewUserPaymentsAccountRepository(masterDb *gorm.DB, clogger *helpers.CustomLogger) UserPaymentsAccountRepository {
	return &userPaymentsAccountRepository{masterDb: masterDb, clogger: clogger}
}
func (r *userPaymentsAccountRepository) GetAllBanksAccountUserPayment(ctx context.Context) ([]entity.UserPaymentsAccount, error) {
	var userPaymentsAccount []entity.UserPaymentsAccount
	err := r.masterDb.WithContext(ctx).
		Order("bank_name ASC"). // pastikan sesuai dengan entity.Bank
		Find(&userPaymentsAccount).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "GetAllBanksAccountUserPayment", err)
		return nil, err
	}
	return userPaymentsAccount, nil
}

func (r *userPaymentsAccountRepository) SearchBanksAccountUserPayment(ctx context.Context, keyword string) ([]entity.UserPaymentsAccount, error) {
	var userPaymentsAccount []entity.UserPaymentsAccount
	err := r.masterDb.WithContext(ctx).
		Where("bank_name ILIKE ? ", "%"+keyword+"%").
		Order("bank_name ASC").
		Find(&userPaymentsAccount).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "SearchBankUserPaymentAccount", err)
		return nil, err
	}
	return userPaymentsAccount, nil
}
