package config

import "time"

type Jwt struct {
	SecreteKey        string        `envconfig:"JWT_SECRET_KEY" required:"true"`
	PublicKey         string        `envconfig:"JWT_PUBLIC_KEY" required:"true"`
	RefreshSecreteKey string        `envconfig:"JWT_REFRESH_SECRET_KEY" required:"true"`
	RefreshPublicKey  string        `envconfig:"JWT_REFRESH_PUBLIC_KEY" required:"true"`
	Expiration        time.Duration `envconfig:"JWT_EXPIRATION" required:"true"`
	RefreshExpiration time.Duration `envconfig:"JWT_REFRESH_EXPIRATION" required:"true"`
}
