package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type CategoryService interface {
	GetAll(query *externalDtos.CategoryQuery) ([]*dtos.Category, error)
	GetAncestorsByID(id uint) ([]*externalDtos.GetCategoriesResponse, error)
	GetMerchantCategories(id uint) ([]*externalDtos.GetCategoriesResponse, error)
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

func (p *categoryService) GetAllSubCategoriesRecursive(initialLevel int, level int, levelLimit int, categoryData []*models.Category) ([]*dtos.Category, error) {
	category := []*dtos.Category{}

	if levelLimit != -1 && level >= levelLimit {
		return category, nil
	}

	for i, cat := range categoryData {
		category = append(category, &dtos.Category{
			ID:    cat.ID,
			Image: cat.MediaURL,
			Name:  cat.Name,
			Level: level,
		})

		subCategoryData, err := p.categoryRepository.GetAllSubCategoriesByParentID(cat.ID)
		if err != nil {
			return nil, err
		}
		if len(subCategoryData) > 0 {
			subCategory, err := p.GetAllSubCategoriesRecursive(initialLevel, level+1, levelLimit, subCategoryData)
			if err != nil {
				return nil, err
			}

			category[i].Category = subCategory
		}
	}

	return category, nil
}

func (p *categoryService) GetCategoryLevel(id uint, level uint) (uint, error) {
	var err error

	category, err := p.categoryRepository.GetByID(id)
	if err != nil && err != e.ErrDataNotFound {
		return 0, err
	}
	if err == e.ErrDataNotFound {
		return level, nil
	}

	if category.ParentID != 0 {
		level, err = p.GetCategoryLevel(category.ParentID, level+1)
		if err != nil && err != e.ErrDataNotFound {
			return 0, err
		}
		if err == e.ErrDataNotFound {
			return level, nil
		}
	}

	return level, err
}

func (p *categoryService) GetAll(query *externalDtos.CategoryQuery) ([]*dtos.Category, error) {
	var err error
	var categoryData []*models.Category
	var currentLevel uint

	currentLevel = 1

	if query.CategoryID != 0 {
		categoryData, err = p.categoryRepository.GetAllByID(query.CategoryID)
		if err != nil {
			return nil, err
		}

		currentLevel, err = p.GetCategoryLevel(query.CategoryID, 1)
		if err != nil {
			return nil, err
		}
	} else {
		categoryData, err = p.categoryRepository.GetAllCategories()
		if err != nil {
			return nil, err
		}
	}

	currentLevelInt := int(currentLevel)
	levelDepth := query.LevelDepth
	if query.LevelDepth != -1 {
		levelDepth += currentLevelInt
	}

	category, err := p.GetAllSubCategoriesRecursive(currentLevelInt, currentLevelInt, levelDepth, categoryData)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (p *categoryService) GetAncestorsByID(id uint) ([]*externalDtos.GetCategoriesResponse, error) {
	category, err := p.categoryRepository.GetAncestorsByID(id)

	if err != nil {
		return nil, err
	}

	return externalDtos.NewGetCategoriesResponse(category), nil
}

func (p *categoryService) GetMerchantCategories(id uint) ([]*externalDtos.GetCategoriesResponse, error) {
	category, err := p.categoryRepository.GetMerchantCategoriesByMerchantID(id)

	if err != nil {
		return nil, err
	}

	return externalDtos.NewGetCategoriesResponse(category), nil
}
