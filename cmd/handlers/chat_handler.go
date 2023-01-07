package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	internalDtos "github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func (h *Handler) ChatGetAllChatList(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	results, err := h.chatService.GetAllChatList(userID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}

func (h *Handler) ChatGetChat(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	toUserID := c.Param("to-user-id")
	touserIDUUID, err := uuid.Parse(toUserID)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	results, err := h.chatService.GetChat(userID, touserIDUUID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}

func (h *Handler) ChatCreate(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *internalDtos.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	request.FromUserID = userID

	err := h.chatService.Create(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, true, nil)
}
