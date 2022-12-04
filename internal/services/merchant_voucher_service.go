package services

import (
	"time"

	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type MerchantVoucherService interface {
	Create(id uuid.UUID, request *dtos.MerchantVoucherCreateRequest) (*dtos.MerchantVoucherCreateResponse, error)
	Update(userID uuid.UUID, voucherID uint, request *dtos.MerchantVoucherUpdateRequest) (*dtos.MerchantVoucherUpdateResponse, error)
	Delete(userID uuid.UUID, voucherID uint) (*dtos.MerchantVoucherDeleteResponse, error)
	GetAll(query *dtos.MerchantVoucherQuery) (*dtos.MerchantVoucherGetAllResponse, *utils.Pagination, error)
	GetByID(id uint) (*dtos.MerchantVoucherGetByIDResponse, error)
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
	if !utils.ValidateMerchantVoucherCode(request.Code[4:]) {
		return nil, e.ErrInvalidMerchantVoucherCode
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

	if int(request.EndDatetime.Sub(request.StartDatetime).Hours()/24) > 90 {
		return nil, e.ErrVoucherPeriodExceedMaximumDays
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

func (s *merchantVoucherService) Update(userID uuid.UUID, voucherID uint, request *dtos.MerchantVoucherUpdateRequest) (*dtos.MerchantVoucherUpdateResponse, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, err
	}

	merchantVoucherDetails, err := s.merchantVoucherRepository.GetDetailsByID(voucherID)

	if err != nil {
		return nil, err
	}

	if merchantVoucherDetails.MerchantID != merchant.ID {
		return nil, e.ErrForbiddenMerchant
	}

	if merchantVoucherDetails.VoucherStatus == "ended" {
		return nil, e.ErrVoucherHasEnded
	}

	if merchantVoucherDetails.VoucherStatus == "ongoing" {
		if int(merchantVoucherDetails.StartDatetime.Sub(request.StartDatetime).Minutes()) != 0 {
			return nil, e.ErrOngoingStartDatetimeUpdated
		}

		if request.MaximumQuota < merchantVoucherDetails.CurrentQuota {
			return nil, e.ErrMaximumQuotaLowerThanCurrentQuota
		}

		if request.EndDatetime.Before(time.Now()) {
			return nil, e.ErrVoucherEndDatetimeEarlierThanNow
		}

		if request.EndDatetime.Sub(merchantVoucherDetails.EndDatetime).Minutes() > 0 {
			return nil, e.ErrUpdatedVoucherEndDatetimeExceedCurrentEndDatetime
		}
	}

	if merchantVoucherDetails.VoucherStatus == "upcoming" {
		if request.StartDatetime.Before(time.Now()) {
			return nil, e.ErrVoucherStartDatetimeEarlierThanNow
		}
	}

	if request.EndDatetime.Before(request.StartDatetime) {
		return nil, e.ErrVoucherEndDatetimeEarlierThanVoucherStartDatetime
	}

	if int(request.EndDatetime.Sub(request.StartDatetime).Hours()) <= 0 {
		return nil, e.ErrVoucherPeriodMinimumOneHour
	}

	if int(request.EndDatetime.Sub(request.StartDatetime).Hours()/24) > 90 {
		return nil, e.ErrVoucherPeriodExceedMaximumDays
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

	model := request.ToModel(voucherID)

	newModel, err := s.merchantVoucherRepository.Update(model)

	if err != nil {
		return nil, err
	}

	return dtos.NewMerchantVoucherUpdateResponse(newModel), nil
}

func (s *merchantVoucherService) Delete(userID uuid.UUID, voucherID uint) (*dtos.MerchantVoucherDeleteResponse, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, err
	}

	merchantVoucherDetails, err := s.merchantVoucherRepository.GetDetailsByID(voucherID)

	if err != nil {
		return nil, err
	}

	if merchantVoucherDetails.MerchantID != merchant.ID {
		return nil, e.ErrForbiddenMerchant
	}

	if merchantVoucherDetails.VoucherStatus != "upcoming" {
		return nil, e.ErrVoucherNotUpcoming
	}

	newModel, err := s.merchantVoucherRepository.Delete(voucherID)

	if err != nil {
		return nil, err
	}

	return dtos.NewMerchantVoucherDeleteResponse(newModel), nil
}

func (s *merchantVoucherService) GetAll(query *dtos.MerchantVoucherQuery) (*dtos.MerchantVoucherGetAllResponse, *utils.Pagination, error) {
	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	merchantVoucherDetails, pagination, err := s.merchantVoucherRepository.GetAllDetails(query, pagination)

	if err != nil {
		return nil, nil, err
	}

	return dtos.NewMerchantVoucherGetAllResponse(merchantVoucherDetails), pagination, nil
}

func (s *merchantVoucherService) GetByID(id uint) (*dtos.MerchantVoucherGetByIDResponse, error) {
	merchantVoucherDetails, err := s.merchantVoucherRepository.GetDetailsByID(id)

	if err != nil {
		return nil, err
	}

	return dtos.NewMerchantVoucherGetByIDResponse(merchantVoucherDetails), nil
}
