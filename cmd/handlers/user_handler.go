package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) UserUpdateProfile(c *gin.Context) {
	var request *dtos.UserProfileUpdateRequest
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	request.ID = id
	responseData, err := h.userService.UpdateProfile(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserUpdateEmail(c *gin.Context) {
	var request *dtos.UserProfileUpdateEmailRequest
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	request.ID = id

	tx := c.MustGet("tx").(*gorm.DB)

	responseData, err := h.userService.WithTx(tx).UpdateProfileEmail(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserUpdateUsername(c *gin.Context) {
	var request *dtos.UserProfileUpdateUsernameRequest
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	request.ID = id
	responseData, err := h.userService.UpdateProfileUsername(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserGetProfile(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	responseData, err := h.userService.GetProfile(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UserRegisterMerchant(c *gin.Context) {
	var request *dtos.MerchantCreateRequest
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	request.UserID = id

	tx := c.MustGet("tx").(*gorm.DB)

	responseData, err := h.merchantService.WithTx(tx).Register(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}
