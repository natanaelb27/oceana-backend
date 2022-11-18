package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func profile(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/profile"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", h.UserGetProfile)
		privateRoute.PUT("", h.UserUpdateProfile)
		privateRoute.POST("/register-shop", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UserRegisterMerchant)
		privateRoute.PATCH("/wallet/activate", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletActivate)
		privateRoute.PATCH("/wallet/pin", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletChangePin)
		privateRoute.POST("/sealabs-pay", h.SealabsPayCreate)
	}
}
