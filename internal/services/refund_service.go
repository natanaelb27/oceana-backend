package services

import (
	"errors"
	"math"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
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
	UpdateApproveRefund(refundID uuid.UUID, userID uuid.UUID) (*models.Refund, error)
	UpdateRejectRefund(refundID uuid.UUID, userID uuid.UUID) (*models.Refund, error)
	UpdateCompleteRefund(refundID uuid.UUID, userID uuid.UUID) (*models.Refund, error)
	CalculateRefundDetails(orderDetailID uuid.UUID, query *dtos.CalculateRefundDetailsQuery) (interface{}, error)
}

type refundService struct {
	merchantRepository                     r.MerchantRepository
	orderRepository                        r.OrderRepository
	orderDetailRepository                  r.OrderDetailRepository
	refundRepository                       r.RefundRepository
	refundDetailRepository                 r.RefundDetailRepository
	refundMediaRepository                  r.RefundMediaRepository
	userRepository                         r.UserRepository
	transactionRepository                  r.TransactionRepository
	walletTransactionPaymentRepository     r.WalletTransactionPaymentRepository
	walletRepository                       r.WalletRepository
	sealabsPayTransactionPaymentRepository r.SealabsPayTransactionPaymentRepository
	orderTrackingRepository                r.OrderTrackingRepository
	sealabsPayService                      SealabsPayService
}

type RefundServiceConfig struct {
	MerchantRepository                     r.MerchantRepository
	OrderRepository                        r.OrderRepository
	OrderDetailRepository                  r.OrderDetailRepository
	RefundRepository                       r.RefundRepository
	RefundDetailRepository                 r.RefundDetailRepository
	RefundMediaRepository                  r.RefundMediaRepository
	UserRepository                         r.UserRepository
	TransactionRepository                  r.TransactionRepository
	WalletTransactionPaymentRepository     r.WalletTransactionPaymentRepository
	WalletRepository                       r.WalletRepository
	SealabsPayTransactionPaymentRepository r.SealabsPayTransactionPaymentRepository
	OrderTrackingRepository                r.OrderTrackingRepository
	SealabsPayService                      SealabsPayService
}

func NewRefundService(c *RefundServiceConfig) RefundService {
	return &refundService{
		merchantRepository:                     c.MerchantRepository,
		orderRepository:                        c.OrderRepository,
		orderDetailRepository:                  c.OrderDetailRepository,
		refundRepository:                       c.RefundRepository,
		refundDetailRepository:                 c.RefundDetailRepository,
		refundMediaRepository:                  c.RefundMediaRepository,
		userRepository:                         c.UserRepository,
		transactionRepository:                  c.TransactionRepository,
		walletTransactionPaymentRepository:     c.WalletTransactionPaymentRepository,
		walletRepository:                       c.WalletRepository,
		sealabsPayTransactionPaymentRepository: c.SealabsPayTransactionPaymentRepository,
		orderTrackingRepository:                c.OrderTrackingRepository,
		sealabsPayService:                      c.SealabsPayService,
	}
}

func (s *refundService) Create(userID uuid.UUID, orderID uuid.UUID, request *dtos.CreateRefundRequest) (*dtos.CreateRefundResponse, error) {
	for _, v := range request.OrderDetailRefund {
		orderDetail, err := s.orderDetailRepository.GetByID(v.OrderDetailID)

		if err != nil {
			return nil, err
		}

		if orderDetail.OrderID != orderID {
			return nil, e.ErrOrderDoesNotHaveThatOrderDetail
		}
	}

	order, err := s.orderRepository.OrderListDetailsGetByID(orderID)

	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, e.ErrForbiddenOrder
	}

	if order.OrderStatus != constants.DeliveredStatus && order.OrderStatus != constants.ComplainedStatus {
		return nil, e.ErrStatusIsNotDelivered
	}

	existingRefunds, err := s.refundRepository.GetByOrderID(orderID)

	if err != nil {
		return nil, err
	}

	for _, v := range existingRefunds {
		if v.Status == constants.RefundStatusApproved {
			return nil, e.ErrCreateRefundIsOngoing
		}

		if v.Status == constants.RefundStatusApproved || v.Status == constants.RefundStatusApprovedBySystem {
			return nil, e.ErrCreateRefundIsApproved
		}

		if v.Status == constants.RefundStatusCompleted {
			return nil, e.ErrCreateRefundIsCompleted
		}
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

	for _, rd := range request.OrderDetailRefund {
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
		defer photoFile.Close()

		imageUrl, err := utils.UploadMedia(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     product_media_constants.TypeImage,
		}

		refundMediaList = append(refundMediaList, refundMedia)
	}

	if request.Photo2 != nil {
		photoFile, err := request.Photo2.Open()
		if err != nil {
			return nil, err
		}
		defer photoFile.Close()

		imageUrl, err := utils.UploadMedia(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     product_media_constants.TypeImage,
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Photo3 != nil {
		photoFile, err := request.Photo3.Open()
		if err != nil {
			return nil, err
		}
		defer photoFile.Close()

		imageUrl, err := utils.UploadMedia(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     product_media_constants.TypeImage,
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Photo4 != nil {
		photoFile, err := request.Photo4.Open()
		if err != nil {
			return nil, err
		}
		defer photoFile.Close()

		imageUrl, err := utils.UploadMedia(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     product_media_constants.TypeImage,
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Photo5 != nil {
		photoFile, err := request.Photo5.Open()
		if err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadMedia(photoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: imageUrl,
			Type:     product_media_constants.TypeImage,
		}

		refundMediaList = append(refundMediaList, refundMedia)

	}

	if request.Video != nil {
		videoFile, err := request.Video.Open()
		if err != nil {
			return nil, err
		}
		defer videoFile.Close()
		videoUrl, err := utils.UploadMedia(videoFile)
		if err != nil {
			return nil, err
		}

		refundMedia := &models.RefundMedia{
			RefundID: refund.ID,
			MediaURL: videoUrl,
			Type:     product_media_constants.TypeVideo,
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
			RefundReason:            rd.RefundReason,
			Images:                  rd.Images,
			Videos:                  rd.Videos,
			RefundItems:             refundItems,
			TotalUserRefundAmount:   rd.TotalUserRefundAmount,
			TotalMerchantLostAmount: rd.TotalMerchantLostAmount,
			RefundSellerConfirmationExpirationDatetime: rd.RefundSellerConfirmationExpirationDatetime,
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

	if refundDetails.UserID != userID && merchant.ID != refundDetails.MerchantID && user.Role != "100" {
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
			RefundReason:            rd.RefundReason,
			Images:                  rd.Images,
			Videos:                  rd.Videos,
			RefundItems:             refundItems,
			TotalUserRefundAmount:   rd.TotalUserRefundAmount,
			TotalMerchantLostAmount: rd.TotalMerchantLostAmount,
			RefundSellerConfirmationExpirationDatetime: rd.RefundSellerConfirmationExpirationDatetime,
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
			RefundReason:            rd.RefundReason,
			Images:                  rd.Images,
			Videos:                  rd.Videos,
			RefundItems:             refundItems,
			TotalUserRefundAmount:   rd.TotalUserRefundAmount,
			TotalMerchantLostAmount: rd.TotalMerchantLostAmount,
			RefundSellerConfirmationExpirationDatetime: rd.RefundSellerConfirmationExpirationDatetime,
		})
	}

	if err != nil {
		return nil, nil, err
	}

	return dtos.NewGetAllRefundDetailsResponse(refundDetailList), pagination, nil
}

func (s *refundService) UpdateApproveRefund(refundID uuid.UUID, userID uuid.UUID) (*models.Refund, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, err
	}

	refund, err := s.refundRepository.GetByID(refundID)
	if err != nil {
		return nil, err
	}

	if merchant.ID != refund.Order.MerchantID {
		return nil, e.ErrForbiddenMerchant
	}

	if refund.Status != constants.RefundStatusOngoing {
		return nil, e.ErrRefundStatusIsNotOngoing
	}

	refund.Status = constants.RefundStatusApproved

	updatedRefund, err := s.refundRepository.Update(refundID, refund)

	if err != nil {
		return nil, err
	}

	return updatedRefund, nil
}

func (s *refundService) UpdateRejectRefund(refundID uuid.UUID, userID uuid.UUID) (*models.Refund, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, err
	}

	refund, err := s.refundRepository.GetByID(refundID)
	if err != nil {
		return nil, err
	}

	if merchant.ID != refund.Order.MerchantID {
		return nil, e.ErrForbiddenMerchant
	}

	if refund.Status != constants.RefundStatusOngoing {
		return nil, e.ErrRefundStatusIsNotOngoing
	}

	refund.Status = constants.RefundStatusRejected

	timeNow := time.Now()
	refund.RejectedDatetime = &timeNow

	updatedRefund, err := s.refundRepository.Update(refundID, refund)

	if err != nil {
		return nil, err
	}

	return updatedRefund, nil
}

func (s *refundService) UpdateCompleteRefund(refundID uuid.UUID, userID uuid.UUID) (*models.Refund, error) {
	user, err := s.userRepository.GetByID(userID)

	if err != nil {
		return nil, err
	}

	if user.Role != "100" {
		return nil, e.ErrForbidden
	}

	refund, err := s.refundRepository.GetByID(refundID)

	if err != nil {
		return nil, err
	}

	if refund.Status != constants.RefundStatusApproved && refund.Status != constants.RefundStatusApprovedBySystem {
		return nil, e.ErrRefundStatusIsNotApproved
	}

	transaction, err := s.transactionRepository.GetByOrderID(refund.OrderID)

	if err != nil {
		return nil, err
	}

	refundDetail, err := s.refundRepository.GetDetailsByID(refundID)

	if err != nil {
		return nil, err
	}

	if transaction.PaymentMethodID == constants.WalletID {
		walletTransactionPayment, err := s.walletTransactionPaymentRepository.GetByTransactionID(transaction.ID)

		if err != nil {
			return nil, err
		}

		newWallet := &models.Wallet{
			Balance: walletTransactionPayment.Wallet.Balance + uint(refundDetail.TotalUserRefundAmount),
		}

		_, err = s.walletRepository.Update(walletTransactionPayment.WalletID, newWallet)

		if err != nil {
			return nil, err
		}

	} else if transaction.PaymentMethodID == constants.SealabsPayID {
		sealabsPayTransactionPayment, err := s.sealabsPayTransactionPaymentRepository.GetByTransactionID(transaction.ID)

		if err != nil {
			return nil, err
		}

		txIDInt, err := strconv.Atoi(sealabsPayTransactionPayment.TxID)

		if err != nil {
			return nil, err
		}

		err = s.sealabsPayService.Refund(&dtos.SealabsPayRefundRequest{
			TxnID:  txIDInt,
			Amount: int(refundDetail.TotalUserRefundAmount),
			Reason: refund.Reason,
		})

		if err != nil {
			return nil, err
		}
	}

	refund.Status = constants.RefundStatusCompleted

	updatedRefund, err := s.refundRepository.Update(refundID, refund)

	if err != nil {
		return nil, err
	}

	_, err = s.orderRepository.Update(refund.OrderID, &models.Order{
		Status: constants.CompletedStatus,
	})

	if err != nil {
		return nil, err
	}

	orderTracking := []*models.OrderTracking{
		{
			OrderID:         refund.OrderID,
			HistoryDatetime: time.Now(),
			Title:           constants.ConfirmedTitle,
			Description:     constants.ConfirmedDescription,
		},
		{
			OrderID:         refund.OrderID,
			HistoryDatetime: time.Now(),
			Title:           constants.CompletedTitle,
			Description:     constants.CompletedDescription,
		},
	}

	_, err = s.orderTrackingRepository.Create(orderTracking)

	if err != nil {
		return nil, err
	}

	return updatedRefund, nil
}

func (s *refundService) CalculateRefundDetails(orderDetailID uuid.UUID, query *dtos.CalculateRefundDetailsQuery) (interface{}, error) {
	refundDetails, err := s.refundDetailRepository.CalculateRefundDetails(orderDetailID, query.Quantity)

	if err != nil {
		return nil, err
	}

	refundDetails.Quantity = query.Quantity

	return dtos.NewCalculateRefundDetailsResponse(refundDetails), err
}
