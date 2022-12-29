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
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UserGetProfile)
		privateRoute.PUT("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UserUpdateProfile)
		privateRoute.POST("/logout", func(c *gin.Context) { middlewares.Role(c, "101") }, h.UserLogout)
		privateRoute.PATCH("/email", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.UserUpdateEmail)
		privateRoute.PATCH("/username", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UserUpdateUsername)
		privateRoute.POST("/register-shop", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.UserRegisterMerchant)
		privateRoute.PATCH("/wallet/activate", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletActivate)
		privateRoute.PATCH("/wallet/pin", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletChangePin)
		privateRoute.GET("/wallet", func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletGetDetails)
		privateRoute.POST("/password/check-verification-code", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProfileCheckPasswordVerificationCode)
		privateRoute.POST("/password/request-verification-code", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProfileRequestChangePasswordVerificationCode)
		privateRoute.PATCH("/password", func(c *gin.Context) { middlewares.Role(c, "001") }, h.UpdateProfilePassword)
		privateRoute.GET("/transaction-history", func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderListGetAll)
		privateRoute.GET("/transaction-history/:orderId", func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderListDetailsGetByID)
		privateRoute.PATCH("/orders/:order-id/complete-order", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.OrderCompleteOrder)
		privateRoute.GET("/pending-payment", func(c *gin.Context) { middlewares.Role(c, "001") }, h.PendingPaymentGetAllByUser)
		privateRoute.POST("/refunds/:order-id", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.CreateRefund)
		privateRoute.GET("/refunds", func(c *gin.Context) { middlewares.Role(c, "001") }, h.GetAllRefundDetailsByBuyer)
		privateRoute.GET("/favorites", func(c *gin.Context) { middlewares.Role(c, "001") }, h.GetAllUserFavoriteProducts)
		privateRoute.DELETE("/favorites/batch-delete", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.DeleteFavoriteProducts)

	}
}
