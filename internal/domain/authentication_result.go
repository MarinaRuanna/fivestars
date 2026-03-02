package domain

import (
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
)

type AuthenticationResult struct {
	UserID string `validate:"required"`
	Name   string
	Token  string `validate:"required"`
}

func NewAuthenticationResult(userID, name, token string) (*AuthenticationResult, error) {
	result := &AuthenticationResult{
		UserID: userID,
		Token:  token,
	}

	if name != "" {
		result.Name = name
	}

	if err := result.Validate(); err != nil {
		return nil, err
	}
	return result, nil
}

func (a *AuthenticationResult) Validate() error {
	if err := validator.Validate(a); err != nil {
		return customerror.NewValidationError("invalid authentication result: " + err.Error())
	}
	return nil
}
