package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
)

type ProductItemService interface {
	CreateProductItem(productItem *models.ProductItem) (*models.ProductItem, error)
}

type productItemService struct {
	productItemRepository r.ProductItemRepository
}

type ProductItemServiceConfig struct {
	ProductItemRepository r.ProductItemRepository
}

func NewProductItemService(c *ProductItemServiceConfig) ProductItemService {
	return &productItemService{
		productItemRepository: c.ProductItemRepository,
	}
}

func (pi *productItemService) CreateProductItem(model *models.ProductItem) (*models.ProductItem, error) {
	return pi.productItemRepository.CreateProductItem(model)
}
