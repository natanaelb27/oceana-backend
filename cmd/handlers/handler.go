package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	roleService services.RoleService
}

type HandlerConfig struct {
	RoleService services.RoleService
}

func New(c *HandlerConfig) *Handler {
	return &Handler{
		roleService: c.RoleService,
	}
}

func (h *Handler) NotFound(c *gin.Context) {
	utils.ResponseError(c, e.ErrResourceDoesNotExist)
}
