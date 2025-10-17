package ppoblistsvc

import (
	repository "backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/entity"
	"context"
)

type PpobListService interface {
	GetPpobList(ctx context.Context) ([]entity.PPOB, error)
	SearchPpobList(ctx context.Context, keyword string) ([]entity.PPOB, error)
}

type ppobListService struct {
	repo repository.PpobListRepository
}

func NewPpobListService(repo repository.PpobListRepository) PpobListService {
	return &ppobListService{repo: repo}
}

func (s *ppobListService) GetPpobList(ctx context.Context) ([]entity.PPOB, error) {
	return s.repo.GetAllPpob(ctx)
}

func (s *ppobListService) SearchPpobList(ctx context.Context, keyword string) ([]entity.PPOB, error) {
	return s.repo.SearchPpob(ctx, keyword)
}
