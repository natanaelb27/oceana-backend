package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

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
		msg = "is incorrect."
	}
	return
}

func TranslateError(err error) (errors string) {
	if err == nil {
		return "-"
	}

	// check if err is coming from gin binding errors
	if _, ok := err.(validator.ValidationErrors); !ok {
		return "-"
	}

	// format error
	for i, e := range err.(validator.ValidationErrors) {
		errors += fmt.Sprintf("%s %s", e.Field(), newErrorMessage(e))
		if i+1 < len(err.(validator.ValidationErrors)) {
			errors += " "
		}
	}

	return
}
