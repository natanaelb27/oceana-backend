package utils

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

func GetAllSubCategoriesRecursive(level int, categoryData []*models.Category) ([]*dtos.Category, error) {
	category := []*dtos.Category{}

	for i, cat := range categoryData {
		category = append(category, &dtos.Category{
			ID:    cat.ID,
			Name:  cat.Name,
			Level: level,
		})

		Db := db.Get()

		var subCategoryData []*models.Category

		result := Db.Model(&models.Category{}).Where("parent_id = ?", cat.ID).Find(&subCategoryData)

		if result.Error != nil {
			return nil, e.ErrDatabaseOperationFailed
		}

		if len(subCategoryData) > 0 {
			subCategory, err := GetAllSubCategoriesRecursive(level+1, subCategoryData)

			if err != nil {
				return nil, err
			}

			category[i].Category = subCategory
		}
	}

	return category, nil
}
