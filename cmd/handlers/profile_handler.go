package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) UpdateProfilePassword(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.UpdateProfilePasswordRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	result, err := h.profileService.UpdateProfilePassword(c, id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateProfilePasswordResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) ProfileCheckPasswordVerificationCode(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.ProfileCheckPasswordVerificationCodeRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	result, err := h.profileService.CheckPasswordVerificationCode(c, id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewProfileCheckPasswordVerificationCodeResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
