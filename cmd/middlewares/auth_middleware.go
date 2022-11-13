package middlewares

import (
	"net/http"
	"strings"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	clientToken := c.GetHeader("Authorization")
	if clientToken == "" {
		utils.ResponseError(c, e.ErrUnauthorized)
		c.Abort()
		return
	}
	extractedToken := strings.Split(clientToken, "Bearer ")
	if len(extractedToken) == 2 {
		clientToken = strings.TrimSpace(extractedToken[1])
	} else {
		utils.ResponseError(c, e.ErrInvalidToken)
		c.Abort()
		return
	}

	jwtWrapper, err := utils.NewJwtWrapper()
	if err != nil {
		return
	}

	claims, err := jwtWrapper.ValidateToken(clientToken)
	if err != nil {
		utils.Response(c, http.StatusUnauthorized, nil, nil, http.StatusText(http.StatusUnauthorized), err.Error())
		c.Abort()
		return
	}

	role, err := utils.GetUserRoleById(claims.User.ID)
	if err != nil {
		utils.ResponseError(c, e.ErrUserNotFound)
		c.Abort()
		return
	}

	if !utils.RolePermissionToResource(role, c.Request.Method, c.Request.RequestURI) {
		utils.ResponseError(c, e.ErrUnauthorized)
		c.Abort()
		return
	}

	c.Set("id", claims.User.ID)

	c.Next()
}
