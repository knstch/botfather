package register

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Register struct {
	registrator Registrator
}

type Registrator interface {
	RegisterUser(ctx context.Context, email string, externalId int) error
	CheckUser(ctx context.Context, externalId int) (bool, error)
	ConfirmEmail(ctx context.Context, email string, externalId int64) error
}

var (
	ErrBadEmail = errors.New("bad email")
)

func NewRegisterService(registrator Registrator) *Register {
	return &Register{
		registrator,
	}
}

func (r Register) CheckUser(ctx context.Context, externalId int) (bool, error) {
	verified, err := r.registrator.CheckUser(ctx, externalId)
	if err != nil {
		return false, err
	}

	return verified, nil
}

func (r Register) ConfirmEmail(ctx context.Context, email string, externalId int64) error {
	if err := r.emailValidator(email); err != nil {
		return ErrBadEmail
	}

	if err := r.registrator.ConfirmEmail(ctx, email, externalId); err != nil {
		return err
	}

	return nil
}

func (r Register) emailValidator(email string) error {
	if err := validation.Validate(email, is.Email); err != nil {
		return ErrBadEmail
	}

	return nil
}
