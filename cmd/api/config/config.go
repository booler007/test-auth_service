package config

import (
	"log"
	"os"
)

type Config struct {
	Port  string `default:"8080"`
	Auth  *auth
	DB    *db
	Email *Email
}

type db struct {
	Pass string
	Name string
	User string
}

type Email struct {
	Account  string
	Password string
	Host     string
	Port     string
}

type auth struct {
	TTLAccess  string `default:"30m"`
	TTLRefresh string `default:"720h"`
	JWTSecret  string `default:"secretWordJWT"`
}

func Init() *Config {
	var (
		cfg *Config
		ok  bool
	)
	cfg.Port = os.Getenv("PORT")

	if cfg.DB.Name, ok = os.LookupEnv("POSTGRES_DB"); !ok {
		log.Fatal("POSTGRES_DB environment variable must be set")
	}

	if cfg.DB.User, ok = os.LookupEnv("POSTGRES_USER"); !ok {
		log.Fatal("POSTGRES_USER environment variable must be set")
	}

	if cfg.DB.Pass, ok = os.LookupEnv("POSTGRES_PASSWORD"); !ok {
		log.Fatal("POSTGRES_PASSWORD environment variable must be set")
	}

	if cfg.Email.Host, ok = os.LookupEnv("EMAIL_HOST"); !ok {
		log.Fatal("EMAIL_HOST environment variable must be set")
	}

	if cfg.Email.Port, ok = os.LookupEnv("EMAIL_PORT"); !ok {
		log.Fatal("EMAIL_PORT environment variable must be set")
	}

	if cfg.Email.Account, ok = os.LookupEnv("EMAIL_ACCOUNT"); !ok {
		log.Fatal("EMAIL_ACCOUNT environment variable must be set")
	}

	if cfg.Email.Password, ok = os.LookupEnv("EMAIL_PASSWORD"); !ok {
		log.Fatal("EMAIL_PASSWORD environment variable must be set")
	}

	cfg.Auth.TTLAccess = os.Getenv("TTL_ACCESS")
	cfg.Auth.TTLRefresh = os.Getenv("TTL_REFRESH")
	cfg.Auth.JWTSecret = os.Getenv("JWT_SECRET")

	return cfg
}
