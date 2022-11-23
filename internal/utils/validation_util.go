package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/go-playground/validator/v10"
)

func newErrorMessage(ve validator.FieldError) (msg string) {
	tag := ve.ActualTag()
	switch tag {
	case "required":
		msg = "is required."
	case "min":
		msg = fmt.Sprintf("length must be at least %v.", ve.Param())
	case "len":
		msg = fmt.Sprintf("length must be %v.", ve.Param())
	case "eqfield":
		msg = fmt.Sprintf("must be equal to %v field.", strings.ToLower(ve.Param()))
	case "email":
		msg = "have incorrect format."
	case "secure_password":
		msg = "length must be at least 6, must be alphanumeric, contains uppercase, lowercase, and symbol."
	default:
		msg = "is incorrect."
	}
	return
}

func TranslateError(err error) (message string, responseErr string, code int) {
	switch {
	// case if error is validator error
	case errors.As(err, &validator.ValidationErrors{}):
		for i, e := range err.(validator.ValidationErrors) {
			message += fmt.Sprintf("%s %s", e.Field(), newErrorMessage(e))
			if i+1 < len(err.(validator.ValidationErrors)) {
				message += " "
			}
		}
		responseErr = http.StatusText(http.StatusBadRequest)
		code = http.StatusBadRequest
		return

	// case if an io.EOF error is returned by Decode() if the request body is empty
	// in some circumstances Decode() may also return an
	// io.ErrUnexpectedEOF error for syntax errors in the JSON
	case errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF):
		err = e.ErrEmptyRequstBody

	// case if the request body is too large.
	case err.Error() == "http: request body too large":
		err = e.ErrRequestBodyTooLarge
	}

	message = err.Error()
	responseErr = http.StatusText(e.GetStatusCode(err))
	code = e.GetStatusCode(err)

	return
}

func IsPasswordContainsUsername(password string, username string) bool {
	return !strings.Contains(strings.ToLower(password), strings.ToLower(username))
}
