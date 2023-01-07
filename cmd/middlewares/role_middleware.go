package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func Role(c *gin.Context, allowedRoles string) {
	if !utils.RoleAllowed(c, allowedRoles) {
		utils.ResponseError(c, e.ErrForbidden)
		c.Abort()
		return
	}

	c.Next()
}
