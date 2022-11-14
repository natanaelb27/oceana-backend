package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
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

func (p *categoryService) GetAll() ([]*dtos.Category, error) {
	categoryData, err := p.categoryRepository.GetAllCategories()

	if err != nil {
		return nil, err
	}

	category, err := utils.GetAllSubCategoriesRecursive(1, categoryData)

	if err != nil {
		return nil, err
	}

	return category, nil
}
