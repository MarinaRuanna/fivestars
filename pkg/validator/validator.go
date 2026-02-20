package validator

import (
	"errors"

	playground "github.com/go-playground/validator/v10"
	"go.uber.org/multierr"
)

func Validate(value interface{}) error {
	v := playground.New()

	err := v.Struct(value)
	if err != nil {
		var validatorErrors playground.ValidationErrors
		if errors.As(err, &validatorErrors) {
			return handlerValidationErrors(validatorErrors)
		}
		return err
	}
	return nil
}

func handlerValidationErrors(err playground.ValidationErrors) error {
	var errorsMessage error
	for _, ex := range err {
		errorsMessage = multierr.Append(errorsMessage, ex)
	}
	return errorsMessage
}
