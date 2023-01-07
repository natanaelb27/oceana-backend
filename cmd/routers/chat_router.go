package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func chat(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/chats"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "111") }, h.ChatGetAllChatList)
		privateRoute.GET("/:to-user-id", func(c *gin.Context) { middlewares.Role(c, "111") }, h.ChatGetChat)
		privateRoute.POST("", func(c *gin.Context) { middlewares.Role(c, "111") }, h.ChatCreate)
	}
}
