package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func image(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/images"
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/upload", func(c *gin.Context) { middlewares.Role(c, "011") }, h.UploadImage)
		privateRoute.POST("/delete", func(c *gin.Context) { middlewares.Role(c, "011")}, h.DeleteImage)
	}
}
