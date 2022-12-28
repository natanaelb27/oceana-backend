package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func role(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/roles"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("", h.RoleGetAll)
		publicRoute.POST("", h.RoleCreate)
		publicRoute.DELETE("/:id", h.RoleDelete)
		publicRoute.PUT("/:id", h.RoleUpdate)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/:id", func(c *gin.Context) { middlewares.Role(c, "001") }, h.RoleGetById)
		privateRoute.GET("/private-example", func(c *gin.Context) { middlewares.Role(c, "001") }, h.RoleGetAll)
		privateRoute.POST("/failing",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "011") },
			h.RoleCreateFailingWithTransaction,
		)
	}
	// We can also define use the middleware on its group like below
	// anotherPrivateRoute := private.Group(path, func(c *gin.Context) { middlewares.Role(c, "001") })
	// {
	// 	anotherPrivateRoute.GET("/another/:id", h.RoleGetById)
	// 	anotherPrivateRoute.GET("/another/private-example", h.RoleGetAll)
	// }

}
