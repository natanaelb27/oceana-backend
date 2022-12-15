package services

import (
	"math"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type RefundService interface {
	Create(userID uuid.UUID, orderID uuid.UUID, request *dtos.CreateRefundRequest) (*dtos.CreateRefundResponse, error)
}

type refundService struct {
	orderRepository        r.OrderRepository
	orderDetailRepository  r.OrderDetailRepository
	refundRepository       r.RefundRepository
	refundDetailRepository r.RefundDetailRepository
	refundMediaRepository  r.RefundMediaRepository
}

type RefundServiceConfig struct {
	OrderRepository        r.OrderRepository
	OrderDetailRepository  r.OrderDetailRepository
	RefundRepository       r.RefundRepository
	RefundDetailRepository r.RefundDetailRepository
	RefundMediaRepository  r.RefundMediaRepository
}

func NewRefundService(c *RefundServiceConfig) RefundService {
	return &refundService{
		orderRepository:        c.OrderRepository,
		orderDetailRepository:  c.OrderDetailRepository,
		refundRepository:       c.RefundRepository,
		refundDetailRepository: c.RefundDetailRepository,
		refundMediaRepository:  c.RefundMediaRepository,
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
