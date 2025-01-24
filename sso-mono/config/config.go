package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = ""

type Configuration struct {
	HTTPServer
	Database
}

type InitialData struct {
	Username     string `envconfig:"ADMIN_USERNAME" default:"admin"`
	Password     string `envconfig:"ADMIN_PASSWORD" default:"admin"`
	Email        string `envconfig:"ADMIN_EMAIL" default:"admin@example.com"`
	ClientName   string `envconfig:"SSO_CLIENT_NAME" default:"admin"`
	ClientSecret string `envconfig:"SSO_CLIENT_SECRET" default:"S3cr3t"`
	RedirectURI  string `envconfig:"SSO_REDIRECT_URI" default:"http://localhost:8080/callback"`
}
type Database struct {
	Host     string `envconfig:"DATABASE_HOST" default:"localhost"`
	Port     int    `envconfig:"DATABASE_PORT" default:"5432"`
	User     string `envconfig:"DATABASE_USER" default:"satish"`
	Password string `envconfig:"DATABASE_PASSWORD" default:"satish"`
	Name     string `envconfig:"DATABASE_NAME" default:"authentication"`
	SSLMode  string `envconfig:"DATABASE_SSL_MODE" default:"disable"`
	Initial  InitialData
}

type HTTPServer struct {
	Debug bool `envconfig:"HTTP_SERVER_DEBUG" default:"true"`
	Port  int  `envconfig:"PORT" default:"8080"`

	IdleTimeout  time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" default:"60s"`
	ReadTimeout  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"1s"`
	WriteTimeout time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"2s"`
}

func Load() (Configuration, error) {
	var cfg Configuration
	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
