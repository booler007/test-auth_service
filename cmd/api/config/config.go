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

func (c *Config) Init() {
	var ok bool
	c.Port = os.Getenv("PORT")

	if c.DB.Name, ok = os.LookupEnv("POSTGRES_DB"); !ok {
		log.Fatal("POSTGRES_DB environment variable must be set")
	}

	if c.DB.User, ok = os.LookupEnv("POSTGRES_USER"); !ok {
		log.Fatal("POSTGRES_USER environment variable must be set")
	}

	if c.DB.Pass, ok = os.LookupEnv("POSTGRES_PASSWORD"); !ok {
		log.Fatal("POSTGRES_PASSWORD environment variable must be set")
	}

	if c.Email.Host = os.Getenv("EMAIL_HOST"); !ok {
		log.Fatal("EMAIL_HOST environment variable must be set")
	}

	if c.Email.Port = os.Getenv("EMAIL_PORT"); !ok {
		log.Fatal("EMAIL_PORT environment variable must be set")
	}

	if c.Email.Account = os.Getenv("EMAIL_ACCOUNT"); !ok {
		log.Fatal("EMAIL_ACCOUNT environment variable must be set")
	}

	if c.Email.Password = os.Getenv("EMAIL_PASSWORD"); !ok {
		log.Fatal("EMAIL_PASSWORD environment variable must be set")
	}

	c.Auth.TTLAccess = os.Getenv("TTL_ACCESS")
	c.Auth.TTLRefresh = os.Getenv("TTL_REFRESH")
	c.Auth.JWTSecret = os.Getenv("JWT_SECRET")
}
