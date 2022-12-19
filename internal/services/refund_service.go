package services

import (
	"errors"
	"math"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type RefundService interface {
	Create(userID uuid.UUID, orderID uuid.UUID, request *dtos.CreateRefundRequest) (*dtos.CreateRefundResponse, error)
	GetAllDetailsByBuyer(userID uuid.UUID, query *dtos.GetAllRefundQuery) ([]*dtos.GetRefundDetailsResponse, *utils.Pagination, error)
	GetDetailsByID(userID uuid.UUID, refundID uuid.UUID) (*dtos.GetRefundDetailsResponse, error)
	GetAllDetailsByMerchant(userID uuid.UUID, query *dtos.GetAllRefundQuery) ([]*dtos.GetRefundDetailsResponse, *utils.Pagination, error)
	GetAllDetailsByAdmin(userID uuid.UUID, query *dtos.GetAllRefundQuery) ([]*dtos.GetRefundDetailsResponse, *utils.Pagination, error)
}

type refundService struct {
	merchantRepository     r.MerchantRepository
	orderRepository        r.OrderRepository
	orderDetailRepository  r.OrderDetailRepository
	refundRepository       r.RefundRepository
	refundDetailRepository r.RefundDetailRepository
	refundMediaRepository  r.RefundMediaRepository
	userRepository         r.UserRepository
}

type RefundServiceConfig struct {
	MerchantRepository     r.MerchantRepository
	OrderRepository        r.OrderRepository
	OrderDetailRepository  r.OrderDetailRepository
	RefundRepository       r.RefundRepository
	RefundDetailRepository r.RefundDetailRepository
	RefundMediaRepository  r.RefundMediaRepository
	UserRepository         r.UserRepository
}

func NewRefundService(c *RefundServiceConfig) RefundService {
	return &refundService{
		merchantRepository:     c.MerchantRepository,
		orderRepository:        c.OrderRepository,
		orderDetailRepository:  c.OrderDetailRepository,
		refundRepository:       c.RefundRepository,
		refundDetailRepository: c.RefundDetailRepository,
		refundMediaRepository:  c.RefundMediaRepository,
		userRepository:         c.UserRepository,
	}
}

func (s *refundService) Create(userID uuid.UUID, orderID uuid.UUID, request *dtos.CreateRefundRequest) (*dtos.CreateRefundResponse, error) {
	order, err := s.orderRepository.OrderListDetailsGetByID(orderID)

	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, e.ErrForbiddenOrder
	}

	if order.OrderStatus != constants.DeliveredStatus {
		return nil, e.ErrStatusIsNotDelivered
	}

	refund, err := s.refundRepository.Create(&models.Refund{
		OrderID:               orderID,
		Reason:                request.Reason,
		RefundRequestDatetime: time.Now(),
		Status:                constants.RefundStatusOngoing,
	})

	if err != nil {
		return nil, err
	}

	var refundDetailList []*models.RefundDetail

	for _, rd := range request.OrderDetail {
		orderDetail, err := s.orderDetailRepository.GetByID(rd.OrderDetailID)

		if err != nil {
			return nil, err
		}

		if orderDetail.OrderID != orderID {
			return nil, e.ErrOrderDoesNotHaveThatOrderDetail
		}

		refundDetail, err := s.refundDetailRepository.CalculateRefundDetails(rd.OrderDetailID, rd.Quantity)

		if err != nil {
			return nil, err
		}

		refundDetailList = append(refundDetailList, &models.RefundDetail{
			RefundID:           refund.ID,
			OrderDetailID:      rd.OrderDetailID,
			Quantity:           rd.Quantity,
			AmountUserGet:      math.Round(refundDetail.AmountUserGet),
			AmountMerchantLost: math.Round(refundDetail.AmountMerchantLost),
		})
	}

	_, err = s.refundDetailRepository.Create(refundDetailList)

	if err != nil {
		return nil, err
	}

	var refundMediaList []*models.RefundMedia

	if request.Photo1 != nil {
		photoFile, err := request.Photo1.Open()
		if err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     "photo",
		}

		refundMediaList = append(refundMediaList, refundMedia)
	}

	if request.Photo2 != nil {
		photoFile, err := request.Photo2.Open()
		if err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     "photo",
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Photo3 != nil {
		photoFile, err := request.Photo3.Open()
		if err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     "photo",
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Photo4 != nil {
		photoFile, err := request.Photo4.Open()
		if err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     "photo",
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Photo5 != nil {
		photoFile, err := request.Photo5.Open()
		if err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     "photo",
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Video != nil {
		videoFile, err := request.Video.Open()
		if err != nil {
			return nil, err
		}
		videoUrl, err := utils.UploadImage(videoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: videoUrl,
			Type:     "video",
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if len(refundMediaList) > 0 {
		_, err = s.refundMediaRepository.Create(refundMediaList)

		if err != nil {
			return nil, err
		}
	}

	updatedOrder := &models.Order{
		Status: constants.ComplainedStatus,
	}

	_, err = s.orderRepository.Update(orderID, updatedOrder)

	if err != nil {
		return nil, err
	}

	return dtos.NewCreateRefundResponse(refund), nil
}

func (s *refundService) GetAllDetailsByBuyer(userID uuid.UUID, query *dtos.GetAllRefundQuery) ([]*dtos.GetRefundDetailsResponse, *utils.Pagination, error) {
	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	refundDetails, pagination, err := s.refundRepository.GetAllDetailsByBuyer(userID, query, pagination)

	if err != nil {
		return nil, nil, err
	}

	refundDetailList := []*internalDtos.RefundDetails{}

	for _, rd := range refundDetails {
		refundItems, err := s.refundDetailRepository.GetDetailsByRefundID(rd.RefundID)

		if err != nil {
			return nil, nil, err
		}

		refundDetailList = append(refundDetailList, &internalDtos.RefundDetails{
			RefundID:                rd.RefundID,
			OrderID:                 rd.OrderID,
			InvoiceNumber:           rd.InvoiceNumber,
			MerchantID:              rd.MerchantID,
			ShopName:                rd.ShopName,
			UserID:                  rd.UserID,
			UserFullName:            rd.UserFullName,
			RefundRequestDatetime:   rd.RefundRequestDatetime,
			RejectedDatetime:        rd.RejectedDatetime,
			RefundStatus:            rd.RefundStatus,
			Images:                  rd.Images,
			Videos:                  rd.Videos,
			RefundItems:             refundItems,
			TotalUserRefundAmount:   rd.TotalUserRefundAmount,
			TotalMerchantLostAmount: rd.TotalMerchantLostAmount,
		})
	}

	if err != nil {
		return nil, nil, err
	}

	return dtos.NewGetAllRefundDetailsResponse(refundDetailList), pagination, nil
}

func (s *refundService) GetDetailsByID(userID uuid.UUID, refundID uuid.UUID) (*dtos.GetRefundDetailsResponse, error) {
	user, err := s.userRepository.GetByID(userID)

	if err != nil {
		return nil, err
	}

	refundDetails, err := s.refundRepository.GetDetailsByID(refundID)

	if err != nil {
		return nil, err
	}

	merchant, err := s.merchantRepository.GetByUserID(userID)

	if errors.Is(err, e.ErrMerchantNotFound) {
	} else if err != nil {
		return nil, err
	}

	if refundDetails.UserID != userID && merchant.UserID != userID && user.Role != "100" {
		return nil, e.ErrForbidden
	}

	refundItems, err := s.refundDetailRepository.GetDetailsByRefundID(refundID)

	if err != nil {
		return nil, err
	}

	refundDetails.RefundItems = refundItems

	return dtos.NewGetRefundDetailsResponse(refundDetails), nil
}

func (s *refundService) GetAllDetailsByMerchant(userID uuid.UUID, query *dtos.GetAllRefundQuery) ([]*dtos.GetRefundDetailsResponse, *utils.Pagination, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, nil, err
	}

	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	refundDetails, pagination, err := s.refundRepository.GetAllDetailsByMerchant(merchant.ID, query, pagination)

	if err != nil {
		return nil, nil, err
	}

	refundDetailList := []*internalDtos.RefundDetails{}

	for _, rd := range refundDetails {
		refundItems, err := s.refundDetailRepository.GetDetailsByRefundID(rd.RefundID)

		if err != nil {
			return nil, nil, err
		}

		refundDetailList = append(refundDetailList, &internalDtos.RefundDetails{
			RefundID:                rd.RefundID,
			OrderID:                 rd.OrderID,
			InvoiceNumber:           rd.InvoiceNumber,
			MerchantID:              rd.MerchantID,
			ShopName:                rd.ShopName,
			UserID:                  rd.UserID,
			UserFullName:            rd.UserFullName,
			RefundRequestDatetime:   rd.RefundRequestDatetime,
			RejectedDatetime:        rd.RejectedDatetime,
			RefundStatus:            rd.RefundStatus,
			Images:                  rd.Images,
			Videos:                  rd.Videos,
			RefundItems:             refundItems,
			TotalUserRefundAmount:   rd.TotalUserRefundAmount,
			TotalMerchantLostAmount: rd.TotalMerchantLostAmount,
		})
	}

	if err != nil {
		return nil, nil, err
	}

	return dtos.NewGetAllRefundDetailsResponse(refundDetailList), pagination, nil
}

func (s *refundService) GetAllDetailsByAdmin(userID uuid.UUID, query *dtos.GetAllRefundQuery) ([]*dtos.GetRefundDetailsResponse, *utils.Pagination, error) {
	user, err := s.userRepository.GetByID(userID)

	if err != nil {
		return nil, nil, err
	}

	if user.Role != "100" {
		return nil, nil, e.ErrForbidden
	}

	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	refundDetails, pagination, err := s.refundRepository.GetAllDetails(query, pagination)

	if err != nil {
		return nil, nil, err
	}

	refundDetailList := []*internalDtos.RefundDetails{}

	for _, rd := range refundDetails {
		refundItems, err := s.refundDetailRepository.GetDetailsByRefundID(rd.RefundID)

		if err != nil {
			return nil, nil, err
		}

		refundDetailList = append(refundDetailList, &internalDtos.RefundDetails{
			RefundID:                rd.RefundID,
			OrderID:                 rd.OrderID,
			InvoiceNumber:           rd.InvoiceNumber,
			MerchantID:              rd.MerchantID,
			ShopName:                rd.ShopName,
			UserID:                  rd.UserID,
			UserFullName:            rd.UserFullName,
			RefundRequestDatetime:   rd.RefundRequestDatetime,
			RejectedDatetime:        rd.RejectedDatetime,
			RefundStatus:            rd.RefundStatus,
			Images:                  rd.Images,
			Videos:                  rd.Videos,
			RefundItems:             refundItems,
			TotalUserRefundAmount:   rd.TotalUserRefundAmount,
			TotalMerchantLostAmount: rd.TotalMerchantLostAmount,
		})
	}

	if err != nil {
		return nil, nil, err
	}

	return dtos.NewGetAllRefundDetailsResponse(refundDetailList), pagination, nil
}
