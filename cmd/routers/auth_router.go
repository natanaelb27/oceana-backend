package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func auth(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/auth"
	publicRoute := public.Group(path)
	{
		publicRoute.POST("/register", middlewares.DBTransactionMiddleware(), h.UserCreate)
		publicRoute.POST("/login", h.UserLogin)
		publicRoute.POST("/register-google", middlewares.DBTransactionMiddleware(), h.UserRegisterWithGoogle)
		publicRoute.POST("/login-google", h.UserLoginWithGoogle)
		publicRoute.POST("/reset-password/verify", h.UserVerifyResetPassword)
		publicRoute.POST("/reset-password/:reset-password-id", h.UserResetPassword)
		publicRoute.POST("/refresh", h.RefreshAccessToken)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/step-up/password", func(c *gin.Context) { middlewares.Role(c, "001") }, h.StepUpPassword)
	}
}
