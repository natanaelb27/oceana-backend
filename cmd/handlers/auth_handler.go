package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UserCreate(c *gin.Context) {
	var request *dtos.UserCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, errors.ErrInvalidEmailAddress)
		return
	}
	responseData, err := h.authService.Create(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}
