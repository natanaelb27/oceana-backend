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
		privateRoute.POST("/logout", h.UserLogout)
		privateRoute.PATCH("/email", h.UserUpdateEmail)
		privateRoute.PATCH("/username", h.UserUpdateUsername)
		privateRoute.POST("/register-shop", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UserRegisterMerchant)
		privateRoute.PATCH("/wallet/activate", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletActivate)
		privateRoute.PATCH("/wallet/pin", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletChangePin)
		privateRoute.GET("/wallet", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletGetDetails)
		privateRoute.POST("/password/check-verification-code", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProfileCheckPasswordVerificationCode)
		privateRoute.POST("/password/request-verification-code", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProfileRequestChangePasswordVerificationCode)
		privateRoute.PATCH("/password", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UpdateProfilePassword)
		privateRoute.GET("/transaction-history", func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderListGetAll)
		privateRoute.GET("/transaction-history/:orderId", func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderListDetailsGetByID)
		privateRoute.PATCH("/orders/:order-id/complete-order", func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderCompleteOrder)
		privateRoute.GET("/pending-payment", func(c *gin.Context) { middlewares.Role(c, "001") }, h.PendingPaymentGetAllByUser)
		privateRoute.POST("/refund/:order-id", func(c *gin.Context) { middlewares.Role(c, "001") }, h.CreateRefund)
		privateRoute.GET("/refund", func(c *gin.Context) { middlewares.Role(c, "001") }, h.GetAllRefundDetailsByBuyer)
	}
}
