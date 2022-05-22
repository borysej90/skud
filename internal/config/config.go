package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Netflix/go-env"
)

type Config struct {
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	HTTPPort   string `env:"HTTP_PORT"`
}

func NewConfig() (Config, error) {
	var cfg Config
	cfg.HTTPPort = "8080"
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	fields := make([]string, 0)
	if c.DBName == "" {
		fields = append(fields, "DB_NAME")
	}
	if c.DBUser == "" {
		fields = append(fields, "DB_USER")
	}
	if c.DBPassword == "" {
		fields = append(fields, "DB_PASSWORD")
	}
	if c.DBHost == "" {
		fields = append(fields, "DB_HOST")
	}
	if c.DBPort == "" {
		fields = append(fields, "DB_PORT")
	}
	if c.HTTPPort == "" {
		fields = append(fields, "HTTP_PORT")
	}
	if len(fields) == 0 {
		return nil
	}
	msg := "validation failed: missing fields "
	msg += strings.Join(fields, ", ")
	return errors.New(msg)
}

func (c Config) DBUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
