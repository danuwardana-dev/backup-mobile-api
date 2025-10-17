package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
)

type Root struct {
	Firebase FirebaseConfig
	App      App
	Postgres Postgres
	Redis    Redis
	Jwt      Jwt
	Smtp     Smtp
	Server   Server
	Verihubs Verihubs
	Minio    Minio
}

func mustLoad(prefix string, spec interface{}) {
	err := envconfig.Process(prefix, spec)
	if err != nil {
		panic(err)
	}
}
func Load(filenames ...string) Root {
	err := godotenv.Overload(filenames...)
	if err != nil {
		log.Errorf(err.Error())
	}

	r := Root{
		App:      App{},
		Postgres: Postgres{},
		Redis:    Redis{},
		Jwt:      Jwt{},
		Smtp:     Smtp{},
		Server:   Server{},
		Firebase: FirebaseConfig{},

		Verihubs: Verihubs{},
		Minio:    Minio{},
	}
	mustLoad("FIREBASE", &r.Firebase)
	mustLoad("SERVER", &r.Server)
	mustLoad("APP", &r.App)
	mustLoad("POSTGRES", &r.Postgres)
	mustLoad("REDIS", &r.Redis)
	mustLoad("JWT", &r.Jwt)
	mustLoad("SMTP", &r.Smtp)
	mustLoad("VERIHUBS", &r.Verihubs)
	mustLoad("MINIO", &r.Minio)

	return r
}
