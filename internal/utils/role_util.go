package utils

import (
	"strings"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

func GetUserRoleById(id string) (role string, err error) {
	Db := db.Get()
	user := &models.User{}
	result := Db.First(&user, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	role = user.Role

	return
}

func RolePermissionToResource(role string, requestMethod string, requestUri string) bool {
	matchedUri := ""
	for uri := range constants.ROLE_RESOURCES_RULE {
		if strings.HasPrefix(requestUri, uri) {
			matchedUri = uri
		}
	}

	for k, v := range []rune(role) {
		if string(v) == "1" {
			resourceRoleArray := []rune(constants.ROLE_RESOURCES_RULE[matchedUri][requestMethod])
			if k >= 0 && k < len(resourceRoleArray) && string(resourceRoleArray[k]) == "1" {
				return true
			}
		}
	}

	return false
}
