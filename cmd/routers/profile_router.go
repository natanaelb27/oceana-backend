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
		privateRoute.PATCH("/wallet/activate", h.WalletActivate)
		privateRoute.PATCH("/wallet/pin", h.WalletChangePin)
		privateRoute.POST("/password/check-verification-code", h.ProfileCheckPasswordVerificationCode)
		privateRoute.POST("/password/request-verification-code", h.ProfileRequestChangePasswordVerificationCode)
		privateRoute.PATCH("/password", h.UpdateProfilePassword)
	}
}
