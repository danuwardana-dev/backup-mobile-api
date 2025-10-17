package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	// transaksi utama
	CreateTransaction(ctx context.Context, tx *entity.Transaction, userUUID string) error
	FindTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error)

	// detail transaksi
	CreateTransactionBankTransfer(ctx context.Context, detail *entity.TransactionBankTransfer) error
	CreateTransactionEwallet(ctx context.Context, detail *entity.TransactionEwallet) error
	CreateTransactionPhoneCredit(ctx context.Context, detail *entity.TransactionPhoneCredit) error
	CreateTransactionInternetTV(ctx context.Context, detail *entity.TransactionInternetTV) error
	CreateTransactionInternational(ctx context.Context, detail *entity.TransactionInternational) error
	//get all data
	GetAllTransactions(ctx context.Context, userID int64) ([]entity.Transaction, error)
	FindAllTransactionsByUserID(ctx context.Context, userID int64) ([]entity.Transaction, error)
	FindUserByUUID(ctx context.Context, uuid string) (*entity.User, error)
	GenerateTransactionID(ctx context.Context) (string, error)
	GenerateUniqueCode(ctx context.Context) (float64, error)
	ExpireOldTransactions(ctx context.Context) error
	GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	GetUserFcmToken(ctx context.Context, userID uint) (string, error)
	FindDeviceByUserUUID(ctx context.Context, uuid string) (*entity.Device, error)
	FindAllTransactionsByUserIDPaginated(
		ctx context.Context,
		userID int64,
		limit, offset int,
		search, status, txType, transactionID, startDate, endDate string,
	) ([]entity.Transaction, int64, error)
}

type transactionRepository struct {
	masterDb          *gorm.DB
	clogger           *helpers.CustomLogger
	lastTransactionID string
}

// TransactionRepository.go
func (r *transactionRepository) UpdateStatus(ctx context.Context, transactionID string, status string) error {
	return r.masterDb.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("transaction_id = ?", transactionID).
		Update("status", status).Error
}

func NewTransactionRepository(masterDb *gorm.DB, clogger *helpers.CustomLogger) TransactionRepository {
	if masterDb == nil {
		log.Println("[ERROR] masterDb nil saat init repository")
	}
	if clogger == nil {
		log.Println("[ERROR] clogger nil saat init repository")
	}
	return &transactionRepository{masterDb: masterDb, clogger: clogger}
}

// =======================
// PRIVATE FUNCTION
// =======================

// generate transaction_id → TXN20250917xxxx
func generateTransactionID() string {
	rand.Seed(time.Now().UnixNano())
	today := time.Now().Format("20060102")
	random := rand.Intn(10000)
	return fmt.Sprintf("TXN%s%d", today, random)
}
func (r *transactionRepository) FindDeviceByUserUUID(ctx context.Context, userUUID string) (*entity.Device, error) {
	var device entity.Device
	if err := r.masterDb.WithContext(ctx).Where("user_uuid = ?", userUUID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}
func (r *transactionRepository) GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error) {
	var tx entity.Transaction
	if err := r.masterDb.WithContext(ctx).Where("id = ?", id).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}
func (r *transactionRepository) GetUserFcmToken(ctx context.Context, userID uint) (string, error) {
	var device entity.Device
	if err := r.masterDb.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&device).Error; err != nil {
		return "", err
	}
	return device.FCMToken, nil
}

// =======================
// IMPLEMENTATION
// =======================

// create transaksi utama
func (r *transactionRepository) GenerateTransactionID(ctx context.Context) (string, error) {
	// ambil count transaksi atau pakai sequence
	var count int64
	if err := r.masterDb.WithContext(ctx).Model(&entity.Transaction{}).Count(&count).Error; err != nil {
		return "", err
	}

	newID := generateTransactionID()
	r.lastTransactionID = newID
	return newID, nil
}
func getCutoffTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	expiredStr := os.Getenv("EXPIRED_TIME")
	if expiredStr == "" {
		expiredStr = "6" // default 6 jam
	}

	expiredInt, err := strconv.Atoi(expiredStr)
	if err != nil {
		log.Printf("invalid EXPIRED_TIME value (%s), fallback to 6 hours", expiredStr)
		expiredInt = 6
	}

	cutoff := now.Add(-time.Duration(expiredInt) * time.Hour)
	return cutoff
}
func generateExpiredAt() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	expiredStr := os.Getenv("EXPIRED_TIME")
	if expiredStr == "" {
		expiredStr = "6" // default jika tidak di-set di env
	}

	expiredInt, err := strconv.Atoi(expiredStr)
	if err != nil {
		log.Printf("invalid EXPIRED_TIME value (%s), fallback to 6 hours", expiredStr)
		expiredInt = 6
	}

	return now.Add(time.Duration(expiredInt) * time.Hour)
}
func (r *transactionRepository) CreateTransaction(ctx context.Context, tx *entity.Transaction, userUUID string) error {
	// cari user_id dari uuid

	var user entity.User
	if err := r.masterDb.WithContext(ctx).Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return err
	}

	tx.TransactionID = r.lastTransactionID
	tx.UserID = user.ID
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	tx.ExpiredAt = generateExpiredAt()

	if tx.PaymentMethod == "bank_transfer" {
		code, err := r.GenerateUniqueCode(ctx)
		if err != nil {
			return err
		}
		tx.UniqueCode = code

	}
	if err := r.masterDb.WithContext(ctx).Create(tx).Error; err != nil {
		r.clogger.ErrorLogger(ctx, "CreateTransaction", err)
		return err
	}

	return nil
}

// generate transaction id

// generate unique code khusus bank transfer

// find transaksi by id
func (r *transactionRepository) FindTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.masterDb.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		Preload("BankTransfer").
		Preload("Ewallet").
		Preload("PhoneCredit").
		Preload("InternetTV").
		Preload("International").
		First(&transaction).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "FindTransactionByID", err)
		return nil, err
	}

	// Convert waktu ke WIB
	loc, _ := time.LoadLocation("Asia/Jakarta")
	transaction.CreatedAt = transaction.CreatedAt.In(loc)
	transaction.UpdatedAt = transaction.UpdatedAt.In(loc)
	transaction.ExpiredAt = transaction.ExpiredAt.In(loc)

	return &transaction, nil
}

// =======================
// DETAIL TRANSAKSI
// =======================

// bank transfer
func (r *transactionRepository) CreateTransactionBankTransfer(ctx context.Context, detail *entity.TransactionBankTransfer) error {
	// if detail.UniqueCode == 0 {
	// 	code, err := r.GenerateUniqueCode(ctx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	detail.UniqueCode = code
	// }
	// err := r.masterDb.WithContext(ctx).Create(detail).Error
	// if err != nil {
	// 	return err
	// }

	// return r.masterDb.WithContext(ctx).
	// 	Model(&entity.TransactionBankTransfer{}).
	// 	Where("transaction_id = ?", detail.TransactionID).
	// 	Update("is_reused", true).Error
	return r.masterDb.WithContext(ctx).Create(detail).Error
}

// ewallet
func (r *transactionRepository) CreateTransactionEwallet(ctx context.Context, detail *entity.TransactionEwallet) error {
	if detail.EwalletName == "" {
		return fmt.Errorf("ewallet name wajib diisi")
	}
	return r.masterDb.WithContext(ctx).Create(detail).Error
}

// phone credit
func (r *transactionRepository) CreateTransactionPhoneCredit(ctx context.Context, detail *entity.TransactionPhoneCredit) error {
	if detail.PhoneNumber == "" {
		return fmt.Errorf("phone number wajib diisi")
	}
	return r.masterDb.WithContext(ctx).Create(detail).Error
}

// internet & tv
func (r *transactionRepository) CreateTransactionInternetTV(ctx context.Context, detail *entity.TransactionInternetTV) error {
	if detail.CustomerName == "" {
		return fmt.Errorf("customer name wajib diisi")
	}
	return r.masterDb.WithContext(ctx).Create(detail).Error
}

// international transfer
func (r *transactionRepository) CreateTransactionInternational(ctx context.Context, detail *entity.TransactionInternational) error {
	if detail.RecipientAcc == "" || detail.RecipientBank == "" {
		return fmt.Errorf("recipient account & bank wajib diisi")
	}
	return r.masterDb.WithContext(ctx).Create(detail).Error
}
func (r *transactionRepository) GetAllTransactions(ctx context.Context, userID int64) ([]entity.Transaction, error) {
	var txns []entity.Transaction
	err := r.masterDb.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&txns).Error

	if err != nil {
		r.clogger.ErrorLogger(ctx, "GetAllTransactions", err)
		return nil, err
	}
	return txns, nil
}

// internal/repository/postgres/transaction_repository.go
func (r *transactionRepository) FindAllTransactionsByUserID(ctx context.Context, userID int64) ([]entity.Transaction, error) {
	var txns []entity.Transaction
	if err := r.masterDb.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("BankTransfer").
		Preload("Ewallet").
		Preload("PhoneCredit").
		Preload("InternetTV").
		Preload("International").
		Find(&txns).Error; err != nil {
		return nil, err
	}
	return txns, nil
}
func (r *transactionRepository) FindUserByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	var user entity.User
	if err := r.masterDb.WithContext(ctx).Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// auto unique code
func (r *transactionRepository) GenerateUniqueCode(ctx context.Context) (float64, error) {
	var expiredList []entity.Transaction
	err := r.masterDb.WithContext(ctx).
		Where("status = ? AND unique_code IS NOT NULL", "expired").
		Order("unique_code ASC").
		Find(&expiredList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	//ini belom benar ya

	// hitung total expired dan total transaksi bank transfer
	var total, expired int64
	if err := r.masterDb.WithContext(ctx).
		Model(&entity.Transaction{}).
		Count(&total).Error; err != nil {
		return 0, err
	}
	if err := r.masterDb.WithContext(ctx).
		Where("status = ? AND unique_code IS NOT NULL", "expired").
		Model(&entity.Transaction{}).
		Count(&expired).Error; err != nil {
		return 0, err
	}

	// kalau semua expired → reset ke 100
	if total > 0 && total == expired {
		return 100, nil
	}

	// kalau tidak ada expired yang reusable → ambil kode terakhir + 1
	var last entity.Transaction
	if err := r.masterDb.WithContext(ctx).
		Order("unique_code DESC").
		Limit(1).
		Find(&last).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	if last.UniqueCode == 0 {
		return 100, nil // mulai dari 100 pertama kali
	}
	if len(expiredList) > 0 {
		return expiredList[0].UniqueCode, nil
	}
	// generate baru
	return last.UniqueCode + 1, nil
}

// ExpireOldTransactions akan mengubah status transaksi pending menjadi expired jika lebih dari 6 jam
func (r *transactionRepository) ExpireOldTransactions(ctx context.Context) error {
	cutoff := getCutoffTime()
	return r.masterDb.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("status = ? AND created_at < ?", "pending", cutoff).
		Update("status", "expired").Error
}
func (r *transactionRepository) FindAllTransactionsByUserIDPaginated(
	ctx context.Context,
	userID int64,
	limit, offset int,
	search, status, txType, transactionID, startDate, endDate string,
) ([]entity.Transaction, int64, error) {

	var txs []entity.Transaction
	var total int64

	// --- Base Query ---
	query := r.masterDb.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("transactions.user_id = ?", userID).
		Preload("BankTransfer").
		Preload("Ewallet").
		Preload("PhoneCredit").
		Preload("InternetTV").
		Preload("International")

	// --- Filter tambahan ---
	if status != "" {
		query = query.Where("transactions.status = ?", status)
	}

	if txType != "" {
		query = query.Where("transactions.type = ?", txType)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("DATE(transactions.created_at) BETWEEN ? AND ?", startDate, endDate)
	}

	// --- Search di semua detail ---
	if search != "" {
		// JOIN semua detail untuk cari recipient_name di semua tabel detail
		query = query.Joins(`
			LEFT JOIN transaction_bank_transfer AS bt 
				ON bt.transaction_id = transactions.transaction_id
			LEFT JOIN transaction_ewallet AS ew 
				ON ew.transaction_id = transactions.transaction_id
			LEFT JOIN transaction_phone_credit AS pc 
				ON pc.transaction_id = transactions.transaction_id
			LEFT JOIN transaction_internet_tv AS itv 
				ON itv.transaction_id = transactions.transaction_id
			LEFT JOIN transaction_international AS intl 
				ON intl.transaction_id = transactions.transaction_id
		`).Where(`
		LOWER(transactions.transaction_id) LIKE LOWER(?) OR
			LOWER(bt.recipient_name) LIKE LOWER(?) OR
			LOWER(ew.recipient_name) LIKE LOWER(?) 
			
		
			
		`,
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	// --- Hitung total ---
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// --- Ambil data ---
	if err := query.
		Order("transactions.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}
