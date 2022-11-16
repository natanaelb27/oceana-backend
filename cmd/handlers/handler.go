package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	roleService     services.RoleService
	walletService   services.WalletService
	authService     services.AuthService
	categoryService services.CategoryService
	userService     services.UserService
}

type HandlerConfig struct {
	RoleService     services.RoleService
	WalletService   services.WalletService
	AuthService     services.AuthService
	CategoryService services.CategoryService
	UserService     services.UserService
}

func New(c *HandlerConfig) *Handler {
	return &Handler{
		roleService:     c.RoleService,
		walletService:   c.WalletService,
		authService:     c.AuthService,
		categoryService: c.CategoryService,
		userService:     c.UserService,
	}
}

func (h *Handler) Root(c *gin.Context) {
	utils.ResponseSuccess(c, "Oceana API", nil)
}

func (h *Handler) NotFound(c *gin.Context) {
	utils.ResponseError(c, e.ErrResourceDoesNotExist)
}
