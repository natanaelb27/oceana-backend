package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func media(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/medias"
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/upload", func(c *gin.Context) { middlewares.Role(c, "111") }, h.UploadMedia)
		privateRoute.POST("/delete", func(c *gin.Context) { middlewares.Role(c, "111") }, h.DeleteMedia)
	}
}
