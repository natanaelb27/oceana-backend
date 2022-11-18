package middlewares

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func Role(c *gin.Context, allowedRoles string) {
	if !utils.RoleAllowed(c, allowedRoles) {
		utils.ResponseError(c, e.ErrForbidden)
		c.Abort()
		return
	}

	c.Next()
}
