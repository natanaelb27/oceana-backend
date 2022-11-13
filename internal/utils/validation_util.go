package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func newValidationError(field string, reason string) ValidationError {
	return ValidationError{Field: field, Reason: reason}
}

func newErrorMessage(ve validator.FieldError) (msg string) {
	tag := ve.ActualTag()
	switch tag {
	case "required":
		msg = "is required."
	case "min":
		msg = fmt.Sprintf("length must be %v.", ve.Param())
	case "eqfield":
		msg = fmt.Sprintf("must be equal to %v field.", strings.ToLower(ve.Param()))
	case "email":
		msg = "have incorrect format."
	default:
		msg = "is required."
	}
	return
}

func TranslateError(err error) (errs []ValidationError) {
	if err == nil {
		return nil
	}

	// check if err is coming from gin binding errors
	if _, ok := err.(validator.ValidationErrors); !ok {
		return
	}

	// format error
	for _, e := range err.(validator.ValidationErrors) {
		msg := newValidationError(e.Field(), newErrorMessage(e))
		errs = append(errs, msg)
	}

	return
}
