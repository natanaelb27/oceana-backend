package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UploadImage(c *gin.Context) {
	var request *externalDtos.UploadImageRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	response, err := h.imageService.Upload(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, response, nil)
}
