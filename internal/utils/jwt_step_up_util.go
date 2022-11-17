package utils

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
)

type JwtStepUp struct {
	Scope string
}

type JwtStepUpClaim struct {
	Data *dtos.JwtStepUp `json:"data"`
	jwt.RegisteredClaims
}

func (j *JwtWrapper) JwtStepUpGenerateToken(jwtData *dtos.JwtStepUp) (signedToken string, err error) {
	claims := &JwtStepUpClaim{
		Data: jwtData,
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

func (j *JwtWrapper) JwtStepUpValidateToken(signedToken string) (claims *JwtStepUpClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtStepUpClaim{},
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

	claims, ok := token.Claims.(*JwtStepUpClaim)
	if !ok {
		err = e.ErrAuthenticationOperationFailed
		return
	}

	return
}
