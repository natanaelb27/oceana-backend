package utils

import (
	"math/rand"
	"strings"
	"time"
)

func randomString(n int, alphabet []rune) string {
	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}

func RandomVerificationCode() string {
	rand.Seed(time.Now().UnixNano())

	var alphabet []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	rs := randomString(6, alphabet)

	return rs
}
