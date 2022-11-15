package handlers

import (
	"os"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
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
