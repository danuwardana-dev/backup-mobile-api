package transactionsvc

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/outbond/smtp"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/service/notification"
	"context"
	"fmt"
	"log"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, tx *entity.Transaction, userUUID string) (*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error)

	// detail insert
	AddTransactionBankTransfer(ctx context.Context, detail *entity.TransactionBankTransfer) error
	AddTransactionEwallet(ctx context.Context, detail *entity.TransactionEwallet) error
	AddTransactionPhoneCredit(ctx context.Context, detail *entity.TransactionPhoneCredit) error
	AddTransactionInternetTV(ctx context.Context, detail *entity.TransactionInternetTV) error
	AddTransactionInternational(ctx context.Context, detail *entity.TransactionInternational) error

	// get all
	GetAllTransactionsPaginated(
		ctx context.Context,
		userUUID string,
		limit, offset int,
		search, status, txType, transactionID, startDate, endDate string,
	) ([]entity.Transaction, int64, error)

	GenerateTransactionCode(ctx context.Context, txType string) (*CodeResponse, error)
	UpdateTransactionStatus(ctx context.Context, transactionID, status string) error
}

type transactionService struct {
	repo     postgres.TransactionRepository
	userRepo postgres.UserRepository
	notifier *notification.FirebaseNotifier
	smtp     *smtp.Smtp
}
type CodeResponse struct {
	TransactionID string   `json:"transaction_id"`
	UniqueCode    *float64 `json:"unique_code,omitempty"`
}

func (s *transactionService) GenerateTransactionCode(ctx context.Context, paymentMethod string) (*CodeResponse, error) {
	// expired-kan transaksi lama
	if err := s.repo.ExpireOldTransactions(ctx); err != nil {
		return nil, err
	}

	// generate transaction ID
	txID, err := s.repo.GenerateTransactionID(ctx)
	if err != nil {
		return nil, err
	}

	// inisialisasi uniqueCode hanya untuk bank_transfer
	var uniqueCode *float64
	if paymentMethod == "bank_transfer" {
		code, err := s.repo.GenerateUniqueCode(ctx)
		if err != nil {
			return nil, err
		}
		uniqueCode = &code
	}

	// return ke controller
	return &CodeResponse{
		TransactionID: txID,
		UniqueCode:    uniqueCode,
	}, nil
}

func NewTransactionService(repo postgres.TransactionRepository, notifier *notification.FirebaseNotifier, smtp *smtp.Smtp, userRepo postgres.UserRepository) TransactionService {
	if repo == nil {
		log.Println("[ERROR] repo nil saat init transaction service")
	}
	return &transactionService{repo: repo,
		notifier: notifier,
		smtp:     smtp,
		userRepo: userRepo,
	}
}

// transaksi utama
func (s *transactionService) CreateTransaction(ctx context.Context, tx *entity.Transaction, userUUID string) (*entity.Transaction, error) {
	// 1. Simpan transaksi ke DB
	if err := s.repo.CreateTransaction(ctx, tx, userUUID); err != nil {
		return nil, err
	}

	// 2. Ambil device user → untuk dapat FCM token
	device, err := s.repo.FindDeviceByUserUUID(ctx, userUUID)
	if err != nil {
		log.Printf("[WARN] gagal ambil device untuk user %s: %v", userUUID, err)
		// kalau device ga ada, tetap return transaksi tanpa notif
		return tx, nil
	}

	// 3. Kirim notif pakai FirebaseNotifier
	if device.FCMToken != "" {
		go func() { // kirim async biar gak nge-block response API
			if err := s.notifier.SendTransactionNotification(ctx, device.FCMToken, tx.TransactionID, tx.Status); err != nil {
				log.Printf("[ERROR] gagal kirim notif: %v", err)
			}
		}()
	} else {
		log.Printf("[WARN] user %s tidak punya FCM token", userUUID)
	}

	return tx, nil
}

type TransactionStatusUpdatePayload struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

// TransactionService.go
func (s *transactionService) UpdateTransactionStatus(ctx context.Context, transactionID, status string) error {
	// 1. Ambil transaksi
	tx, err := s.repo.FindTransactionByID(ctx, transactionID)

	if err != nil {
		return err
	}

	// 2. Update status transaksi di DB
	if err := s.repo.UpdateStatus(ctx, transactionID, status); err != nil {
		return err
	}

	// 3. Cari device FCM token user
	fcmToken, err := s.repo.GetUserFcmToken(ctx, uint(tx.UserID))
	if err != nil {
		return err
	}
	if s.notifier != nil {
		_ = s.notifier.SendTransactionNotification(ctx, fcmToken, tx.TransactionID, status)
	}
	user, err := s.userRepo.SelectUserByID(ctx, tx.UserID)
	// 3. Kirim Email
	if s.smtp != nil {
		body := fmt.Sprintf(
			"Halo %s,\n\nStatus transaksi kamu dengan ID %s sekarang adalah: %s.\n\nTerima kasih sudah menggunakan layanan kami.",
			user.FullName,
			tx.TransactionID,
			status,
		)
		fmt.Println(user.ID)
		to := []string{user.Email} // alamat email penerima

		if err := s.smtp.SendMail(ctx, to, enum.EmailSubject("Notifikasi Transaksi"), body); err != nil {
			helpers.CustomeLogger(ctx, &dto.CustomLoggerRequest{
				Error:   err.Error(),
				Remarks: " gagal kirim email transaksi",
			})
		}
	}

	return nil
}
func (s *transactionService) GetTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error) {
	tx, err := s.repo.FindTransactionByID(ctx, transactionID)
	if err != nil {
		helpers.CustomeLogger(ctx, &dto.CustomLoggerRequest{
			Error:   err.Error(),
			Remarks: "[Service][GetTransactionByID] gagal ambil transaksi",
		})
		return nil, err
	}
	return tx, nil
}

func (s *transactionService) AddTransactionBankTransfer(ctx context.Context, detail *entity.TransactionBankTransfer) error {
	if err := s.repo.CreateTransactionBankTransfer(ctx, detail); err != nil {
		helpers.CustomeLogger(ctx, &dto.CustomLoggerRequest{
			Error:   err.Error(),
			Remarks: "[Service][AddTransactionBankTransfer] gagal insert detail",
		})

	}
	return nil
}

func (s *transactionService) AddTransactionEwallet(ctx context.Context, detail *entity.TransactionEwallet) error {
	if err := s.repo.CreateTransactionEwallet(ctx, detail); err != nil {
		helpers.CustomeLogger(ctx, &dto.CustomLoggerRequest{
			Error:   err.Error(),
			Remarks: "[Service][AddTransactionEwallet] gagal insert detail",
		})

	}
	return nil
}

func (s *transactionService) AddTransactionPhoneCredit(ctx context.Context, detail *entity.TransactionPhoneCredit) error {
	if err := s.repo.CreateTransactionPhoneCredit(ctx, detail); err != nil {
		helpers.CustomeLogger(ctx, &dto.CustomLoggerRequest{
			Error:   err.Error(),
			Remarks: "[Service][AddTransactionPhoneCredit] gagal insert detail",
		})

	}
	return nil
}

func (s *transactionService) AddTransactionInternetTV(ctx context.Context, detail *entity.TransactionInternetTV) error {
	if err := s.repo.CreateTransactionInternetTV(ctx, detail); err != nil {
		helpers.CustomeLogger(ctx, &dto.CustomLoggerRequest{
			Error:   err.Error(),
			Remarks: "[Service][AddTransactionInternetTV] gagal insert detail",
		})

	}
	return nil
}

func (s *transactionService) AddTransactionInternational(ctx context.Context, detail *entity.TransactionInternational) error {
	if err := s.repo.CreateTransactionInternational(ctx, detail); err != nil {
		helpers.CustomeLogger(ctx, &dto.CustomLoggerRequest{
			Error:   err.Error(),
			Remarks: "[Service][AddTransactionInternational] gagal insert detail",
		})

	}
	return nil
}

// ✅ fix ambil transaksi by userUUID
func (s *transactionService) GetAllTransactionsPaginated(
	ctx context.Context,
	userUUID string,
	limit, offset int,
	search, status, txType, transactionID, startDate, endDate string,
) ([]entity.Transaction, int64, error) {

	user, err := s.repo.FindUserByUUID(ctx, userUUID)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.FindAllTransactionsByUserIDPaginated(ctx, user.ID, limit, offset, search, status, txType, transactionID, startDate, endDate)
}
