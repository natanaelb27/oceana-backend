package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	roleService    services.RoleService
	authService    services.AuthService
	addressService services.AddressService
}

type HandlerConfig struct {
	RoleService    services.RoleService
	AuthService    services.AuthService
	AddressService services.AddressService
}

func New(c *HandlerConfig) *Handler {
	return &Handler{
		roleService:    c.RoleService,
		authService:    c.AuthService,
		addressService: c.AddressService,
	}
}

func (h *Handler) Root(c *gin.Context) {
	utils.ResponseSuccess(c, "Oceana API", nil)
}

func (h *Handler) NotFound(c *gin.Context) {
	utils.ResponseError(c, e.ErrResourceDoesNotExist)
}
