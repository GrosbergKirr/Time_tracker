package internal

import (
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig
	HttpConfig
}

type DatabaseConfig struct {
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Adress   string `env:"DB_ADDRESS"`
	Database string `env:"DB_DATABASE"`
	Mode     string `env:"DB_MODE"`
}

type HttpConfig struct {
	Host        string        `env:"HTTP_HOST"`
	Port        string        `env:"HTTP_PORT"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT"`
}

func SetupConfig(log *slog.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
		return nil
	}
	var cfg Config
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		log.Error("can't read config: %s", err)
		return nil
	}
	return &cfg

}
