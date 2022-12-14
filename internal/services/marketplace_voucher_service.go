package services

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type MarketplaceVoucherService interface {
	Create(marketplaceVoucher *models.MarketplaceVoucher) (*models.MarketplaceVoucher, error)
	GetAll() ([]*models.MarketplaceVoucher, error)
	GetAllByActive() ([]*models.MarketplaceVoucher, error)
	GetById(id uint) (*models.MarketplaceVoucher, error)
	Update(id uint, newModel *models.MarketplaceVoucher) (*models.MarketplaceVoucher, error)
	Delete(id uint) (*models.MarketplaceVoucher, error)
	GetAllByBuyer(id uuid.UUID, query *dtos.BuyerMarketplaceVoucherQuery) (*dtos.MarketplaceVoucherGetAllResponse, *utils.Pagination, error)
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
	if model.StartDatetime.Before(time.Now()) {
		return nil, e.ErrVoucherStartDatetimeEarlierThanNow
	}
	if model.EndDatetime.Before(model.StartDatetime) {
		return nil, e.ErrVoucherEndDatetimeEarlierThanVoucherStartDatetime
	}

	if int(time.Until(model.StartDatetime).Hours()/24) <= 0 {
		return nil, e.ErrVoucherStartDatetimeMinimumOneDayFromNow
	}

	if int(model.EndDatetime.Sub(model.StartDatetime).Hours()) <= 0 {
		return nil, e.ErrVoucherPeriodMinimumOneHour
	}

	if int(model.EndDatetime.Sub(model.StartDatetime).Hours()/24) > 90 {
		return nil, e.ErrVoucherPeriodExceedMaximumDays
	}
	if model.DiscountType == "nominal" {
		if model.Discount > model.MinimumSpendAmount {
			return nil, e.ErrDiscountExceedMinimumSpendAmount
		}

		model.MaximumDiscountAmount = model.Discount
	}

	if model.DiscountType == "percentage" {
		if model.Discount > 100 {
			return nil, e.ErrExceedDiscountPercentage
		}

		if model.MaximumDiscountAmount > model.MinimumSpendAmount {
			return nil, e.ErrMaximumDiscountAmountExceedMinimumSpendAmount
		}
	}

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
	marketplaceVoucherDetails, err := p.marketplaceVoucherRepository.GetDetailsByID(id)

	if err != nil {
		return nil, err
	}

	if marketplaceVoucherDetails.VoucherStatus == constants.EndedPromotionOrVoucher {
		return nil, e.ErrVoucherHasEnded
	}

	if marketplaceVoucherDetails.VoucherStatus == constants.OngoingPromotionOrVoucher {
		if int(marketplaceVoucherDetails.StartDatetime.Sub(newModel.StartDatetime).Minutes()) != 0 {
			return nil, e.ErrOngoingStartDatetimeUpdated
		}

		if newModel.MaximumQuota < marketplaceVoucherDetails.CurrentQuota {
			return nil, e.ErrMaximumQuotaLowerThanCurrentQuota
		}

		if newModel.EndDatetime.Before(time.Now()) {
			return nil, e.ErrVoucherEndDatetimeEarlierThanNow
		}

		if newModel.EndDatetime.Sub(marketplaceVoucherDetails.EndDatetime).Minutes() > 0 {
			return nil, e.ErrUpdatedVoucherEndDatetimeExceedCurrentEndDatetime
		}
	}

	if marketplaceVoucherDetails.VoucherStatus == "upcoming" {
		if newModel.StartDatetime.Before(time.Now()) {
			return nil, e.ErrVoucherStartDatetimeEarlierThanNow
		}
	}

	if newModel.EndDatetime.Before(newModel.StartDatetime) {
		return nil, e.ErrVoucherEndDatetimeEarlierThanVoucherStartDatetime
	}

	if int(newModel.EndDatetime.Sub(newModel.StartDatetime).Hours()) <= 0 {
		return nil, e.ErrVoucherPeriodMinimumOneHour
	}

	if int(newModel.EndDatetime.Sub(newModel.StartDatetime).Hours()/24) > 90 {
		return nil, e.ErrVoucherPeriodExceedMaximumDays
	}

	if newModel.DiscountType == "nominal" {
		if newModel.Discount > newModel.MinimumSpendAmount {
			return nil, e.ErrDiscountExceedMinimumSpendAmount
		}

		newModel.MaximumDiscountAmount = newModel.Discount
	}

	if newModel.DiscountType == "percentage" {
		if newModel.Discount > 100 {
			return nil, e.ErrExceedDiscountPercentage
		}

		if newModel.MaximumDiscountAmount > newModel.MinimumSpendAmount {
			return nil, e.ErrMaximumDiscountAmountExceedMinimumSpendAmount
		}
	}

	err = p.marketplaceVoucherRepository.Update(id, newModel)
	if err != nil {
		return nil, err
	}

	return p.marketplaceVoucherRepository.GetById(id)
}

func (p *marketplaceVoucherService) Delete(id uint) (*models.MarketplaceVoucher, error) {
	marketplaceVoucherDetails, err := p.marketplaceVoucherRepository.GetDetailsByID(id)

	if err != nil {
		return nil, err
	}

	if marketplaceVoucherDetails.VoucherStatus != "upcoming" {
		return nil, e.ErrVoucherNotUpcoming
	}

	return p.marketplaceVoucherRepository.Delete(id)
}

func (p *marketplaceVoucherService) GetAllByBuyer(id uuid.UUID, query *dtos.BuyerMarketplaceVoucherQuery) (*dtos.MarketplaceVoucherGetAllResponse, *utils.Pagination, error) {
	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	marketplaceVoucherDetails, pagination, err := p.marketplaceVoucherRepository.GetAllDetailsByBuyer(id, query, pagination)

	if err != nil {
		return nil, nil, err
	}

	return dtos.NewMarketplaceVoucherGetAllResponse(marketplaceVoucherDetails), pagination, nil
}
