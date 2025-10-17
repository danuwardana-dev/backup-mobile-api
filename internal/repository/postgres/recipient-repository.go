package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"log"

	"gorm.io/gorm"
)

type RecipientRepository interface {
	GetAllRecipients(ctx context.Context) ([]entity.RecipientWithBank, error)
	SearchRecipients(ctx context.Context, keyword string) ([]entity.RecipientWithBank, error)
	InsertRecipient(ctx context.Context, recipient *entity.Recipient) error
	FindUserByUUID(ctx context.Context, uuid string) (*entity.User, error) // ✅ tambahan

}

type recipientRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewRecipientRepository(masterDb *gorm.DB, clogger *helpers.CustomLogger) RecipientRepository {
	if masterDb == nil {
		log.Println("[ERROR] masterDb nil saat init repository")
	}
	if clogger == nil {
		log.Println("[ERROR] clogger nil saat init repository")
	}
	return &recipientRepository{masterDb: masterDb, clogger: clogger}
}

// ✅ Ambil semua recipient beserta url_image bank
func (r *recipientRepository) GetAllRecipients(ctx context.Context) ([]entity.RecipientWithBank, error) {
	var results []entity.RecipientWithBank

	err := r.masterDb.WithContext(ctx).
		Table("tb_recipient as r").
		Select(`r.recipient_id, 
		        r.nama_penerima, 
		        r.no_rekening, 
		        r.user_id,
		        r.bank_id, 
		        b.url_image as bank_image_url,
		        b.nama_bank as nama_bank
				`).
		Joins("LEFT JOIN tb_bank_list b ON r.bank_id = b.bank_id").
		Scan(&results).Error

	if err != nil {
		r.clogger.ErrorLogger(ctx, "GetAllRecipients", err)
		return nil, err
	}
	return results, nil
}

// ✅ Search berdasarkan nama penerima / no rekening
func (r *recipientRepository) SearchRecipients(ctx context.Context, keyword string) ([]entity.RecipientWithBank, error) {
	var recipients []entity.RecipientWithBank

	err := r.masterDb.WithContext(ctx).
		Table("tb_recipient as r").
		Select(`recipient_id, 
		        r.nama_penerima, 
		        r.no_rekening, 
		        r.user_id,
		        r.bank_id, 
		        b.url_image as bank_image_url,
		        b.nama_bank as nama_bank
				`).
		Joins("LEFT JOIN tb_bank_list b ON r.bank_id = b.bank_id").
		Where("LOWER(r.nama_penerima) LIKE ? OR r.no_rekening LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Scan(&recipients).Error

	if err != nil {
		r.clogger.ErrorLogger(ctx, "SearchRecipients", err)
		return nil, err
	}
	return recipients, nil
}

// ✅ Insert hanya ke tb_recipient
func (r *recipientRepository) InsertRecipient(ctx context.Context, recipient *entity.Recipient) error {
	err := r.masterDb.WithContext(ctx).Create(recipient).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "InsertRecipient", err)
	}
	return err
}

// ✅ Cari user by UUID (biar dapet user_id dari DB)
func (r *recipientRepository) FindUserByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	var user entity.User
	err := r.masterDb.WithContext(ctx).Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "FindUserByUUID", err)
		return nil, err
	}
	return &user, nil
}
