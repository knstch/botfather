package repository

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tg_botfather/config"
	"time"
)

type DbRepository struct {
	db *gorm.DB
}

var (
	ErrAccountAlreadyExists = errors.New("account with this id/email already exists")
)

func InitConnect(cfg config.Config) (*DbRepository, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)

	return &DbRepository{
		db: db,
	}, nil
}
