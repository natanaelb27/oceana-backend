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
		privateRoute.POST("/password/check-verification-code", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProfileCheckPasswordVerificationCode)
		privateRoute.POST("/password/request-verification-code", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProfileRequestChangePasswordVerificationCode)
		privateRoute.PATCH("/password", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UpdateProfilePassword)
		privateRoute.POST("/sealabs-pay", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayCreate)
		privateRoute.GET("/transaction-history", func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderListGetAll)
		privateRoute.GET("/transaction-history/:orderId", func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderListDetailsGetByID)
	}
}
