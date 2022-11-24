package dtos

import "time"

type UserLoginData struct {
	AccessToken    string        `json:"access_token"`
	RefreshToken   string        `json:"refresh_token"`
	ExpireDuration time.Duration `json:"expire_duration"`
	Role           string        `json:"role"`
}
