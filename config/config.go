package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type Config struct {
	Port  string `envconfig:"PORT" default:"8080"`
	DSN   string `envconfig:"DATABASE_URL"`
	Token string `envconfig:"TG_TOKEN"`
}

var (
	config Config
	once   sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
	})
	return &config
}

func InitENV(dir string) error {
	if err := godotenv.Load(filepath.Join(dir, ".env.local")); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Print("файл .env.local не найден")
		} else {
			return err
		}
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		return err
	}
	return nil
}
