package services

import (
	"time"

	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type MerchantVoucherService interface {
	Create(id uuid.UUID, request *dtos.MerchantVoucherCreateRequest) (*dtos.MerchantVoucherCreateResponse, error)
}

type merchantVoucherService struct {
	merchantVoucherRepository r.MerchantVoucherRepository
	merchantRepository        r.MerchantRepository
}

type MerchantVoucherServiceConfig struct {
	MerchantVoucherRepository r.MerchantVoucherRepository
	MerchantRepository        r.MerchantRepository
}

func NewMerchantVoucherService(c *MerchantVoucherServiceConfig) MerchantVoucherService {
	return &merchantVoucherService{
		merchantVoucherRepository: c.MerchantVoucherRepository,
		merchantRepository:        c.MerchantRepository,
	}
}

func (s *merchantVoucherService) Create(id uuid.UUID, request *dtos.MerchantVoucherCreateRequest) (*dtos.MerchantVoucherCreateResponse, error) {
	if len(request.Code[4:]) > 5 {
		return nil, e.ErrVoucherCodeSuffixMoreThanFive
	}

	if request.StartDatetime.Before(time.Now()) {
		return nil, e.ErrVoucherStartDatetimeEarlierThanNow
	}

	if request.EndDatetime.Before(request.StartDatetime) {
		return nil, e.ErrVoucherEndDatetimeEarlierThanVoucherStartDatetime
	}

	if int(time.Until(request.StartDatetime).Hours()/24) <= 0 {
		return nil, e.ErrVoucherStartDatetimeMinimumOneDayFromNow
	}

	if int(request.EndDatetime.Sub(request.StartDatetime).Hours()) <= 0 {
		return nil, e.ErrVoucherPeriodMinimumOneHour
	}

	if request.DiscountType == "nominal" {
		if request.Discount > request.MinimumSpendAmount {
			return nil, e.ErrDiscountExceedMinimumSpendAmount
		}

		request.MaximumDiscountAmount = request.Discount
	}

	if request.DiscountType == "percentage" {
		if request.Discount > 100 {
			return nil, e.ErrExceedDiscountPercentage
		}

		if request.MaximumDiscountAmount > request.MinimumSpendAmount {
			return nil, e.ErrMaximumDiscountAmountExceedMinimumSpendAmount
		}
	}

	merchant, err := s.merchantRepository.GetByUserID(id)

	if err != nil {
		return nil, err
	}

	model := request.ToModel(merchant.ID)

	newModel, err := s.merchantVoucherRepository.Create(model)

	if err != nil {
		return nil, err
	}

	return dtos.NewMerchantVoucherCreateResponse(newModel), nil
}
