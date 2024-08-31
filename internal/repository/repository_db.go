package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrUserNotExists = errors.New("user not exists")
)

func (r DbRepository) RegisterUser(ctx context.Context, email string, externalId int) error {
	tx := r.db.WithContext(ctx).Begin()

	user := &User{
		Email:      &email,
		ExternalId: externalId,
	}

	if err := tx.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAccountAlreadyExists
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r DbRepository) CheckUser(ctx context.Context, externalId int) (bool, error) {
	tx := r.db.WithContext(ctx).Begin()

	user := new(User)
	if err := tx.First(&user, "external_id = ?", externalId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &User{
				ExternalId: externalId,
			}

			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				return false, err
			}

			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				return false, err
			}

			return false, ErrUserNotExists
		}
		return false, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return user.Verified, nil
}

func (r DbRepository) ConfirmEmail(ctx context.Context, email string, externalId int64) error {
	tx := r.db.WithContext(ctx).Begin()

	if err := tx.Model(&User{}).Where("external_id = ?", externalId).Update("email", email).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
