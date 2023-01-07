package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func sealabsPay(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	callbackPath := "/callback/sealabs-pay"
	callbackRoute := public.Group(callbackPath)
	{
		callbackRoute.POST("", h.SealabsPayCallback)
	}

	path := "/sealabs-pay"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayGetAll)
		privateRoute.POST("", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayCreate)
		privateRoute.GET("/:sealabs-pay-id", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayGetById)
		privateRoute.DELETE("/:sealabs-pay-id", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayDelete)
		privateRoute.PATCH("/:sealabs-pay-id", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPaySetDefault)
		privateRoute.POST("/register-otp", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayRegisterOTP)
		privateRoute.POST("/topup-otp", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayTopupOTP)
		privateRoute.POST("/payment-otp", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayPaymentOTP)
	}
}
