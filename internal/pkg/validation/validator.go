package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Failure struct {
	Field   string
	Message string
}

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (v *Validator) Struct(value any) []Failure {
	if v == nil || v.validate == nil {
		v = New()
	}
	if err := v.validate.Struct(value); err != nil {
		return Normalize(err)
	}
	return nil
}

func Normalize(err error) []Failure {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []Failure{{
			Field:   "request",
			Message: "invalid request",
		}}
	}

	failures := make([]Failure, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		failures = append(failures, Failure{
			Field:   fieldErr.Field(),
			Message: fieldErr.Tag(),
		})
	}
	return failures
}
