package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) RoleGetAll(c *gin.Context) {
	results, err := h.roleService.GetAll()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*dtos.RoleGetResponse = make([]*dtos.RoleGetResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewRoleGetResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) RoleGetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.roleService.GetById(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewRoleGetResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) RoleCreate(c *gin.Context) {
	var request *dtos.RoleCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel()

	result, err := h.roleService.Create(model)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewRoleCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) RoleUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	var request *dtos.RoleCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel()

	result, err := h.roleService.Update(id, model)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewRoleCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) RoleDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.roleService.Delete(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	responseData := dtos.NewRoleDeleteResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
