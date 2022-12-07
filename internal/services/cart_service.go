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
	PayCartWithWallet(c *gin.Context, request *dtos.PayCartWithWalletRequest) (interface{}, error)
}

type cartService struct {
	addressSnapshotRepository          r.AddressSnapshotRepository
	cartRepository                     r.CartRepository
	productRepository                  r.ProductRepository
	productItemRepository              r.ProductItemRepository
	walletRepository                   r.WalletRepository
	productMerchantPromotionRepository r.ProductMerchantPromotionRepository
	productSnapshotRepository          r.ProductSnapshotRepository
	productMediaSnapshotRepository     r.ProductMediaSnapshotRepository
	productVariationSnapshotRepository r.ProductVariationSnapshotRepository
	merchantVoucherRepository          r.MerchantVoucherRepository
	merchantPromotionRepository        r.MerchantPromotionRepository
	marketplaceVoucherRepository       r.MarketplaceVoucherRepository
	redisRepository                    r.RedisRepository
	userMerchantPromotionRepository    r.UserMerchantPromotionRepository
}

type CartServiceConfig struct {
	AddressSnapshotRepository          r.AddressSnapshotRepository
	CartRepository                     r.CartRepository
	ProductRepository                  r.ProductRepository
	ProductItemRepository              r.ProductItemRepository
	WalletRepository                   r.WalletRepository
	ProductMerchantPromotionRepository r.ProductMerchantPromotionRepository
	ProductSnapshotRepository          r.ProductSnapshotRepository
	ProductVariationSnapshotRepository r.ProductVariationSnapshotRepository
	ProductMediaSnapshotRepository     r.ProductMediaSnapshotRepository
	MerchantVoucherRepository          r.MerchantVoucherRepository
	MerchantPromotionRepository        r.MerchantPromotionRepository
	MarketplaceVoucherRepository       r.MarketplaceVoucherRepository
	RedisRepository                    r.RedisRepository
	UserMerchantPromotionRepository    r.UserMerchantPromotionRepository
}

func NewCartService(c *CartServiceConfig) CartService {
	return &cartService{
		addressSnapshotRepository:          c.AddressSnapshotRepository,
		cartRepository:                     c.CartRepository,
		productRepository:                  c.ProductRepository,
		productItemRepository:              c.ProductItemRepository,
		walletRepository:                   c.WalletRepository,
		productMerchantPromotionRepository: c.ProductMerchantPromotionRepository,
		productSnapshotRepository:          c.ProductSnapshotRepository,
		productVariationSnapshotRepository: c.ProductVariationSnapshotRepository,
		productMediaSnapshotRepository:     c.ProductMediaSnapshotRepository,
		merchantVoucherRepository:          c.MerchantVoucherRepository,
		merchantPromotionRepository:        c.MerchantPromotionRepository,
		marketplaceVoucherRepository:       c.MarketplaceVoucherRepository,
		redisRepository:                    c.RedisRepository,
		userMerchantPromotionRepository:    c.UserMerchantPromotionRepository,
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

func (p *cartService) PayCartWithWallet(c *gin.Context, request *dtos.PayCartWithWalletRequest) (interface{}, error) {
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

	// validasi marketplace voucher

	transaction := &models.Transaction{
		ID:                  uuid.New(),
		UserID:              request.UserID,
		TransactionDatetime: time.Now(),
		PaymentMethodID:     1, //ganti constant
		IsPaid:              true,
	}

	// fmt.Println(transaction)

	orderList := []*models.Order{}
	orderDetailList := []*models.OrderDetail{}
	orderDetailMerchantPromotionList := []*models.OrderDetailMerchantPromotion{}
	orderMerchantVoucher := []*models.OrderMerchantVoucher{}

	var sumOrder float64 = 0
	var sumShippingFee uint = 0

	addressSnapshot, err := p.addressSnapshotRepository.Create(request.AddressID)

	if err != nil {
		return nil, err
	}

	for _, cartOrderItem := range request.CartOrderItems {
		merchantVoucher := &internalDtos.MerchantVoucherDetails{}

		if cartOrderItem.MerchantVoucherID != nil {
			resultMerchantVoucherID, err := p.merchantVoucherRepository.ValidateMerchantVoucherByBuyer(*cartOrderItem.MerchantVoucherID, request.UserID)

			if errors.Is(err, e.ErrDataNotFound) {
				break
			} else if err != nil {
				return nil, err
			}

			if resultMerchantVoucherID != nil {
				return nil, errors.New("merchant voucher used by user")
			}

			merchantVoucher, err = p.merchantVoucherRepository.GetDetailsByID(*cartOrderItem.MerchantVoucherID)

			if err != nil {
				return nil, err
			}

			if cartOrderItem.MerchantID != merchantVoucher.MerchantID {
				return nil, errors.New("error merchant id")
			}

			if merchantVoucher.VoucherStatus != constants.OngoingVoucherStatus {
				return nil, errors.New("merchant voucher not ongoing")
			}

			if merchantVoucher.CurrentQuota >= merchantVoucher.MaximumQuota {
				return nil, errors.New("merchant voucher udah full quotanya")
			}
		}

		newOrder := &models.Order{
			ID:                uuid.New(),
			InvoiceNumber:     "get count order dl trus ntr disambung",
			TransactionID:     transaction.ID,
			MerchantID:        cartOrderItem.MerchantID,
			AddressSnapshotID: addressSnapshot.ID,
			CourierServiceID:  cartOrderItem.CourierServiceID,
			ShippingFee:       cartOrderItem.ShippingFee,
			Status:            constants.PendingSellerStatus,
			ClaimStatus:       false,
		}

		var sumProduct float64 = 0

		for _, product := range cartOrderItem.Products {
			productMerchantPromotion := &internalDtos.ProductMerchantPromotionDetails{}

			if product.ProductMerchantPromotionID != nil {
				productMerchantPromotion, err := p.productMerchantPromotionRepository.GetDetailsByID(*product.ProductMerchantPromotionID)

				if err != nil {
					return nil, err
				}

				if productMerchantPromotion.PromotionStatus != constants.OngoingVoucherStatus {
					return nil, errors.New("merchant promotion not ongoing")
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
					return nil, errors.New("MERCHANT GA PNY PRODUCT ITEM ITU")
				}

				if pi.StockQuantity < productItems.Quantity {
					return nil, errors.New("PRODUK STOCK KURANG")
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
					// if !reflect.DeepEqual(productMerchantPromotion, &internalDtos.ProductMerchantPromotionDetails{}) {
					if strings.EqualFold(productMerchantPromotion.DiscountType, constants.NominalDiscountType) {
						totalDiscountODMP = float64(productMerchantPromotion.Discount * productItems.Quantity)
						totalPriceAfterDiscountOrderDetails = math.Max((float64(totalPriceOrderDetails) - totalDiscountODMP), 0)
					} else {
						totalDiscountODMP = math.Min(float64(productMerchantPromotion.Discount/100.0*pi.Price), float64(productMerchantPromotion.MaximumDiscountAmount)) * float64(productItems.Quantity)
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
					TotalPriceAfterDiscount: totalPriceAfterDiscountOrderDetails,
				}

				orderDetailList = append(orderDetailList, newOrderDetail)

				if product.ProductMerchantPromotionID != nil {
					// if !reflect.DeepEqual(productMerchantPromotion, &internalDtos.ProductMerchantPromotionDetails{}) {
					newOrderDetailMerchantPromotion := &models.OrderDetailMerchantPromotion{
						OrderDetailID:       newOrderDetail.ID,
						MerchantPromotionID: productMerchantPromotion.ProductMerchantPromotionID,
						TotalDiscount:       totalDiscountODMP,
					}

					orderDetailMerchantPromotionList = append(orderDetailMerchantPromotionList, newOrderDetailMerchantPromotion)
				}

				fmt.Println("order detail: ", totalPriceOrderDetails)
				fmt.Println("order detail disc: ", totalDiscountODMP)
				fmt.Println("order detail a4 disc: ", totalPriceAfterDiscountOrderDetails)
			}

			if product.ProductMerchantPromotionID != nil {
				// if !reflect.DeepEqual(productMerchantPromotion, &internalDtos.ProductMerchantPromotionDetails{}) {
				if productMerchantPromotion.MaximumQuota != -1 && int(productMerchantPromotion.CurrentQuota+sumQuantity) > productMerchantPromotion.MaximumQuota {
					return nil, errors.New("PURCHASED QTY MELEBIHI PMP GLOBAL QUOTA LIMIT")
				}

				userMerchantPromotion, err := p.userMerchantPromotionRepository.GetUserUsage(request.UserID, *product.ProductMerchantPromotionID)

				if errors.Is(err, e.ErrDataNotFound) {
					break
				} else if err != nil {
					return nil, err
				}

				if productMerchantPromotion.MaximumUserLimit != -1 && int(userMerchantPromotion.Usage+sumQuantity) > productMerchantPromotion.MaximumUserLimit {
					return nil, errors.New("PURCHASED QTY MELEBIHI PMP USER LIMIT")
				}

				newUserMerchantPromotionMap := map[string]any{
					"user_id":                       request.UserID,
					"product_merchant_promotion_id": product.ProductMerchantPromotionID,
					"usage":                         userMerchantPromotion.Usage + sumQuantity,
				}

				// ADD USER USAGE LIMIT DI USER MERCHANT PROMO & ADD CURRENT LIMIT TO PMP BY sumQuantity

				_, err = p.productMerchantPromotionRepository.AddCurrentQuota(*product.ProductMerchantPromotionID, sumQuantity)

				if err != nil {
					return nil, err
				}

				_, err = p.userMerchantPromotionRepository.UpsertUsingMap(newUserMerchantPromotionMap)

				if err != nil {
					return nil, err
				}

			}

			fmt.Println("order sum: ", sumProductItem)

			sumProduct += sumProductItem
		}

		var totalPriceAfterDiscountOrders float64 = sumProduct
		var totalDiscountOMV float64 = 0

		if cartOrderItem.MerchantVoucherID != nil {
			if merchantVoucher.MinimumSpendAmount > uint(sumProduct) {
				return nil, errors.New("min spend amount > sum order")
			}

			if strings.EqualFold(merchantVoucher.DiscountType, constants.NominalDiscountType) {
				totalDiscountOMV = float64(merchantVoucher.Discount)
				totalPriceAfterDiscountOrders = math.Max((sumProduct - totalDiscountOMV), 0)
			} else {
				totalDiscountOMV = math.Min(float64(merchantVoucher.Discount/100.0*uint(sumProduct)), float64(merchantVoucher.MaximumDiscountAmount))
				totalPriceAfterDiscountOrders = sumProduct - totalDiscountOMV
			}

		}

		newOrder.TotalPrice = sumProduct
		newOrder.TotalPriceAfterDiscount = totalPriceAfterDiscountOrders

		orderList = append(orderList, newOrder)

		if cartOrderItem.MerchantVoucherID != nil {
			newOrderMerchantVoucher := &models.OrderMerchantVoucher{
				OrderID:           newOrder.ID,
				MerchantVoucherID: *cartOrderItem.MerchantVoucherID,
				TotalDiscount:     totalDiscountOMV,
			}

			orderMerchantVoucher = append(orderMerchantVoucher, newOrderMerchantVoucher)
		}

		sumOrder += sumProduct
		sumShippingFee += cartOrderItem.ShippingFee
	}

	transaction.TotalPrice = sumOrder

	// TRANSACTION IS PAID SET TO TRUE KARNA WALLET

	// for _, cartOrderItem := range request.CartOrderItems {
	// 	merchantVoucher, err := p.merchantVoucherRepository.GetDetailsByID(*cartOrderItem.MerchantVoucherID)

	// 	if err != nil && err != e.ErrDataNotFound {
	// 		return nil, err
	// 	}

	// 	for _, orderCartItem := range cartOrderItem.OrderCartItems {
	// 		orderCartItemPrice := 0
	// 		if cartOrderItem.MerchantPromotionID != nil {
	// 			productItem, err := p.productItemRepository.GetById(orderCartItem.ProductItemID)

	// 			if err != nil {
	// 				return nil, err
	// 			}

	// 			productMerchantPromotion, err := p.productMerchantPromotionRepository.GetOnGoingByProductID(productItem.ProductID)

	// 			if err != nil && err != e.ErrDataNotFound {
	// 				return nil, err
	// 			}

	// 			if err == nil {
	// 				totalPrice := productItem.Price * *orderCartItem.Quantity
	// 				discount := uint(0)

	// 				if productMerchantPromotion.DiscountType == constants.NominalDiscountType {
	// 					discount = productMerchantPromotion.Discount
	// 				}

	// 				if productMerchantPromotion.DiscountType == constants.PercentageDiscountType {
	// 					discount = uint(float32(float32(totalPrice) * float32(discount/100)))

	// 					if discount > productMerchantPromotion.MaximumDiscountAmount {
	// 						discount = productMerchantPromotion.MaximumDiscountAmount
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// TO-DO: validate user wallet balance

	// TO-DO: validate user wallet balance

	// TO-DO: validate product quantity

	// TO-DO: create snapshot for bunch of stuffs

	// TO-DO: do payment stuffs

	return request, nil
}
