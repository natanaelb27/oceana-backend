package utils

import (
	"encoding/json"
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
	fmt.Println(tag)
	switch tag {
	case "required":
		msg = "is required."
	case "gte":
		msg = fmt.Sprintf("must be greater than equal to %v.", ve.Param())
	case "lte":
		msg = fmt.Sprintf("must be less than equal to %v.", ve.Param())
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
	var unmarshalTypeError *json.UnmarshalTypeError

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

	case errors.As(err, &unmarshalTypeError):
		err = fmt.Errorf("request body contains an invalid value for the %q field, type should be %v, instead got %v", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)

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
