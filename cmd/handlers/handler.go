package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
<<<<<<< HEAD
	roleService    services.RoleService
	authService    services.AuthService
	addressService services.AddressService
}

type HandlerConfig struct {
	RoleService    services.RoleService
	AuthService    services.AuthService
	AddressService services.AddressService
=======
	roleService     services.RoleService
	walletService   services.WalletService
	authService     services.AuthService
	categoryService services.CategoryService
}

type HandlerConfig struct {
	RoleService     services.RoleService
	WalletService   services.WalletService
	AuthService     services.AuthService
	CategoryService services.CategoryService
>>>>>>> a817cfb680e773bccf30b513f4bfde6aa7066b5e
}

func New(c *HandlerConfig) *Handler {
	return &Handler{
<<<<<<< HEAD
		roleService:    c.RoleService,
		authService:    c.AuthService,
		addressService: c.AddressService,
=======
		roleService:     c.RoleService,
		walletService:   c.WalletService,
		authService:     c.AuthService,
		categoryService: c.CategoryService,
>>>>>>> a817cfb680e773bccf30b513f4bfde6aa7066b5e
	}
}

func (h *Handler) Root(c *gin.Context) {
	utils.ResponseSuccess(c, "Oceana API", nil)
}

func (h *Handler) NotFound(c *gin.Context) {
	utils.ResponseError(c, e.ErrResourceDoesNotExist)
}
