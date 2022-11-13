package utils

import (
	"os"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

func NewJwtWrapper() (*JwtWrapper, error) {
	expirationHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		return nil, e.ErrAuthenticationOperationFailed
	}
	return &JwtWrapper{
		SecretKey:       os.Getenv("JWT_SECRET_KEY"),
		Issuer:          os.Getenv("JWT_ISSUER"),
		ExpirationHours: int64(expirationHours),
	}, nil

}

type JwtClaim struct {
	User *dtos.JwtUser `json:"user"`
	jwt.RegisteredClaims
}

func (j *JwtWrapper) GenerateToken(jwtUser *dtos.JwtUser) (signedToken string, err error) {

	claims := &JwtClaim{
		User: jwtUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", e.ErrAuthenticationOperationFailed
	}

	return
}

func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err = e.ErrInvalidToken
				return nil, e.ErrInvalidToken
			}
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = e.ErrAuthenticationOperationFailed
		return
	}

	return
}
