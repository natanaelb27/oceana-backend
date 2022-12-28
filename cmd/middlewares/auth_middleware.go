package middlewares

import (
	"net/http"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	clientToken, err := c.Cookie("access_token")
	if err != nil {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			utils.ResponseError(c, e.ErrUnauthorized)
			c.Abort()
			return
		}

		jwtWrapper, err := utils.NewJwtWrapper()
		if err != nil {
			return
		}

		claims, err := jwtWrapper.ValidateToken(refreshToken, nil)
		if err != nil {
			utils.Response(c, http.StatusUnauthorized, nil, nil, http.StatusText(http.StatusUnauthorized), err.Error())
			c.Abort()
			return
		}
		c.Set("id", claims.User.ID)
		c.Next()
		return
	}

	jwtWrapper, err := utils.NewJwtWrapper()
	if err != nil {
		return
	}

	claims, err := jwtWrapper.ValidateToken(clientToken, nil)
	if err != nil {
		utils.Response(c, http.StatusUnauthorized, nil, nil, http.StatusText(http.StatusUnauthorized), err.Error())
		c.Abort()
		return
	}

	role, err := utils.GetUserRoleById(claims.User.ID.String())
	if err != nil {
		utils.ResponseError(c, e.ErrUserNotFound)
		c.Abort()
		return
	}

	c.Set("id", claims.User.ID)
	c.Set("role", role)

	c.Next()
}
