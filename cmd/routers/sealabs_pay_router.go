package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
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
		privateRoute.POST("/register", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayCreate)
		privateRoute.POST("/register-otp", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayRegisterOTP)
		privateRoute.POST("/topup-otp", func(c *gin.Context) { middlewares.Role(c, "011") }, h.SealabsPayTopupOTP)
	}
}
