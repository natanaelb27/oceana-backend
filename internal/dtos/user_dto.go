package dtos

import "time"

type UserLoginData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Role         string    `json:"role"`
}
