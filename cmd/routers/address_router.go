package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func address(private *gin.RouterGroup, h *handlers.Handler) {
	path := "/profile/addresses"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/", h.AddressGetAll)
		privateRoute.GET("/:address_id", h.AddressGetById)
		privateRoute.POST("/", h.AddressCreate)
	}
}
