package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func profile(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/profile"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", h.UserGetProfile)
		privateRoute.PUT("", h.UserUpdateProfile)
		privateRoute.POST("/register-shop", h.UserRegisterMerchant)
		privateRoute.PATCH("/wallet/activate", h.WalletActivate)
		privateRoute.POST("/sealabs-pay", h.SealabsPayCreate)
	}
}
