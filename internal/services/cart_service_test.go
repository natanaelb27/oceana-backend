package services_test

import (
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/constants"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	in_cart_dtos "github.com/natanaelb27/oceana-backend/internal/dtos/cart_dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	"github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	externalDtos "github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_cartService_Update(t *testing.T) {
	var requestQuantity uint = 2
	var zeroRequestQuantity uint = 0
	var requestActionIcrement string = "increment"
	err := errors.New("error")

	requestCreate := &in_cart_dtos.CartRequest{
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		ProductItemID: 1,
		Quantity:      &requestQuantity,
		Action:        nil,
	}

	requestUpdate := &in_cart_dtos.CartRequest{
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		ProductItemID: 1,
		Quantity:      &requestQuantity,
		Action:        &requestActionIcrement,
	}

	productItemModel := &models.ProductItem{
		ID:            1,
		ProductID:     1,
		StockQuantity: 10,
	}

	modelCreate := &models.Cart{
		ID:            1,
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		ProductItemID: 1,
		Quantity:      2,
	}

	modelUpdate := &models.Cart{
		ID:            1,
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		ProductItemID: 1,
		Quantity:      3,
	}

	type fields struct {
		cartRepository        *repositories.MockCartRepository
		productItemRepository *repositories.MockProductItemRepository
	}
	type args struct {
		args1 *in_cart_dtos.CartRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args, w *models.Cart)
		want    *models.Cart
		wantErr bool
	}{
		{
			args: args{
				args1: requestCreate,
			},
			mock: func(r *fields, a args, w *models.Cart) {
				r.cartRepository.On("GetByUserIdProductItemId", mock.Anything, mock.Anything).Return(&models.Cart{}, e.ErrDataNotFound)
				r.productItemRepository.On("GetById", mock.Anything).Return(productItemModel, nil)
				r.cartRepository.On("Create", mock.Anything).Return(modelCreate, nil)
			},
			want: modelCreate,
		},
		{
			args: args{
				args1: requestUpdate,
			},
			mock: func(r *fields, a args, w *models.Cart) {
				r.cartRepository.On("GetByUserIdProductItemId", mock.Anything, mock.Anything).Return(&models.Cart{
					Quantity: 1,
				}, nil)
				r.productItemRepository.On("GetById", mock.Anything).Return(productItemModel, nil)
				r.cartRepository.On("Update", mock.Anything, mock.Anything).Return(modelUpdate, nil)
			},
			want: modelUpdate,
		},
		{
			args: args{
				args1: requestCreate,
			},
			mock: func(r *fields, a args, w *models.Cart) {
				r.cartRepository.On("GetByUserIdProductItemId", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: requestCreate,
			},
			mock: func(r *fields, a args, w *models.Cart) {
				r.cartRepository.On("GetByUserIdProductItemId", mock.Anything, mock.Anything).Return(&models.Cart{}, nil)
				r.productItemRepository.On("GetById", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: requestCreate,
			},
			mock: func(r *fields, a args, w *models.Cart) {
				r.cartRepository.On("GetByUserIdProductItemId", mock.Anything, mock.Anything).Return(&models.Cart{}, e.ErrDataNotFound)
				r.productItemRepository.On("GetById", mock.Anything).Return(&models.ProductItem{
					ID:            1,
					ProductID:     1,
					StockQuantity: 1,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: &in_cart_dtos.CartRequest{
					UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
					ProductItemID: 1,
					Quantity:      &zeroRequestQuantity,
					Action:        nil,
				},
			},
			mock: func(r *fields, a args, w *models.Cart) {
				r.cartRepository.On("GetByUserIdProductItemId", mock.Anything, mock.Anything).Return(&models.Cart{}, e.ErrDataNotFound)
				r.productItemRepository.On("GetById", mock.Anything).Return(productItemModel, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: &in_cart_dtos.CartRequest{
					UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
					ProductItemID: 1,
					Quantity:      &zeroRequestQuantity,
					Action:        &requestActionIcrement,
				},
			},
			mock: func(r *fields, a args, w *models.Cart) {
				r.cartRepository.On("GetByUserIdProductItemId", mock.Anything, mock.Anything).Return(&models.Cart{}, nil)
				r.productItemRepository.On("GetById", mock.Anything).Return(productItemModel, nil)
				r.cartRepository.On("Delete", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				repositories.NewMockCartRepository(t),
				repositories.NewMockProductItemRepository(t),
			}

			cartService := services.NewCartService(&services.CartServiceConfig{
				CartRepository:        tt.fields.cartRepository,
				ProductItemRepository: tt.fields.productItemRepository,
			})

			tt.mock(&tt.fields, tt.args, tt.want)

			got, err := cartService.Update(tt.args.args1)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cartService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cartService_GetAllByUserId(t *testing.T) {
	userID := uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851")

	model := []*in_cart_dtos.Merchant{{
		MerchantID: 1,
		Products:   []*in_cart_dtos.Product{},
	}}

	type fields struct {
		cartRepository *repositories.MockCartRepository
	}
	type args struct {
		args1 uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args, w []*in_cart_dtos.Merchant)
		want    []*in_cart_dtos.Merchant
		wantErr bool
	}{
		{
			args: args{
				args1: userID,
			},
			mock: func(r *fields, a args, w []*in_cart_dtos.Merchant) {
				r.cartRepository.On("GetCartsMerchants", mock.Anything).Return([]*in_cart_dtos.Merchant{{
					MerchantID: 1,
				}}, nil)
				r.cartRepository.On("GetCartsProductsByUserIdMerchants", mock.Anything, mock.Anything).Return([]*in_cart_dtos.Product{}, nil)
			},
			want: model,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				repositories.NewMockCartRepository(t),
			}

			cartService := services.NewCartService(&services.CartServiceConfig{
				CartRepository: tt.fields.cartRepository,
			})

			tt.mock(&tt.fields, tt.args, tt.want)

			got, err := cartService.GetAllByUserId(tt.args.args1)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartService.GetAllByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cartService.GetAllByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cartService_PayCartWithWallet(t *testing.T) {
	hashedPIN, _ := utils.Hash("123456")
	falseHashedPIN, _ := utils.Hash("098765")
	var marketplaceVoucherID uint = 1
	var merchantVoucherID uint = 1
	var productMerchantPromotionID uint = 1
	var countOrder int64 = 1
	var message string = "message"
	var cartIDs []*uint = make([]*uint, 1)
	var cartID uint = 1
	cartIDs[0] = &cartID

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := &externalDtos.PayCartWithWalletRequest{
		UserID:       uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		Pin:          "123456",
		TotalPayment: 10000,
		CartIDs:      cartIDs,
		CartOrderItems: []*externalDtos.CartOrderItem{
			{
				MerchantID:        1,
				MerchantVoucherID: &merchantVoucherID,
				CourierServiceID:  1,
				ShippingFee:       10000,
				Products: []*externalDtos.Product{
					{
						ProductID:                  1,
						ProductMerchantPromotionID: &productMerchantPromotionID,
						ProductItems: []*externalDtos.ProductItem{
							{
								ProductItemID: 1,
								Quantity:      5,
								Message:       &message,
							},
						},
					},
				},
			},
		},
		MarketplaceVoucherID: &marketplaceVoucherID,
	}

	falseRequestTotalPayment := &externalDtos.PayCartWithWalletRequest{
		UserID:       uuid.UUID{},
		Pin:          "123456",
		TotalPayment: 123,
		CartOrderItems: []*externalDtos.CartOrderItem{
			{
				MerchantID:        1,
				MerchantVoucherID: &merchantVoucherID,
				CourierServiceID:  1,
				ShippingFee:       10000,
				Products: []*externalDtos.Product{
					{
						ProductID:                  1,
						ProductMerchantPromotionID: &productMerchantPromotionID,
						ProductItems: []*externalDtos.ProductItem{
							{
								ProductItemID: 1,
								Quantity:      5,
								Message:       &message,
							},
						},
					},
				},
			},
		},
		MarketplaceVoucherID: &marketplaceVoucherID,
	}

	err := errors.New("error")
	model := &externalDtos.PaymentResponse{}

	type fields struct {
		addressSnapshotRepository               *repositories.MockAddressSnapshotRepository
		cartRepository                          *repositories.MockCartRepository
		orderRepository                         *repositories.MockOrderRepository
		orderDetailRepository                   *repositories.MockOrderDetailRepository
		orderDetailMerchantPromotionRepository  *repositories.MockOrderDetailMerchantPromotionRepository
		orderMerchantVoucherRepository          *repositories.MockOrderMerchantVoucherRepository
		orderTrackingRepository                 *repositories.MockOrderTrackingRepository
		productRepository                       *repositories.MockProductRepository
		productItemRepository                   *repositories.MockProductItemRepository
		productCourierServiceRepository         *repositories.MockProductCourierServiceRepository
		productMerchantPromotionRepository      *repositories.MockProductMerchantPromotionRepository
		productSnapshotRepository               *repositories.MockProductSnapshotRepository
		productMediaSnapshotRepository          *repositories.MockProductMediaSnapshotRepository
		productVariationSnapshotRepository      *repositories.MockProductVariationSnapshotRepository
		merchantVoucherRepository               *repositories.MockMerchantVoucherRepository
		merchantPromotionRepository             *repositories.MockMerchantPromotionRepository
		marketplaceVoucherRepository            *repositories.MockMarketplaceVoucherRepository
		redisRepository                         *repositories.MockRedisRepository
		transactionRepository                   *repositories.MockTransactionRepository
		transactionMarketplaceVoucherRepository *repositories.MockTransactionMarketplaceVoucherRepository
		userMerchantPromotionRepository         *repositories.MockUserMerchantPromotionRepository
		walletRepository                        *repositories.MockWalletRepository
		walletTransactionPaymentRepository      *repositories.MockWalletTransactionPaymentRepository
	}

	type args struct {
		args1 *gin.Context
		args2 *externalDtos.PayCartWithWalletRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args, w *externalDtos.PaymentResponse)
		want    *externalDtos.PaymentResponse
		wantErr bool
	}{
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(&models.Transaction{}, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderMerchantVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailMerchantPromotionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletRepository.On("Update", mock.Anything, mock.Anything).Return(nil, nil)
				r.cartRepository.On("Delete", mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			want: model,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderMerchantVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailMerchantPromotionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletRepository.On("Update", mock.Anything, mock.Anything).Return(nil, nil)
				r.cartRepository.On("Delete", mock.Anything).Return(&models.Cart{
					UserID: uuid.MustParse("76747c91-6b8b-42b2-bbac-d76a8e2090db"),
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderMerchantVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailMerchantPromotionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletRepository.On("Update", mock.Anything, mock.Anything).Return(nil, nil)
				r.cartRepository.On("Delete", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderMerchantVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailMerchantPromotionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletRepository.On("Update", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderMerchantVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailMerchantPromotionRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderMerchantVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderDetailRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderMerchantVoucherRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderTrackingRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, nil)
				r.orderRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionMarketplaceVoucherRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, nil)
				r.marketplaceVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, nil)
				r.walletTransactionPaymentRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 1000000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.transactionRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 0,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: falseRequestTotalPayment,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.NominalDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 1000000,
					DiscountType:       constants.NominalDiscountType,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.PercentageDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 0,
					DiscountType:       constants.NominalDiscountType,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
				r.merchantVoucherRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:         1,
					VoucherStatus:      constants.OngoingPromotionOrVoucher,
					CurrentQuota:       100,
					MaximumQuota:       120,
					MinimumSpendAmount: 1000000,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("UpsertUsingMap", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 1,
				}, nil)
				r.productMerchantPromotionRepository.On("AddCurrentQuota", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(&models.UserMerchantPromotion{
					Usage: 5,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:        1,
					PromotionStatus:  constants.OngoingPromotionOrVoucher,
					DiscountType:     constants.PercentageDiscountType,
					MaximumQuota:     10,
					MaximumUserLimit: 7,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
				r.userMerchantPromotionRepository.On("GetUserUsage", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
					DiscountType:    constants.PercentageDiscountType,
					MaximumQuota:    2,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
					DiscountType:    constants.NominalDiscountType,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				r.productItemRepository.On("SubtractProductStock", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMediaSnapshotRepository.On("CreateProductMediaSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(&models.ProductSnapshot{}, nil)
				r.productVariationSnapshotRepository.On("CreateProductVariationSnapshot", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 10,
				}, nil)
				r.productSnapshotRepository.On("CreateProductSnapshot", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID:    1,
					ProductID:     1,
					Status:        constants.ProductActiveStatus,
					StockQuantity: 2,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID: 1,
					ProductID:  1,
					Status:     constants.ProductInactiveStatus,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID: 1,
					ProductID:  2,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(&dtos.ProductItemDetails{
					MerchantID: 2,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.OngoingPromotionOrVoucher,
				}, nil)
				r.productItemRepository.On("GetDetailsById", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID:       1,
					PromotionStatus: constants.EndedPromotionOrVoucher,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.ProductMerchantPromotionDetails{
					ProductID: 2,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, nil)
				r.productMerchantPromotionRepository.On("GetDetailsByID", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(&countOrder, nil)
				r.productCourierServiceRepository.On("GetAllByProductIDAndCourierServiceID", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.orderRepository.On("GetOrderCount").Return(nil, e.ErrDataNotFound)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(&merchantVoucherID, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.merchantVoucherRepository.On("ValidateMerchantVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  100,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID:    1,
					VoucherStatus: constants.EndedPromotionOrVoucher,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(&dtos.MerchantVoucherDetails{
					MerchantID: 2,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(&models.AddressSnapshot{}, nil)
				r.merchantVoucherRepository.On("GetDetailsByID", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					ID:      uuid.UUID{},
					Balance: 100000,
					Pin:     hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
				r.addressSnapshotRepository.On("Create", mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					Pin: hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(&marketplaceVoucherID, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					Pin: hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  120,
				}, nil)
				r.marketplaceVoucherRepository.On("ValidateMarketplaceVoucherByBuyer", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					Pin: hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.OngoingPromotionOrVoucher,
					CurrentQuota:  100,
					MaximumQuota:  100,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					Pin: hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(&dtos.MarketplaceVoucherDetails{
					VoucherStatus: constants.EndedVoucherStatus,
				}, nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					Pin: hashedPIN,
				}, nil)
				r.marketplaceVoucherRepository.On("GetDetailsByID", mock.Anything, mock.Anything).Return(nil, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					Pin: falseHashedPIN,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrDataNotFound)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{
					Pin: falseHashedPIN,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("nan", nil)
				r.walletRepository.On("GetByUserId", mock.Anything).Return(&models.Wallet{}, err)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("3", nil)
			},
			wantErr: true,
		},
		{
			args: args{
				args1: c,
				args2: request,
			},
			mock: func(r *fields, a args, w *externalDtos.PaymentResponse) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", err)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				repositories.NewMockAddressSnapshotRepository(t),
				repositories.NewMockCartRepository(t),
				repositories.NewMockOrderRepository(t),
				repositories.NewMockOrderDetailRepository(t),
				repositories.NewMockOrderDetailMerchantPromotionRepository(t),
				repositories.NewMockOrderMerchantVoucherRepository(t),
				repositories.NewMockOrderTrackingRepository(t),
				repositories.NewMockProductRepository(t),
				repositories.NewMockProductItemRepository(t),
				repositories.NewMockProductCourierServiceRepository(t),
				repositories.NewMockProductMerchantPromotionRepository(t),
				repositories.NewMockProductSnapshotRepository(t),
				repositories.NewMockProductMediaSnapshotRepository(t),
				repositories.NewMockProductVariationSnapshotRepository(t),
				repositories.NewMockMerchantVoucherRepository(t),
				repositories.NewMockMerchantPromotionRepository(t),
				repositories.NewMockMarketplaceVoucherRepository(t),
				repositories.NewMockRedisRepository(t),
				repositories.NewMockTransactionRepository(t),
				repositories.NewMockTransactionMarketplaceVoucherRepository(t),
				repositories.NewMockUserMerchantPromotionRepository(t),
				repositories.NewMockWalletRepository(t),
				repositories.NewMockWalletTransactionPaymentRepository(t),
			}

			cartService := services.NewCartService(&services.CartServiceConfig{
				AddressSnapshotRepository:               tt.fields.addressSnapshotRepository,
				CartRepository:                          tt.fields.cartRepository,
				OrderRepository:                         tt.fields.orderRepository,
				OrderDetailRepository:                   tt.fields.orderDetailRepository,
				OrderDetailMerchantPromotionRepository:  tt.fields.orderDetailMerchantPromotionRepository,
				OrderMerchantVoucherRepository:          tt.fields.orderMerchantVoucherRepository,
				OrderTrackingRepository:                 tt.fields.orderTrackingRepository,
				ProductRepository:                       tt.fields.productRepository,
				ProductItemRepository:                   tt.fields.productItemRepository,
				ProductCourierServiceRepository:         tt.fields.productCourierServiceRepository,
				ProductMerchantPromotionRepository:      tt.fields.productMerchantPromotionRepository,
				ProductSnapshotRepository:               tt.fields.productSnapshotRepository,
				ProductMediaSnapshotRepository:          tt.fields.productMediaSnapshotRepository,
				ProductVariationSnapshotRepository:      tt.fields.productVariationSnapshotRepository,
				MerchantVoucherRepository:               tt.fields.merchantVoucherRepository,
				MerchantPromotionRepository:             tt.fields.merchantPromotionRepository,
				MarketplaceVoucherRepository:            tt.fields.marketplaceVoucherRepository,
				RedisRepository:                         tt.fields.redisRepository,
				TransactionRepository:                   tt.fields.transactionRepository,
				TransactionMarketplaceVoucherRepository: tt.fields.transactionMarketplaceVoucherRepository,
				UserMerchantPromotionRepository:         tt.fields.userMerchantPromotionRepository,
				WalletRepository:                        tt.fields.walletRepository,
				WalletTransactionPaymentRepository:      tt.fields.walletTransactionPaymentRepository,
			})

			tt.mock(&tt.fields, tt.args, tt.want)

			got, err := cartService.PayCartWithWallet(tt.args.args1, tt.args.args2)
			if (err != nil) != tt.wantErr {
				t.Errorf("cartService.PayCartWithWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cartService.PayCartWithWallet() = %v, want %v", got, tt.want)
			}
		})
	}
}
