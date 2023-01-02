package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
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
