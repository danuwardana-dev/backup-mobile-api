package banklistsvc

import (
	repository "backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/entity"
	"context"
)

type BankService interface {
	GetBanks(ctx context.Context) ([]entity.Bank, error)
	SearchBanks(ctx context.Context, keyword string) ([]entity.Bank, error)
	GetBanksByType(ctx context.Context, bankType string) ([]entity.Bank, error) // NEW
	SearchBanksByType(ctx context.Context, bankType, keyword string) ([]entity.Bank, error)
}

type bankService struct {
	repo repository.BankListRepository
}

func (s *bankService) SearchBanksByType(ctx context.Context, bankType, keyword string) ([]entity.Bank, error) {
	return s.repo.SearchBanksByType(ctx, bankType, keyword)
}

func NewBankListService(repo repository.BankListRepository) BankService {
	return &bankService{repo: repo}
}

func (s *bankService) GetBanks(ctx context.Context) ([]entity.Bank, error) {
	return s.repo.GetAllBanks(ctx)
}

func (s *bankService) SearchBanks(ctx context.Context, keyword string) ([]entity.Bank, error) {
	return s.repo.SearchBanks(ctx, keyword)
}

func (s *bankService) GetBanksByType(ctx context.Context, bankType string) ([]entity.Bank, error) {
	return s.repo.GetBanksByType(ctx, bankType)
}
