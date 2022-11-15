package handlers

import (
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
	responseData, err := h.authService.Login(request) 
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}
