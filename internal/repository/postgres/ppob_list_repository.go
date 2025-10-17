package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"

	"gorm.io/gorm"
)

type PpobListRepository interface {
	GetAllPpob(ctx context.Context) ([]entity.PPOB, error)
	SearchPpob(ctx context.Context, keyword string) ([]entity.PPOB, error)
}

type ppobListRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewPpobListRepository(masterDb *gorm.DB, clogger *helpers.CustomLogger) PpobListRepository {
	return &ppobListRepository{masterDb: masterDb, clogger: clogger}
}

func (r *ppobListRepository) GetAllPpob(ctx context.Context) ([]entity.PPOB, error) {
	var banks []entity.PPOB
	err := r.masterDb.WithContext(ctx).
		Order("name_provider ASC"). // pastikan sesuai dengan entity.Bank
		Find(&banks).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "GetAllProvider", err)
		return nil, err
	}
	return banks, nil
}

func (r *ppobListRepository) SearchPpob(ctx context.Context, keyword string) ([]entity.PPOB, error) {
	var banks []entity.PPOB
	err := r.masterDb.WithContext(ctx).
		Where("name_provider ILIKE ?", "%"+keyword+"%").
		Order("name_provider ASC").
		Find(&banks).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "SearchPpob", err)
		return nil, err
	}
	return banks, nil
}
