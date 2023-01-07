package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
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

func NewRefreshTokenJwtWrapper() (*JwtWrapper, error) {
	expirationDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_DURATION"))
	expirationHours := expirationDays * 24
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

func (j *JwtWrapper) GenerateToken(jwtUser *dtos.JwtUser, customJwtWrapper *JwtWrapper) (signedToken string, err error) {
	jw := j.processCustomJwtWrapper(customJwtWrapper)

	claims := &JwtClaim{
		User: jwtUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jw.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(jw.ExpirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(jw.SecretKey))
	if err != nil {
		return "", e.ErrAuthenticationOperationFailed
	}

	return
}

func (j *JwtWrapper) ValidateToken(signedToken string, customJwtWrapper *JwtWrapper) (claims *JwtClaim, err error) {
	jw := j.processCustomJwtWrapper(customJwtWrapper)

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err = e.ErrInvalidToken
				return nil, e.ErrInvalidToken
			}
			return []byte(jw.SecretKey), nil
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

func (j *JwtWrapper) processCustomJwtWrapper(customJwtWrapper *JwtWrapper) *JwtWrapper {
	jw := &JwtWrapper{
		SecretKey:       j.SecretKey,
		Issuer:          j.Issuer,
		ExpirationHours: j.ExpirationHours,
	}

	if customJwtWrapper != nil {
		if customJwtWrapper.ExpirationHours != 0 {
			jw.ExpirationHours = customJwtWrapper.ExpirationHours
		}
		if customJwtWrapper.Issuer != "" {
			jw.Issuer = customJwtWrapper.Issuer
		}
		if customJwtWrapper.SecretKey != "" {
			jw.SecretKey = customJwtWrapper.SecretKey
		}
	}
	return jw
}
