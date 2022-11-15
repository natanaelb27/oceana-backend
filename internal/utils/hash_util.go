package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	hashedPassword := string(bytes)
	return hashedPassword, nil
}

func CheckHash(password string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
