package handlers

import (
	"net/http"
	"os"
	"strconv"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) UserCreate(c *gin.Context) {
	var request *externalDtos.UserCreateRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	responseData, err := h.authService.Register(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserLogin(c *gin.Context) {
	var request *externalDtos.UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	userLoginData, err := h.authService.Login(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	userLoginResponse := &externalDtos.UserLoginResponse{
		Role: userLoginData.Role,
	}

	err = h.setLoginCookies(c, userLoginData)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, userLoginResponse, nil)
}

func (h *Handler) UserRegisterWithGoogle(c *gin.Context) {
	var request *externalDtos.UserRegisterWithGoogleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	responseData, err := h.authService.RegisterWithGoogle(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserLoginWithGoogle(c *gin.Context) {
	var request *externalDtos.UserLoginWithGoogleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	userLoginData, err := h.authService.LoginWithGoogle(c, request)
	if err != nil {
		if err.Error() == e.ErrRegisterGoogle.Error() {
			utils.Response(c, http.StatusContinue, nil, nil, nil, http.StatusText(http.StatusContinue))
			return
		}
		utils.ResponseError(c, err)
		return
	}
	userLoginResponse := &externalDtos.UserLoginResponse{
		Role: userLoginData.Role,
	}

	err = h.setLoginCookies(c, userLoginData)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, userLoginResponse, nil)
}

func (h *Handler) setLoginCookies(c *gin.Context, userLoginData *internalDtos.UserLoginData) error {
	jwtExpirationHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		return err
	}
	accessTokenExpirationSeconds := 3600 * jwtExpirationHours
	c.SetCookie("access_token", userLoginData.AccessToken, accessTokenExpirationSeconds, "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("refresh_token", userLoginData.RefreshToken, int(userLoginData.ExpireDuration), "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("access_token_verifier", "", accessTokenExpirationSeconds, "/", os.Getenv("DOMAIN"), false, false)
	return nil
}

func (h *Handler) StepUpPassword(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *externalDtos.AuthStepUpPasswordRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	result, err := h.authService.StepUpPassword(c, id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := externalDtos.NewAuthStepUpPasswordResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserVerifyResetPassword(c *gin.Context) {
	var request *externalDtos.UserVerifyResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	result, err := h.authService.SendResetPasswordEmail(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	responseData := externalDtos.NewUserVerifyResetPasswordResponse(result)
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserResetPassword(c *gin.Context) {
	var request *externalDtos.ResetPasswordRequest
	resetPasswordID := c.Param("reset-password-id")
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	responseData, err := h.authService.ResetPassword(c, resetPasswordID, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) RefreshAccessToken(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrUnauthorized)
		return
	}

	accessToken, err := h.authService.RefreshAccessToken(c, id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	jwtExpirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	accessTokenExpirationSeconds := 3600 * jwtExpirationHours
	c.SetCookie("access_token", accessToken, accessTokenExpirationSeconds, "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("access_token_verifier", "", accessTokenExpirationSeconds, "/", os.Getenv("DOMAIN"), false, false)

	utils.ResponseSuccess(c, nil, nil)
}
