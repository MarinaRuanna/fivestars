package domain

import (
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
)

type UserCredentials struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

func NewUserCredentials(email, password string) (*UserCredentials, error) {
	creds := &UserCredentials{
		Email:    email,
		Password: password,
	}
	if err := creds.Validate(); err != nil {
		return nil, err
	}
	return creds, nil
}

func (c *UserCredentials) Validate() error {
	if err := validator.Validate(c); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}
