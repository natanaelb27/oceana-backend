package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	externalDtos "github.com/natanaelb27/oceana-backend/pkg/dtos"
)

func (h *Handler) UploadMedia(c *gin.Context) {
	var request *externalDtos.UploadMediaRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	response, err := h.mediaService.Upload(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, response, nil)
}

func (h *Handler) DeleteMedia(c *gin.Context) {
	var request *externalDtos.DeleteMediaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	response, err := h.mediaService.Delete(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, response, nil)
}
