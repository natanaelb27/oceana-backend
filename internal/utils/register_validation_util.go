package utils

import (
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("secure_password", securePassword)
	}
}

func securePassword(fl validator.FieldLevel) bool {
	pass := fl.Field().String()
	containsNumber := false
	containsUpperCase := false
	containsLowerCase := false
	containsSymbol := false

	for _, char := range pass {
		r := rune(char)
		if unicode.IsNumber(r) {
			containsNumber = true
		}

		if unicode.IsUpper(r) {
			containsUpperCase = true
		}

		if unicode.IsLower(r) {
			containsLowerCase = true
		}

		if unicode.IsSymbol(r) || unicode.IsPunct(r) {
			containsSymbol = true
		}

	}

	return containsNumber && containsLowerCase &&
		containsUpperCase && containsSymbol && len([]rune(pass)) >= 6
}
