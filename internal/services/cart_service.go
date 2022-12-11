package services

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	in_cart_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/cart_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartService interface {
	Update(request *in_cart_dtos.CartRequest) (*models.Cart, error)
	GetAllByUserId(userId uuid.UUID) ([]*in_cart_dtos.Merchant, error)
	PayCartWithWallet(c *gin.Context, request *dtos.PayCartWithWalletRequest) (*dtos.WalletPaymentResponse, error)
}

type cartService struct {
	addressSnapshotRepository               r.AddressSnapshotRepository
	cartRepository                          r.CartRepository
	orderRepository                         r.OrderRepository
	orderDetailRepository                   r.OrderDetailRepository
	orderDetailMerchantPromotionRepository  r.OrderDetailMerchantPromotionRepository
	orderMerchantVoucherRepository          r.OrderMerchantVoucherRepository
	orderTrackingRepository                 r.OrderTrackingRepository
	productRepository                       r.ProductRepository
	productItemRepository                   r.ProductItemRepository
	productMerchantPromotionRepository      r.ProductMerchantPromotionRepository
	productSnapshotRepository               r.ProductSnapshotRepository
	productMediaSnapshotRepository          r.ProductMediaSnapshotRepository
	productVariationSnapshotRepository      r.ProductVariationSnapshotRepository
	merchantVoucherRepository               r.MerchantVoucherRepository
	merchantPromotionRepository             r.MerchantPromotionRepository
	marketplaceVoucherRepository            r.MarketplaceVoucherRepository
	redisRepository                         r.RedisRepository
	transactionRepository                   r.TransactionRepository
	transactionMarketplaceVoucherRepository r.TransactionMarketplaceVoucherRepository
	userMerchantPromotionRepository         r.UserMerchantPromotionRepository
	walletRepository                        r.WalletRepository
	walletTransactionPaymentRepository      r.WalletTransactionPaymentRepository
}

type CartServiceConfig struct {
	AddressSnapshotRepository               r.AddressSnapshotRepository
	CartRepository                          r.CartRepository
	OrderRepository                         r.OrderRepository
	OrderDetailRepository                   r.OrderDetailRepository
	OrderDetailMerchantPromotionRepository  r.OrderDetailMerchantPromotionRepository
	OrderMerchantVoucherRepository          r.OrderMerchantVoucherRepository
	OrderTrackingRepository                 r.OrderTrackingRepository
	ProductRepository                       r.ProductRepository
	ProductItemRepository                   r.ProductItemRepository
	ProductMerchantPromotionRepository      r.ProductMerchantPromotionRepository
	ProductSnapshotRepository               r.ProductSnapshotRepository
	ProductMediaSnapshotRepository          r.ProductMediaSnapshotRepository
	ProductVariationSnapshotRepository      r.ProductVariationSnapshotRepository
	MerchantVoucherRepository               r.MerchantVoucherRepository
	MerchantPromotionRepository             r.MerchantPromotionRepository
	MarketplaceVoucherRepository            r.MarketplaceVoucherRepository
	RedisRepository                         r.RedisRepository
	TransactionRepository                   r.TransactionRepository
	TransactionMarketplaceVoucherRepository r.TransactionMarketplaceVoucherRepository
	UserMerchantPromotionRepository         r.UserMerchantPromotionRepository
	WalletRepository                        r.WalletRepository
	WalletTransactionPaymentRepository      r.WalletTransactionPaymentRepository
}

func NewCartService(c *CartServiceConfig) CartService {
	return &cartService{
		addressSnapshotRepository:               c.AddressSnapshotRepository,
		cartRepository:                          c.CartRepository,
		orderRepository:                         c.OrderRepository,
		orderDetailRepository:                   c.OrderDetailRepository,
		orderDetailMerchantPromotionRepository:  c.OrderDetailMerchantPromotionRepository,
		orderMerchantVoucherRepository:          c.OrderMerchantVoucherRepository,
		orderTrackingRepository:                 c.OrderTrackingRepository,
		productRepository:                       c.ProductRepository,
		productItemRepository:                   c.ProductItemRepository,
		productMerchantPromotionRepository:      c.ProductMerchantPromotionRepository,
		productSnapshotRepository:               c.ProductSnapshotRepository,
		productMediaSnapshotRepository:          c.ProductMediaSnapshotRepository,
		productVariationSnapshotRepository:      c.ProductVariationSnapshotRepository,
		merchantVoucherRepository:               c.MerchantVoucherRepository,
		merchantPromotionRepository:             c.MerchantPromotionRepository,
		marketplaceVoucherRepository:            c.MarketplaceVoucherRepository,
		redisRepository:                         c.RedisRepository,
		transactionRepository:                   c.TransactionRepository,
		transactionMarketplaceVoucherRepository: c.TransactionMarketplaceVoucherRepository,
		userMerchantPromotionRepository:         c.UserMerchantPromotionRepository,
		walletRepository:                        c.WalletRepository,
		walletTransactionPaymentRepository:      c.WalletTransactionPaymentRepository,
	}
}

func (p *cartService) Update(request *in_cart_dtos.CartRequest) (result *models.Cart, err error) {
	cart, cartErr := p.cartRepository.GetByUserIdProductItemId(request.UserID, request.ProductItemID)
	if cartErr != nil && cartErr != e.ErrDataNotFound {
		err = cartErr
		return
	}

	// check if stock available
	productItem, err := p.productItemRepository.GetById(request.ProductItemID)
	if err != nil {
		return
	}
	if productItem.StockQuantity < *request.Quantity {
		err = e.ErrStockIsNotAvailable
		return
	}

	if cartErr == e.ErrDataNotFound {
		if *request.Quantity == 0 {
			err = e.ErrEmptyQuantity
			return
		}
		return p.cartRepository.Create(request.ToModel())
	}

	if *request.Quantity == 0 && err != e.ErrDataNotFound {
		return p.cartRepository.Delete(cart.ID)
	}

	cart.Quantity = *request.Quantity

	return p.cartRepository.Update(cart.ID, cart)
}

func (p *cartService) GetAllByUserId(userId uuid.UUID) ([]*in_cart_dtos.Merchant, error) {
	cartsMerchants, err := p.cartRepository.GetCartsMerchants(userId)
	for k, v := range cartsMerchants {
		cartsMerchants[k].Products, err = p.cartRepository.GetCartsProductsByUserIdMerchants(v.MerchantID, userId)
	}
	return cartsMerchants, err
}

func (p *cartService) PayCartWithWallet(c *gin.Context, request *dtos.PayCartWithWalletRequest) (*dtos.WalletPaymentResponse, error) {
	redisKey := "wallet-pin-failed-count:" + request.UserID.String()
	redisValue, err := p.redisRepository.Get(c, redisKey)

	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	var failedAttemptCount int

	if err == e.ErrDataNotFound {
		failedAttemptCount = 0
	} else {
		value, err := strconv.Atoi(redisValue)

		if err != nil {
			failedAttemptCount = 0
		} else {
			failedAttemptCount = value
		}
	}

	if failedAttemptCount >= constants.MaxWalletFailAttempt {
		return nil, e.ErrBlockedWallet
	}

	wallet, err := p.walletRepository.GetByUserId(request.UserID)

	if err != nil {
		return nil, err
	}

	if err := utils.CheckHash(wallet.Pin, request.Pin); err != nil {
		if err := p.redisRepository.Set(c, redisKey, fmt.Sprintf("%d", failedAttemptCount+1), constants.BlockDuration); err != nil {
			return nil, err
		}

		return nil, e.ErrWrongWalletPin
	}

	marketplaceVoucher := &internalDtos.MarketplaceVoucherDetails{}

	if request.MarketplaceVoucherID != nil {
		marketplaceVoucher, err = p.marketplaceVoucherRepository.GetDetailsByID(*request.MarketplaceVoucherID)

		if err != nil {
			return nil, err
		}

		if marketplaceVoucher.VoucherStatus != constants.OngoingPromotionOrVoucher {
			return nil, e.ErrMarketplaceVoucherNotOngoing
		}

		if marketplaceVoucher.CurrentQuota >= marketplaceVoucher.MaximumQuota {
			return nil, e.ErrMarketplaceVoucherQuotaFull
		}

	}

	newTransaction := &models.Transaction{
		ID:                  uuid.New(),
		UserID:              request.UserID,
		TransactionDatetime: time.Now(),
		PaymentMethodID:     constants.WalletID,
		IsPaid:              true,
	}

	newOrderTrackingList := []*models.OrderTracking{}
	newOrderList := []*models.Order{}
	newOrderDetailList := []*models.OrderDetail{}
	newOrderDetailMerchantPromotionList := []*models.OrderDetailMerchantPromotion{}
	newOrderMerchantVoucherList := []*models.OrderMerchantVoucher{}
	newTransactionMarketplaceVoucher := &models.TransactionMarketplaceVoucher{}

	var sumOrder float64 = 0
	var sumShippingFee float64 = 0

	addressSnapshot, err := p.addressSnapshotRepository.Create(request.AddressID)

	if err != nil {
		return nil, err
	}

	for _, cartOrderItem := range request.CartOrderItems {
		merchantVoucher := &internalDtos.MerchantVoucherDetails{}

		if cartOrderItem.MerchantVoucherID != nil {
			merchantVoucher, err = p.merchantVoucherRepository.GetDetailsByID(*cartOrderItem.MerchantVoucherID)

			if err != nil {
				return nil, err
			}

			if cartOrderItem.MerchantID != merchantVoucher.MerchantID {
				return nil, e.ErrRequestedVoucherNotMerchant
			}

			if merchantVoucher.VoucherStatus != constants.OngoingPromotionOrVoucher {
				return nil, e.ErrMerchantVoucherNotOngoing
			}

			if merchantVoucher.CurrentQuota >= merchantVoucher.MaximumQuota {
				return nil, e.ErrMerchantVoucherQuotaFull
			}

			validateMerchantVoucherID, err := p.merchantVoucherRepository.ValidateMerchantVoucherByBuyer(*cartOrderItem.MerchantVoucherID, request.UserID)

			if errors.Is(err, e.ErrDataNotFound) {

			} else if err != nil {
				return nil, err
			}

			if validateMerchantVoucherID != nil {
				return nil, e.ErrMerchantVoucherUsedByUser
			}
		}

		countOrder, err := p.orderRepository.GetOrderCount()

		if err != nil {
			return nil, err
		}

		newOrder := &models.Order{
			ID:                uuid.New(),
			InvoiceNumber:     fmt.Sprintf("INV/%s/%s", time.Now().Format("20060201"), strconv.Itoa(int(*countOrder)+1)),
			TransactionID:     newTransaction.ID,
			MerchantID:        cartOrderItem.MerchantID,
			AddressSnapshotID: addressSnapshot.ID,
			CourierServiceID:  cartOrderItem.CourierServiceID,
			ShippingFee:       cartOrderItem.ShippingFee,
			Status:            constants.PendingSellerStatus,
			ClaimStatus:       false,
		}

		newOrderTracking := &models.OrderTracking{
			OrderID:         newOrder.ID,
			HistoryDatetime: time.Now(),
			Title:           constants.PaymentVerifiedTitle,
			Description:     constants.PaymentVerifiedDescription,
		}

		newOrderTrackingList = append(newOrderTrackingList, newOrderTracking)

		var sumProduct float64 = 0

		for _, product := range cartOrderItem.Products {
			productMerchantPromotion := &internalDtos.ProductMerchantPromotionDetails{}

			if product.ProductMerchantPromotionID != nil {
				productMerchantPromotion, err = p.productMerchantPromotionRepository.GetDetailsByID(*product.ProductMerchantPromotionID)

				if err != nil {
					return nil, err
				}

				if productMerchantPromotion.ProductID != product.ProductID {
					return nil, e.ErrRequestedPromotionNotProduct
				}

				if productMerchantPromotion.PromotionStatus != constants.OngoingPromotionOrVoucher {
					return nil, e.ErrMerchantPromotionNotOngoing
				}
			}

			var sumProductItem float64 = 0
			var sumQuantity uint = 0

			for _, productItems := range product.ProductItems {
				pi, err := p.productItemRepository.GetDetailsById(productItems.ProductItemID)

				if err != nil {
					return nil, err
				}

				if pi.MerchantID != cartOrderItem.MerchantID {
					return nil, e.ErrMerchantDoesNotHaveThatProductItem
				}

				if pi.ProductID != product.ProductID {
					return nil, e.ErrProductDoesNotHaveThatProductItem
				}

				if pi.StockQuantity < productItems.Quantity {
					return nil, e.ErrProductStockNotEnough
				}

				newProductSnapshot := &models.ProductSnapshot{
					CategoryID:  pi.CategoryID,
					ProductID:   pi.ProductID,
					Name:        pi.Name,
					Description: pi.Description,
					IsHazardous: pi.IsHazardous,
					Weight:      pi.Weight,
					PackageSize: pi.PackageSize,
					Condition:   pi.Condition,
					SKU:         pi.SKU,
					Price:       pi.Price,
				}

				productSnapshot, err := p.productSnapshotRepository.CreateProductSnapshot(newProductSnapshot)

				if err != nil {
					return nil, err
				}

				_, err = p.productVariationSnapshotRepository.CreateProductVariationSnapshot(productSnapshot.ID, pi.ProductItemID)

				if err != nil {
					return nil, err
				}

				_, err = p.productMediaSnapshotRepository.CreateProductMediaSnapshot(productSnapshot.ID, pi.ProductItemID, pi.ProductID)

				if err != nil {
					return nil, err
				}

				totalPriceOrderDetails := pi.Price * productItems.Quantity

				sumQuantity += productItems.Quantity

				totalPriceAfterDiscountOrderDetails := float64(totalPriceOrderDetails)
				var totalDiscountODMP float64 = 0

				if product.ProductMerchantPromotionID != nil {
					if strings.EqualFold(productMerchantPromotion.DiscountType, constants.NominalDiscountType) {
						totalDiscountODMP = float64(productMerchantPromotion.Discount * productItems.Quantity)
						totalPriceAfterDiscountOrderDetails = math.Max((float64(totalPriceOrderDetails) - totalDiscountODMP), 0)
					} else {
						totalDiscountODMP = math.Min(float64(productMerchantPromotion.Discount)/100.0*float64(pi.Price), float64(productMerchantPromotion.MaximumDiscountAmount)) * float64(productItems.Quantity)
						totalPriceAfterDiscountOrderDetails = float64(totalPriceOrderDetails) - totalDiscountODMP
					}
				}

				sumProductItem += totalPriceAfterDiscountOrderDetails

				newOrderDetail := &models.OrderDetail{
					ID:                      uuid.New(),
					OrderID:                 newOrder.ID,
					Message:                 *productItems.Message,
					ProductSnapshotID:       productSnapshot.ID,
					Quantity:                productItems.Quantity,
					TotalPrice:              totalPriceOrderDetails,
					TotalPriceAfterDiscount: math.Round(totalPriceAfterDiscountOrderDetails),
				}

				newOrderDetailList = append(newOrderDetailList, newOrderDetail)

				if product.ProductMerchantPromotionID != nil {
					newOrderDetailMerchantPromotion := &models.OrderDetailMerchantPromotion{
						OrderDetailID:       newOrderDetail.ID,
						MerchantPromotionID: productMerchantPromotion.ProductMerchantPromotionID,
						TotalDiscount:       math.Round(totalDiscountODMP),
					}

					newOrderDetailMerchantPromotionList = append(newOrderDetailMerchantPromotionList, newOrderDetailMerchantPromotion)
				}

				_, err = p.productItemRepository.SubtractProductStock(productItems.ProductItemID, productItems.Quantity)

				if err != nil {
					return nil, err
				}

			}

			if product.ProductMerchantPromotionID != nil {

				if productMerchantPromotion.MaximumQuota != -1 && int(productMerchantPromotion.CurrentQuota+sumQuantity) > productMerchantPromotion.MaximumQuota {
					return nil, e.ErrPurchasedQuantityExceedRemainingMerchantPromotionQuota
				}

				userMerchantPromotion, err := p.userMerchantPromotionRepository.GetUserUsage(request.UserID, *product.ProductMerchantPromotionID)

				if errors.Is(err, e.ErrDataNotFound) {

				} else if err != nil {
					return nil, err
				} else {
					if productMerchantPromotion.MaximumUserLimit != -1 && int(userMerchantPromotion.Usage+sumQuantity) > productMerchantPromotion.MaximumUserLimit {
						return nil, e.ErrPurchasedQuantityExceedRemainingUserLimit
					}
				}

				newUserMerchantPromotionMap := map[string]any{
					"user_id":                       request.UserID,
					"product_merchant_promotion_id": product.ProductMerchantPromotionID,
					"usage":                         userMerchantPromotion.Usage + sumQuantity,
				}

				_, err = p.productMerchantPromotionRepository.AddCurrentQuota(*product.ProductMerchantPromotionID, sumQuantity)

				if err != nil {
					return nil, err
				}

				_, err = p.userMerchantPromotionRepository.UpsertUsingMap(newUserMerchantPromotionMap)

				if err != nil {
					return nil, err
				}

			}

			sumProduct += sumProductItem
		}

		var totalPriceAfterDiscountOrders float64 = sumProduct
		var totalDiscountOMV float64 = 0

		if cartOrderItem.MerchantVoucherID != nil {
			if merchantVoucher.MinimumSpendAmount > uint(sumProduct) {
				return nil, e.ErrPurchasedAmountNotExceedMerchantVoucherMinimumAmount
			}

			if strings.EqualFold(merchantVoucher.DiscountType, constants.NominalDiscountType) {
				totalDiscountOMV = float64(merchantVoucher.Discount)
				totalPriceAfterDiscountOrders = math.Max((sumProduct - totalDiscountOMV), 0)
			} else {
				totalDiscountOMV = math.Min(float64(merchantVoucher.Discount)/100.0*sumProduct, float64(merchantVoucher.MaximumDiscountAmount))
				totalPriceAfterDiscountOrders = sumProduct - totalDiscountOMV
			}
		}

		newOrder.TotalPrice = math.Round(sumProduct)
		newOrder.TotalPriceAfterDiscount = math.Round(totalPriceAfterDiscountOrders)

		newOrderList = append(newOrderList, newOrder)

		if cartOrderItem.MerchantVoucherID != nil {
			_, err = p.merchantVoucherRepository.AddCurrentQuota(*cartOrderItem.MerchantVoucherID, 1)

			if err != nil {
				return nil, err
			}

			newOrderMerchantVoucher := &models.OrderMerchantVoucher{
				OrderID:           newOrder.ID,
				MerchantVoucherID: *cartOrderItem.MerchantVoucherID,
				TotalDiscount:     math.Round(totalDiscountOMV),
			}

			newOrderMerchantVoucherList = append(newOrderMerchantVoucherList, newOrderMerchantVoucher)
		}

		sumOrder += totalPriceAfterDiscountOrders
		sumShippingFee += float64(cartOrderItem.ShippingFee)
	}

	var totalPriceAfterDiscountTransactions float64 = sumOrder
	var totalDiscountTMV float64 = 0

	if request.MarketplaceVoucherID != nil {
		if marketplaceVoucher.MinimumSpendAmount > uint(sumOrder) {
			return nil, e.ErrPurchasedAmountNotExceedMarketplaceVoucherMinimumAmount
		}

		if strings.EqualFold(marketplaceVoucher.DiscountType, constants.NominalDiscountType) {
			totalDiscountTMV = float64(marketplaceVoucher.Discount)
			totalPriceAfterDiscountTransactions = math.Max((sumOrder - totalDiscountTMV), 0)
		} else {
			totalDiscountTMV = math.Min(float64(marketplaceVoucher.Discount)/100.0*sumOrder, float64(marketplaceVoucher.MaximumDiscountAmount))
			totalPriceAfterDiscountTransactions = sumOrder - totalDiscountTMV
		}

		newTransactionMarketplaceVoucher = &models.TransactionMarketplaceVoucher{
			TransactionID:        newTransaction.ID,
			MarketplaceVoucherID: *request.MarketplaceVoucherID,
			TotalDiscount:        math.Round(totalDiscountTMV),
		}
	}

	newTransaction.TotalPrice = math.Round(sumOrder)
	newTransaction.TotalPriceAfterDiscount = math.Round(totalPriceAfterDiscountTransactions)
	newTransaction.TotalPriceWithShippingFee = math.Round(totalPriceAfterDiscountTransactions + sumShippingFee)

	if request.TotalPayment != newTransaction.TotalPriceWithShippingFee {
		return nil, e.ErrDifferentTotalPayment
	}

	if float64(wallet.Balance) < newTransaction.TotalPriceWithShippingFee {
		return nil, e.ErrWalletBalanceNotEnough
	}

	newTransactionCreated, err := p.transactionRepository.Create(newTransaction)

	if err != nil {
		return nil, err
	}

	newWalletTransactionPayment := &models.WalletTransactionPayment{
		TransactionID: newTransaction.ID,
		WalletID:      wallet.ID,
	}

	_, err = p.walletTransactionPaymentRepository.Create(newWalletTransactionPayment)

	if err != nil {
		return nil, err
	}

	if request.MarketplaceVoucherID != nil {
		_, err = p.marketplaceVoucherRepository.AddCurrentQuota(*request.MarketplaceVoucherID, 1)

		if err != nil {
			return nil, err
		}

		_, err = p.transactionMarketplaceVoucherRepository.Create(newTransactionMarketplaceVoucher)

		if err != nil {
			return nil, err
		}
	}

	_, err = p.orderRepository.Create(newOrderList)

	if err != nil {
		return nil, err
	}

	_, err = p.orderTrackingRepository.Create(newOrderTrackingList)

	if err != nil {
		return nil, err
	}

	if len(newOrderMerchantVoucherList) > 0 {
		_, err = p.orderMerchantVoucherRepository.Create(newOrderMerchantVoucherList)

		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	_, err = p.orderDetailRepository.Create(newOrderDetailList)

	if err != nil {
		return nil, err
	}

	if len(newOrderDetailMerchantPromotionList) > 0 {
		_, err = p.orderDetailMerchantPromotionRepository.Create(newOrderDetailMerchantPromotionList)

		if err != nil {
			return nil, err
		}
	}

	newWallet := &models.Wallet{
		Balance: wallet.Balance - uint(math.Round(newTransaction.TotalPriceWithShippingFee)),
	}

	_, err = p.walletRepository.Update(wallet.ID, newWallet)

	if err != nil {
		return nil, err
	}

	return dtos.NewWalletPaymentResponse(newTransactionCreated), nil
}
