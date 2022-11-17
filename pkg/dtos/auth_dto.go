package dtos

type AuthStepUpPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type AuthStepUpPasswordResponse struct {
	UpgradeToken string `json:"upgrade_token"`
}

func NewAuthStepUpPasswordResponse(token string) (target *AuthStepUpPasswordResponse) {
	target = &AuthStepUpPasswordResponse{
		UpgradeToken: token,
	}
	return
}
