package config

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Postgres struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Dbname   string `envconfig:"POSTGRES_DATABASE" required:"true"`

	MaxConnectionLifetime time.Duration `envconfig:"DB_MAX_CONN_LIFE_TIME" default:"300s"`
	MaxOpenConnection     int           `envconfig:"DB_MAX_OPEN_CONNECTION" default:"100"`
	MaxIdleConnection     int           `envconfig:"DB_MAX_IDLE_CONNECTION" default:"10"`
}

func (pg Postgres) OpenPostgresDatabaseConnection() (*gorm.DB, error) {

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", pg.Host, pg.Port, pg.User, pg.Password, pg.Dbname),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Errorf("Error opening database connection: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("Error opening database connection: %v", err)
		return nil, err
	}
	// Lakukan Ping untuk memastikan koneksi aktif
	if err := sqlDB.Ping(); err != nil {
		log.Errorf("Error opening database connection: %v", err)
		return nil, fmt.Errorf("invalid ping database: %w", err)
	}

	return db, nil
}
func LoadPostgres(postgres Postgres) *Postgres {
	return &postgres
}
