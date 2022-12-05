package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	roleService              services.RoleService
	walletService            services.WalletService
	authService              services.AuthService
	categoryService          services.CategoryService
	userService              services.UserService
	profileService           services.ProfileService
	merchantService          services.MerchantService
	sealabsPayService        services.SealabsPayService
	merchantVoucherService   services.MerchantVoucherService
	productService           services.ProductService
	productItemService       services.ProductItemService
	addressService           services.AddressService
	orderService             services.OrderService
	productReviewService     services.ProductReviewService
	merchantPromotionService services.MerchantPromotionService
	cartService              services.CartService
	subDistrictService       services.SubDistrictService
}

type HandlerConfig struct {
	RoleService              services.RoleService
	WalletService            services.WalletService
	AuthService              services.AuthService
	CategoryService          services.CategoryService
	ProfileService           services.ProfileService
	UserService              services.UserService
	MerchantService          services.MerchantService
	SealabsPayService        services.SealabsPayService
	MerchantVoucherService   services.MerchantVoucherService
	ProductService           services.ProductService
	ProductItemService       services.ProductItemService
	AddressService           services.AddressService
	OrderService             services.OrderService
	ProductReviewService     services.ProductReviewService
	MerchantPromotionService services.MerchantPromotionService
	CartService              services.CartService
	SubDistrictService       services.SubDistrictService
}

func New(c *HandlerConfig) *Handler {
	return &Handler{
		roleService:              c.RoleService,
		walletService:            c.WalletService,
		authService:              c.AuthService,
		categoryService:          c.CategoryService,
		userService:              c.UserService,
		profileService:           c.ProfileService,
		merchantService:          c.MerchantService,
		sealabsPayService:        c.SealabsPayService,
		merchantVoucherService:   c.MerchantVoucherService,
		productService:           c.ProductService,
		productItemService:       c.ProductItemService,
		addressService:           c.AddressService,
		orderService:             c.OrderService,
		productReviewService:     c.ProductReviewService,
		merchantPromotionService: c.MerchantPromotionService,
		cartService:              c.CartService,
		subDistrictService:       c.SubDistrictService,
	}
}

func (h *Handler) Root(c *gin.Context) {
	utils.ResponseSuccess(c, "Oceana API", nil)
}

func (h *Handler) NotFound(c *gin.Context) {
	utils.ResponseError(c, e.ErrResourceDoesNotExist)
}
