package config

type FirebaseConfig struct {
	CredentialsFile string `envconfig:"FIREBASE_CREDENTIALS"`
}
