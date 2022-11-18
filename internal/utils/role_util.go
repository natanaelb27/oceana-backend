package utils

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func GetUserRoleById(id string) (role string, err error) {
	Db := db.Get()
	user := &models.User{}
	result := Db.Where("id = ?", id).First(&user)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	role = user.Role

	return
}

func RoleAllowed(c *gin.Context, allowedRole string) bool {
	for k, v := range []rune(c.GetString("role")) {
		if string(v) == "1" {
			resourceRoleArray := []rune(allowedRole)
			if k >= 0 && k < len(resourceRoleArray) && string(resourceRoleArray[k]) == "1" {
				return true
			}
		}
	}

	return false
}
