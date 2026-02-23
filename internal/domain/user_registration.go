package domain

import (
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
)

type UserRegistration struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
	Name     string
}

func NewUserRegistration(email, password, name string) (*UserRegistration, error) {
	reg := &UserRegistration{
		Email:    email,
		Password: password,
		Name:     name,
	}
	if err := reg.Validate(); err != nil {
		return nil, err
	}
	return reg, nil
}

func (r *UserRegistration) Validate() error {
	if err := validator.Validate(r); err != nil {
		return customerror.NewValidationError(err.Error())

	}
	return nil
}
