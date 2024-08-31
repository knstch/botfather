package repository

import "gorm.io/gorm"

type Bot struct {
	gorm.Model
	Token  string `gorm:"unique"`
	Name   string `gorm:"unique"`
	User   User
	UserId uint
}

type User struct {
	gorm.Model
	Email      *string `gorm:"unique"`
	ExternalId int     `gorm:"unique"`
	Verified   bool    `gorm:"default:false"`
}

func (r DbRepository) AutoMigrate() error {
	if err := r.db.AutoMigrate(
		&Bot{},
		&User{},
	); err != nil {
		return err
	}

	return nil
}
