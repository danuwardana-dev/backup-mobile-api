package config

import "github.com/joho/godotenv"

type Server struct {
	Port   string `envconfig:"SERVER_PORT" required:"true" default:":9090"`
	Host   string `envconfig:"SERVER_HOST" required:"true" default:"127.0.0.1"`
	Domain string `envconfig:"SERVER_DOMAIN" required:"true" default:"127.0.0.1:9090"`
}

func LoadForServer(filenames ...string) Server {
	// we do not care if there is no .env file.
	_ = godotenv.Overload(filenames...)

	r := Server{}

	mustLoad("SERVER", &r)

	return r
}
