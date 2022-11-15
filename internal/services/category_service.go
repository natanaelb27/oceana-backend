package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
)

type CategoryService interface {
	GetAll() ([]*dtos.Category, error)
}

type categoryService struct {
	categoryRepository r.CategoryRepository
}

type CategoryServiceConfig struct {
	CategoryRepository r.CategoryRepository
}

func NewCategoryService(c *CategoryServiceConfig) CategoryService {
	return &categoryService{
		categoryRepository: c.CategoryRepository,
	}
}

func (p *categoryService) GetAllSubCategoriesRecursive(level int, categoryData []*models.Category) ([]*dtos.Category, error) {
	category := []*dtos.Category{}

	for i, cat := range categoryData {
		category = append(category, &dtos.Category{
			ID:    cat.ID,
			Name:  cat.Name,
			Level: level,
		})

		subCategoryData, err := p.categoryRepository.GetAllSubCategoriesByParentID(cat.ID)

		if err != nil {
			return nil, err
		}

		if len(subCategoryData) > 0 {
			subCategory, err := p.GetAllSubCategoriesRecursive(level+1, subCategoryData)

			if err != nil {
				return nil, err
			}

			category[i].Category = subCategory
		}
	}

	return category, nil
}

func (p *categoryService) GetAll() ([]*dtos.Category, error) {
	categoryData, err := p.categoryRepository.GetAllCategories()

	if err != nil {
		return nil, err
	}

	category, err := p.GetAllSubCategoriesRecursive(1, categoryData)

	if err != nil {
		return nil, err
	}

	return category, nil
}
