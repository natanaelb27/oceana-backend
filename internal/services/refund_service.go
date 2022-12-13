package services

import (
	"math"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
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
}

type RefundServiceConfig struct {
	OrderRepository        r.OrderRepository
	OrderDetailRepository  r.OrderDetailRepository
	RefundRepository       r.RefundRepository
	RefundDetailRepository r.RefundDetailRepository
}

func NewRefundService(c *RefundServiceConfig) RefundService {
	return &refundService{
		orderRepository:        c.OrderRepository,
		orderDetailRepository:  c.OrderDetailRepository,
		refundRepository:       c.RefundRepository,
		refundDetailRepository: c.RefundDetailRepository,
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

	// refundMedia, err := s.refundMediaRepository.Create()

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

	return nil, nil

	// return dtos.NewCreateRefundResponse(refund), nil
}
