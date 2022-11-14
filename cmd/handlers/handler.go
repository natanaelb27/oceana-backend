package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	roleService     services.RoleService
	categoryService services.CategoryService
}

type HandlerConfig struct {
	RoleService     services.RoleService
	CategoryService services.CategoryService
}

func New(c *HandlerConfig) *Handler {
	return &Handler{
		roleService:     c.RoleService,
		categoryService: c.CategoryService,
	}
}

func (h *Handler) Root(c *gin.Context) {
	utils.ResponseSuccess(c, "Oceana API", nil)
}

func (h *Handler) NotFound(c *gin.Context) {
	utils.ResponseError(c, e.ErrResourceDoesNotExist)
}
