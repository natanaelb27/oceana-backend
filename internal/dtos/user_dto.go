package dtos

type UserLoginData struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	ExpireDuration int    `json:"expire_duration"`
	Role           string `json:"role"`
}
