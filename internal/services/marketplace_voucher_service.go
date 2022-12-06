package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
)

type MarketplaceVoucherService interface {
	Create(marketplaceVoucher *models.MarketplaceVoucher) (*models.MarketplaceVoucher, error)
	GetAll() ([]*models.MarketplaceVoucher, error)
	GetAllByActive() ([]*models.MarketplaceVoucher, error)
	GetById(id uint) (*models.MarketplaceVoucher, error)
	Update(id uint, newModel *models.MarketplaceVoucher) (*models.MarketplaceVoucher, error)
	Delete(id uint) (*models.MarketplaceVoucher, error)
}

type marketplaceVoucherService struct {
	marketplaceVoucherRepository r.MarketplaceVoucherRepository
}

type MarketplaceVoucherServiceConfig struct {
	MarketplaceVoucherRepository r.MarketplaceVoucherRepository
}

func NewMarketplaceVoucherService(c *MarketplaceVoucherServiceConfig) MarketplaceVoucherService {
	return &marketplaceVoucherService{
		marketplaceVoucherRepository: c.MarketplaceVoucherRepository,
	}
}

func (p *marketplaceVoucherService) Create(model *models.MarketplaceVoucher) (*models.MarketplaceVoucher, error) {
	return p.marketplaceVoucherRepository.Create(model)
}

func (p *marketplaceVoucherService) GetAll() ([]*models.MarketplaceVoucher, error) {
	return p.marketplaceVoucherRepository.GetAll()
}

func (p *marketplaceVoucherService) GetAllByActive() ([]*models.MarketplaceVoucher, error) {
	return p.marketplaceVoucherRepository.GetAllByActive()
}

func (p *marketplaceVoucherService) GetById(id uint) (*models.MarketplaceVoucher, error) {
	return p.marketplaceVoucherRepository.GetById(id)
}

func (p *marketplaceVoucherService) Update(id uint, newModel *models.MarketplaceVoucher) (*models.MarketplaceVoucher, error) {
	err := p.marketplaceVoucherRepository.Update(id, newModel)
	if err != nil {
		return nil, err
	}

	return p.marketplaceVoucherRepository.GetById(id)
}

func (p *marketplaceVoucherService) Delete(id uint) (*models.MarketplaceVoucher, error) {
	return p.marketplaceVoucherRepository.Delete(id)
}
