package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"

	"gorm.io/gorm"
)

type BankListRepository interface {
	GetAllBanks(ctx context.Context) ([]entity.Bank, error)
	SearchBanks(ctx context.Context, keyword string) ([]entity.Bank, error)
	GetBanksByType(ctx context.Context, bankType string) ([]entity.Bank, error)
	SearchBanksByType(ctx context.Context, bankType, keyword string) ([]entity.Bank, error)
}

type bankListRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewBankListRepository(masterDb *gorm.DB, clogger *helpers.CustomLogger) BankListRepository {
	return &bankListRepository{masterDb: masterDb, clogger: clogger}
}

// ✅ ambil semua bank/ewallet
func (r *bankListRepository) GetAllBanks(ctx context.Context) ([]entity.Bank, error) {
	var banks []entity.Bank
	err := r.masterDb.WithContext(ctx).
		Order("nama_bank ASC"). // pastikan sesuai dengan entity.Bank
		Find(&banks).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "GetAllBanks", err)
		return nil, err
	}
	return banks, nil
}

// ✅ search by nama bank / va_name tapi dengan filter type (BANK / EWALLET)
func (r *bankListRepository) SearchBanksByType(ctx context.Context, bankType, keyword string) ([]entity.Bank, error) {
	var banks []entity.Bank
	err := r.masterDb.WithContext(ctx).
		Where("type = ? AND (nama_bank ILIKE ? OR va_name ILIKE ?)", bankType, "%"+keyword+"%", "%"+keyword+"%").
		Order("nama_bank ASC").
		Find(&banks).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "SearchBanksByType", err)
		return nil, err
	}
	return banks, nil
}

// ✅ search by nama bank atau va_name
func (r *bankListRepository) SearchBanks(ctx context.Context, keyword string) ([]entity.Bank, error) {
	var banks []entity.Bank
	err := r.masterDb.WithContext(ctx).
		Where("nama_bank ILIKE ? OR va_name ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("nama_bank ASC").
		Find(&banks).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "SearchBanks", err)
		return nil, err
	}
	return banks, nil
}

// ✅ filter by type (BANK / EWALLET)
func (r *bankListRepository) GetBanksByType(ctx context.Context, bankType string) ([]entity.Bank, error) {
	var banks []entity.Bank
	err := r.masterDb.WithContext(ctx).
		Where("type = ?", bankType).
		Order("nama_bank ASC").
		Find(&banks).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "GetBanksByType", err)
		return nil, err
	}
	return banks, nil
}
