package config

type Smtp struct {
	Host     string `env:"SMTP_HOST" envDefault:"smtp.example.com"`
	Port     string `env:"SMTP_PORT" envDefault:"587"`
	From     string `env:"SMTP_FROM" required:"true"`
	Password string `env:"SMTP_PASSWORD" required:"true"`
}
