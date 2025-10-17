package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebaseNotifier struct {
	client *messaging.Client
}

// InitFirebaseNotifier inisialisasi Firebase Admin SDK
func InitFirebaseNotifier(ctx context.Context, credFile string) (*FirebaseNotifier, error) {
	// load service account key dari file json
	data, err := os.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read firebase credentials file: %w", err)
	}
	var creds struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(data, &creds); err == nil {
		fmt.Println("[DEBUG] Firebase project_id dari credFile:", creds.ProjectID)
	}

	opt := option.WithCredentialsFile(credFile)

	// inisialisasi firebase app
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to init firebase app: %w", err)
	}

	// ambil messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init messaging client: %w", err)
	}
	fmt.Println("[DEBUG] Firebase app initialized with credentials:", credFile)
	return &FirebaseNotifier{client: client}, nil
}

// SendPushNotification kirim notif ke device tertentu
func (f *FirebaseNotifier) SendPushNotification(token, title, body, status string) error {
	ctx := context.Background() // jangan pakai ctx dari HTTP request
	fmt.Printf("[DEBUG] Preparing to send push notif\n")
	fmt.Printf("        Token (len=%d): %s\n", len(token), token)
	fmt.Printf("        Title: %s\n", title)
	fmt.Printf("        Body: %s\n", body)
	fmt.Printf("        Status: %s\n", status)

	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"click_action":       "FLUTTER_NOTIFICATION_CLICK",
			"transaction_status": status,
		},
	}
	fmt.Printf("[DEBUG] Message Payload: %+v\n", msg)

	resp, err := f.client.Send(ctx, msg)
	if err != nil {
		fmt.Printf("[ERROR] Failed to send message: %v\n", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Println("Successfully sent message:", resp)
	return nil
}

// SendTransactionNotification kirim notif berdasarkan status transaksi
func (f *FirebaseNotifier) SendTransactionNotification(ctx context.Context, fcmToken string, transactionID string, status string) error {
	var body string
	var title string
	switch status {
	case "pending":
		title = "Transfer Pending"
		body = "Selesaikan pembayaranmu yuk"
	case "success":
		title = "Transfer Berhasil 🎉"
		body = "Pembayaranmu berhasil 🎉"
	case "failed":
		title = "Transfer Gagal ❌"
		body = "Pembayaranmu gagal ❌"
	case "canceled":
		title = "Transfer Dibatalkan"
		body = "Transaksi kamu dibatalkan"
	case "expired":
		title = "Transfer Expired"
		body = "Transaksimu sudah kadaluarsa ⏰"
	default:
		title = "Transfer Update"
		body = fmt.Sprintf("Status transaksi %s: %s", transactionID, status)
	}

	err := f.SendPushNotification(
		fcmToken,
		title,
		body,
		status,
	)
	if err != nil {
		log.Printf("[ERROR] gagal kirim notif transaksi %s (status: %s): %v", transactionID, status, err)
		return err
	}
	log.Printf("[INFO] Notif transaksi %s status %s berhasil dikirim", transactionID, status)
	return nil
}
