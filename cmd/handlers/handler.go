package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	roleService     services.RoleService
	authService     services.AuthService
	addressService  services.AddressService
	walletService   services.WalletService
	categoryService services.CategoryService
}

type HandlerConfig struct {
	RoleService     services.RoleService
	AuthService     services.AuthService
	AddressService  services.AddressService
	WalletService   services.WalletService
	CategoryService services.CategoryService
}

func New(c *HandlerConfig) *Handler {
	return &Handler{
		roleService:     c.RoleService,
		walletService:   c.WalletService,
		authService:     c.AuthService,
		categoryService: c.CategoryService,
	}
}

func (h *Handler) Root(c *gin.Context) {
	utils.ResponseSuccess(c, "Oceana API", nil)
}

func (h *Handler) NotFound(c *gin.Context) {
	utils.ResponseError(c, e.ErrResourceDoesNotExist)
}
