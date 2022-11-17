package handlers

import (
	"os"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) UserCreate(c *gin.Context) {
	var request *dtos.UserCreateRequest
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
	var request *dtos.UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	userLoginData, err := h.authService.Login(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	userLoginResponse := &dtos.UserLoginResponse{
		Role: userLoginData.Role,
	}
	expirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	accessTokenExpirationSeconds := 3600 * expirationHours
	refreshTokenExpirationSeconds := int(time.Time(userLoginData.ExpiresAt).Unix()) - int(time.Now().Unix())
	c.SetCookie("access_token", userLoginData.AccessToken, accessTokenExpirationSeconds, "/", "localhost", false, true)
	c.SetCookie("refresh_token", userLoginData.RefreshToken, refreshTokenExpirationSeconds, "/", "localhost", false, true)
	utils.ResponseSuccess(c, userLoginResponse, nil)
}

func (h *Handler) StepUpPassword(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.AuthStepUpPasswordRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	result, err := h.authService.StepUpPassword(c, id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewAuthStepUpPasswordResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserVerifyResetPassword(c *gin.Context) {
	var request *dtos.UserVerifyResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	err := h.authService.SendResetPasswordEmail(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, nil, nil)
}
