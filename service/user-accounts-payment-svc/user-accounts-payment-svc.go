package useraccountspaymentsvc

import (
	repository "backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/entity"
	"context"
)

type UserAccountsPaymentService interface {
	SearchUserAccountsPaymentService(ctx context.Context, keyword string) ([]entity.UserPaymentsAccount, error)
	GetUserAccountsPaymentService(ctx context.Context) ([]entity.UserPaymentsAccount, error)
}

type userAccountsPaymentService struct {
	repo repository.UserPaymentsAccountRepository
}

func NewUserAccountsPaymentService(repo repository.UserPaymentsAccountRepository) UserAccountsPaymentService {
	return &userAccountsPaymentService{repo: repo}
}
func (s *userAccountsPaymentService) SearchUserAccountsPaymentService(ctx context.Context, keyword string) ([]entity.UserPaymentsAccount, error) {
	return s.repo.SearchBanksAccountUserPayment(ctx, keyword)
}
func (s *userAccountsPaymentService) GetUserAccountsPaymentService(ctx context.Context) ([]entity.UserPaymentsAccount, error) {
	return s.repo.GetAllBanksAccountUserPayment(ctx)
}
