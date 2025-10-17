package recipientsvc

import (
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/entity"
	"context"
	"log"
)

type RecipientService interface {
	GetRecipients(ctx context.Context) ([]entity.RecipientWithBank, error)
	SearchRecipients(ctx context.Context, keyword string) ([]entity.RecipientWithBank, error)
	AddRecipient(ctx context.Context, recipient *entity.Recipient) error
	GetUserByUUID(ctx context.Context, uuid string) (*entity.User, error)
}

type recipientService struct {
	repo postgres.RecipientRepository
}

func NewRecipientService(repo postgres.RecipientRepository) RecipientService {
	if repo == nil {
		log.Println("[ERROR] repo nil saat init service")
	}
	return &recipientService{repo: repo}
}

// ✅ ambil semua recipient beserta url_image bank
func (s *recipientService) GetRecipients(ctx context.Context) ([]entity.RecipientWithBank, error) {
	if s == nil {
		log.Println("[ERROR] service struct nil")
	}
	// if s.repo == nil {
	// 	log.Println("[ERROR] service.repo nil sebelum dipakai")
	// }
	return s.repo.GetAllRecipients(ctx)
}

// ✅ search recipient berdasarkan nama / no rekening
func (s *recipientService) SearchRecipients(ctx context.Context, keyword string) ([]entity.RecipientWithBank, error) {
	return s.repo.SearchRecipients(ctx, keyword)
}

// ✅ insert recipient baru
func (s *recipientService) AddRecipient(ctx context.Context, recipient *entity.Recipient) error {
	return s.repo.InsertRecipient(ctx, recipient)
}
func (s *recipientService) GetUserByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	return s.repo.FindUserByUUID(ctx, uuid) // panggil repo
}
