package entity

import "time"

// ========================
// TABEL UTAMA
// ========================
type Transaction struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	TransactionID string    `gorm:"unique;not null;type:varchar(50)" json:"transaction_id" db:"transaction_id"`
	UserID        int64     `gorm:"not null" json:"user_id" db:"user_id"`
	Type          string    `gorm:"not null;type:varchar(50)" json:"type" db:"type"`                     // bank_transfer, ewallet, phone_credit, etc
	PaymentMethod string    `gorm:"not null;type:varchar(30)" json:"payment_method" db:"payment_method"` // bank_transfer | virtual_account
	Description   string    `gorm:"type:text" json:"description" db:"description"`
	Nominal       float64   `gorm:"not null" json:"nominal" db:"nominal"`
	AdminFee      float64   `gorm:"default:0" json:"admin_fee" db:"admin_fee"`
	UniqueCode    float64   `gorm:"default:0" json:"unique_code" db:"unique_code"`
	Total         float64   `gorm:"not null" json:"total" db:"total"`
	Status        string    `gorm:"not null;type:varchar(30)" json:"status" db:"status"` // pending, success, failed
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at" db:"updated_at"`
	ExpiredAt     time.Time `gorm:"column:expired_at" json:"expired_at"`

	// RELASI DETAIL TRANSAKSI
	BankTransfer  *TransactionBankTransfer  `gorm:"foreignKey:TransactionID;references:TransactionID" json:"bank_transfer,omitempty"`
	Ewallet       *TransactionEwallet       `gorm:"foreignKey:TransactionID;references:TransactionID" json:"ewallet,omitempty"`
	PhoneCredit   *TransactionPhoneCredit   `gorm:"foreignKey:TransactionID;references:TransactionID" json:"phone_credit,omitempty"`
	InternetTV    *TransactionInternetTV    `gorm:"foreignKey:TransactionID;references:TransactionID" json:"internet_tv,omitempty"`
	International *TransactionInternational `gorm:"foreignKey:TransactionID;references:TransactionID" json:"international,omitempty"`
}

func (Transaction) TableName() string { return "transactions" }

// ========================
// DETAIL PER TIPE TRANSAKSI
// ========================

// Bank Transfer
type TransactionBankTransfer struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	TransactionID string `gorm:"not null;index" json:"transaction_id" db:"transaction_id"`
	RecipientName string `gorm:"not null;type:varchar(100)" json:"recipient_name" db:"recipient_name"`
	AccountNumber string `gorm:"size:50" json:"account_number"`
	ImageURL      string `gorm:"size:255" json:"image_url"`
	BankName      string `gorm:"not null;type:varchar(100)" json:"bank_name" db:"bank_name"`
	Notes         string `gorm:"type:text" json:"notes" db:"notes"`
}

func (TransactionBankTransfer) TableName() string { return "transaction_bank_transfer" }

// E-Wallet
type TransactionEwallet struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	TransactionID string `gorm:"not null;index" json:"transaction_id" db:"transaction_id"`
	RecipientName string `gorm:"not null;type:varchar(100)" json:"recipient_name" db:"recipient_name"`
	AccountNumber string `gorm:"size:50" json:"account_number"`
	ImageURL      string `gorm:"size:255" json:"image_url"`
	EwalletName   string `gorm:"not null;type:varchar(50)" json:"ewallet_name" db:"ewallet_name"`
}

func (TransactionEwallet) TableName() string { return "transaction_ewallet" }

// Phone Credit
type TransactionPhoneCredit struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	TransactionID string `gorm:"not null;index" json:"transaction_id" db:"transaction_id"`
	PhoneNumber   string `gorm:"not null;type:varchar(20)" json:"phone_number" db:"phone_number"`
	ProductName   string `gorm:"not null;type:varchar(100)" json:"product_name" db:"product_name"`
}

func (TransactionPhoneCredit) TableName() string { return "transaction_phone_credit" }

// Internet & TV
type TransactionInternetTV struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	TransactionID string `gorm:"not null;index" json:"transaction_id" db:"transaction_id"`
	CustomerName  string `gorm:"not null;type:varchar(100)" json:"customer_name" db:"customer_name"`
	Description   string `gorm:"type:text" json:"description" db:"description"`
}

func (TransactionInternetTV) TableName() string { return "transaction_internet_tv" }

// International Transfer
type TransactionInternational struct {
	ID             int64   `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	TransactionID  string  `gorm:"not null;index" json:"transaction_id" db:"transaction_id"`
	RecipientFirst string  `gorm:"not null;type:varchar(50)" json:"recipient_first_name" db:"recipient_first_name"`
	RecipientLast  string  `gorm:"not null;type:varchar(50)" json:"recipient_last_name" db:"recipient_last_name"`
	RecipientBank  string  `gorm:"not null;type:varchar(100)" json:"recipient_bank" db:"recipient_bank"`
	RecipientAcc   string  `gorm:"not null;type:varchar(50)" json:"recipient_account" db:"recipient_account"`
	SenderName     string  `gorm:"not null;type:varchar(100)" json:"sender_name" db:"sender_name"`
	SenderBank     string  `gorm:"not null;type:varchar(100)" json:"sender_bank" db:"sender_bank"`
	Country        string  `gorm:"not null;type:varchar(50)" json:"country" db:"country"`
	Currency       string  `gorm:"not null;type:varchar(10)" json:"currency" db:"currency"`
	TransferMethod string  `gorm:"not null;type:varchar(50)" json:"transfer_method" db:"transfer_method"`
	YouSend        float64 `gorm:"not null" json:"you_send" db:"you_send"`
	RecipientGets  float64 `gorm:"not null" json:"recipient_gets" db:"recipient_gets"`
}

func (TransactionInternational) TableName() string { return "transaction_international" }
