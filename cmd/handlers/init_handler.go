package handlers

import (
	"net/http"

	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
)

func Init() *Handler {
	redisRepository := repositories.NewRedisRepository(&repositories.RedisRepositoryConfig{
		RedisClient: db.RedisGet(),
	})

	roleRepository := repositories.NewRoleRepository(&repositories.RoleRepositoryConfig{
		DB: db.Get(),
	})

	productRepository := repositories.NewProductRepository(&repositories.ProductRepositoryConfig{
		DB: db.Get(),
	})

	productItemRepository := repositories.NewProductItemRepository(&repositories.ProductItemRepositoryConfig{
		DB: db.Get(),
	})

	productVariationRepository := repositories.NewProductVariationRepository(&repositories.ProductVariationRepositoryConfig{
		DB: db.Get(),
	})

	productViewRepository := repositories.NewProductViewRepository(&repositories.ProductViewRepositoryConfig{
		DB: db.Get(),
	})

	productFavoriteRepository := repositories.NewProductFavoriteRepository(&repositories.ProductFavoriteRepositoryConfig{
		DB: db.Get(),
	})

	walletRepository := repositories.NewWalletRepository(&repositories.WalletRepositoryConfig{
		DB: db.Get(),
	})
	walletTopupRepository := repositories.NewWalletTopupRepository(&repositories.WalletTopupRepositoryConfig{
		DB: db.Get(),
	})
	sealabsPayRepository := repositories.NewSealabsPayRepository(&repositories.SealabsPayRepositoryConfig{
		DB: db.Get(),
	})
	walletService := services.NewWalletService(&services.WalletServiceConfig{
		WalletRepository:      walletRepository,
		RedisRepository:       redisRepository,
		WalletTopupRepository: walletTopupRepository,
		SealabsPayRepository:  sealabsPayRepository,
	})

	userRepository := repositories.NewUserRepository(&repositories.UserRepositoryConfig{
		DB: db.Get(),
	})

	emailRepository := repositories.NewEmailRepository(&repositories.EmailRepositoryConfig{
		DB: db.Get(),
	})

	addressRepository := repositories.NewAddressRepository(&repositories.AddressRepositoryConfig{
		DB: db.Get(),
	})

	categoryRepository := repositories.NewCategoryRepository(&repositories.CategoryRepositoryConfig{
		DB: db.Get(),
	})

	merchantRepository := repositories.NewMerchantRepository(&repositories.MerchantRepositoryConfig{
		DB: db.Get(),
	})

	orderRepository := repositories.NewOrderRepository(&repositories.OrderRepositoryConfig{
		DB: db.Get(),
	})

	orderDetailRepository := repositories.NewOrderDetailRepository(&repositories.OrderDetailRepositoryConfig{
		DB: db.Get(),
	})

	addressSnapshotRepository := repositories.NewAddressSnapshotRepository(&repositories.AddressSnapshotRepositoryConfig{
		DB: db.Get(),
	})

	orderTrackingRepository := repositories.NewOrderTrackingRepository(&repositories.OrderTrackingRepositoryConfig{
		DB: db.Get(),
	})

	productReviewRepository := repositories.NewProductReviewRepository(&repositories.ProductReviewRepositoryConfig{
		DB: db.Get(),
	})

	productReviewMediaRepository := repositories.NewProductReviewMediaRepository(&repositories.ProductReviewMediaRepositoryConfig{
		DB: db.Get(),
	})

	userService := services.NewUserService(&services.UserServiceConfig{
		UserRepository:  userRepository,
		EmailRepository: emailRepository,
	})
	sealabsPayService := services.NewSealabsPayService(&services.SealabsPayServiceConfig{
		SealabsPayRepository: sealabsPayRepository,
		RedisRepository:      redisRepository,
	})
	merchantVoucherRepository := repositories.NewMerchantVoucherRepository(&repositories.MerchantVoucherRepositoryConfig{
		DB: db.Get(),
	})

	cartRepository := repositories.NewCartRepository(&repositories.CartRepositoryConfig{
		DB: db.Get(),
	})

	courierServiceRepository := repositories.NewCourierServiceRepository(&repositories.CourierServiceRepositoryConfig{
		DB: db.Get(),
	})

	courierRepository := repositories.NewCourierRepository(&repositories.CourierRepositoryConfig{
		DB: db.Get(),
	})

	merchantConfigurationRepository := repositories.NewMerchantConfigurationRepository(&repositories.MerchantConfigurationRepositoryConfig{
		DB: db.Get(),
	})

	merchantCourierServiceRepository := repositories.NewMerchantCourierServiceRepository(&repositories.MerchantCourierServiceRepositoryConfig{
		DB: db.Get(),
	})

	merchantPromotionRepository := repositories.NewMerchantPromotionRepository(&repositories.MerchantPromotionRepositoryConfig{
		DB: db.Get(),
	})

	productMerchantPromotionRepository := repositories.NewProductMerchantPromotionRepositoryConfig(&repositories.ProductMerchantPromotionRepositoryConfig{
		DB: db.Get(),
	})

	provinceRepository := repositories.NewProvinceRepository(&repositories.ProvinceRepositoryConfig{
		DB: db.Get(),
	})

	cityRepository := repositories.NewCityRepository(&repositories.CityRepositoryConfig{
		DB: db.Get(),
	})

	districtRepository := repositories.NewDistrictRepository(&repositories.DistrictRepositoryConfig{
		DB: db.Get(),
	})

	subDistrictRepository := repositories.NewSubDistrictRepository(&repositories.SubDistrictRepositoryConfig{
		DB: db.Get(),
	})

	rajaOngkirRepository := repositories.NewRajaOngkirRepository(&repositories.RajaOngkirRepositoryConfig{
		Client: http.Client{},
	})

	marketplaceVoucherRepository := repositories.NewMarketplaceVoucherRepository(&repositories.MarketplaceVoucherRepositoryConfig{
		DB: db.Get(),
	})

	productSnapshotRepository := repositories.NewProductSnapshotRepository(&repositories.ProductSnapshotRepositoryConfig{
		DB: db.Get(),
	})

	productMediaSnapshotRepository := repositories.NewProductMediaSnapshotRepository(&repositories.ProductMediaSnapshotRepositoryConfig{
		DB: db.Get(),
	})

	productVariationSnapshotRepository := repositories.NewProductVariationSnapshotRepository(&repositories.ProductVariationSnapshotRepositoryConfig{
		DB: db.Get(),
	})

	userMerchantPromotionRepository := repositories.NewUserMerchantPromotionRepository(&repositories.UserMerchantPromotionRepositoryConfig{
		DB: db.Get(),
	})

	walletTransactionPaymentRepository := repositories.NewWalletTransactionPaymentRepository(&repositories.WalletTransactionPaymentRepositoryConfig{
		DB: db.Get(),
	})

	transactionRepository := repositories.NewTransactionRepository(&repositories.TransactionRepositoryConfig{
		DB: db.Get(),
	})

	transactionMarketplaceVoucherRepository := repositories.NewTransactionMarketplaceVoucherRepository(&repositories.TransactionMarketplaceVoucherRepositoryConfig{
		DB: db.Get(),
	})

	orderDetailMerchantPromotionRepository := repositories.NewOrderDetailMerchantPromotionRepository(&repositories.OrderDetailMerchantPromotionRepositoryConfig{
		DB: db.Get(),
	})

	orderMerchantVoucherRepository := repositories.NewOrderMerchantVoucherRepository(&repositories.OrderMerchantVoucherRepositoryConfig{
		DB: db.Get(),
	})

	sealabsPayTransactionPaymentRepository := repositories.NewSealabsPayTransactionPaymentRepository(&repositories.SealabsPayTransactionPaymentRepositoryConfig{
		DB: db.Get(),
	})

	refundRepository := repositories.NewRefundRepository(&repositories.RefundRepositoryConfig{
		DB: db.Get(),
	})

	refundDetailRepository := repositories.NewRefundDetailRepository(&repositories.RefundDetailRepositoryConfig{
		DB: db.Get(),
	})

	refundMediaRepository := repositories.NewRefundMediaRepository(&repositories.RefundMediaRepositoryConfig{
		DB: db.Get(),
	})

	productMediaRepository := repositories.NewProductMediaRepository(&repositories.ProductMediaRepositoryConfig{
		DB: db.Get(),
	})

	productVariantTypeRepository := repositories.NewProductVariantTypeRepository(&repositories.ProductVariantTypeRepositoryConfig{
		DB: db.Get(),
	})

	productVariantMediaRepository := repositories.NewProductVariantMediaRepository(&repositories.ProductVariantMediaRepositoryConfig{
		DB: db.Get(),
	})

	productVariantRepository := repositories.NewProductVariantRepository(&repositories.ProductVariantRepositoryConfig{
		DB: db.Get(),
	})

	productCourierServiceRepository := repositories.NewProductCourierServiceRepository(&repositories.ProductCourierServiceRepositoryConfig{
		DB: db.Get(),
	})

	roleService := services.NewRoleService(&services.RoleServiceConfig{
		RoleRepository: roleRepository,
	})

	authService := services.NewAuthService(&services.AuthServiceConfig{
		UserRepository:   userRepository,
		EmailRepository:  emailRepository,
		WalletRepository: walletRepository,
		RedisRepository:  redisRepository,
	})

	categoryService := services.NewCategoryService(&services.CategoryServiceConfig{
		CategoryRepository: categoryRepository,
	})

	addressService := services.NewAddressService(&services.AddressServiceConfig{
		AddressRepository:  addressRepository,
		MerchantRepository: merchantRepository,
	})

	productReviewService := services.NewProductReviewService(&services.ProductReviewServiceConfig{
		OrderDetailRepository:        orderDetailRepository,
		ProductReviewRepository:      productReviewRepository,
		ProductReviewMediaRepository: productReviewMediaRepository,
	})

	orderService := services.NewOrderService(&services.OrderServiceConfig{
		OrderRepository:                        orderRepository,
		OrderDetailRepository:                  orderDetailRepository,
		AddressSnapshotRepository:              addressSnapshotRepository,
		OrderTrackingRepository:                orderTrackingRepository,
		MerchantRepository:                     merchantRepository,
		RefundRepository:                       refundRepository,
		RefundDetailRepository:                 refundDetailRepository,
		TransactionRepository:                  transactionRepository,
		WalletTransactionPaymentRepository:     walletTransactionPaymentRepository,
		WalletRepository:                       walletRepository,
		SealabsPayTransactionPaymentRepository: sealabsPayTransactionPaymentRepository,
	})

	merchantVoucherService := services.NewMerchantVoucherService(&services.MerchantVoucherServiceConfig{
		MerchantRepository:        merchantRepository,
		MerchantVoucherRepository: merchantVoucherRepository,
	})

	merchantService := services.NewMerchantService(&services.MerchantServiceConfig{
		MerchantRepository: merchantRepository,
		UserRepository:     userRepository,
	})

	profileService := services.NewProfileService(&services.ProfileServiceConfig{
		UserRepository:  userRepository,
		RedisRepository: redisRepository,
	})

	courierService := services.NewCourierService(&services.CourierServiceConfig{
		CourierRepository:        courierRepository,
		CourierServiceRepository: courierServiceRepository,
	})

	logisticService := services.NewLogisticService(&services.LogisticServiceConfig{
		CourierRepository:                courierRepository,
		CourierServiceRepository:         courierServiceRepository,
		MerchantCourierServiceRepository: merchantCourierServiceRepository,
		UserRepository:                   userRepository,
		MerchantRepository:               merchantRepository,
		MerchantConfigurationRepository:  merchantConfigurationRepository,
		ProductRepository:                productRepository,
	})

	productService := services.NewProductService(&services.ProductServiceConfig{
		ProductRepository:                productRepository,
		MerchantRepository:               merchantRepository,
		ProductItemRepository:            productItemRepository,
		ProductReviewRepository:          productReviewRepository,
		ProductVariationRepository:       productVariationRepository,
		ProductViewRepository:            productViewRepository,
		ProductFavoriteRepository:        productFavoriteRepository,
		CategoryRepository:               categoryRepository,
		MerchantCourierServiceRepository: merchantCourierServiceRepository,
		ProductVariantTypeRepository:     productVariantTypeRepository,
		ProductVariantMediaRepository:    productVariantMediaRepository,
		ProductVariantRepository:         productVariantRepository,
		ProductCourierServiceRepository:  productCourierServiceRepository,
		ProductMediaRepository:           productMediaRepository,
	})

	merchantPromotionService := services.NewMerchantPromotionService(&services.MerchantPromotionServiceConfig{
		MerchantRepository:                 merchantRepository,
		MerchantPromotionRepository:        merchantPromotionRepository,
		ProductMerchantPromotionRepository: productMerchantPromotionRepository,
	})

	subDistrictService := services.NewSubDistrictService(&services.SubDistrictServiceConfig{
		ProvinceRepository:    provinceRepository,
		CityRepository:        cityRepository,
		DistrictRepository:    districtRepository,
		SubDistrictRepository: subDistrictRepository,
	})

	rajaOngkirService := services.NewRajaOngkirService(&services.RajaOngkirServiceConfig{
		RajaOngkirRepository: rajaOngkirRepository,
	})

	marketplaceVoucherService := services.NewMarketplaceVoucherService(&services.MarketplaceVoucherServiceConfig{
		MarketplaceVoucherRepository: marketplaceVoucherRepository,
	})

	cartService := services.NewCartService(&services.CartServiceConfig{
		AddressSnapshotRepository:               addressSnapshotRepository,
		CartRepository:                          cartRepository,
		OrderRepository:                         orderRepository,
		OrderDetailRepository:                   orderDetailRepository,
		OrderDetailMerchantPromotionRepository:  orderDetailMerchantPromotionRepository,
		OrderMerchantVoucherRepository:          orderMerchantVoucherRepository,
		OrderTrackingRepository:                 orderTrackingRepository,
		ProductRepository:                       productRepository,
		ProductItemRepository:                   productItemRepository,
		ProductMerchantPromotionRepository:      productMerchantPromotionRepository,
		ProductSnapshotRepository:               productSnapshotRepository,
		ProductMediaSnapshotRepository:          productMediaSnapshotRepository,
		ProductVariationSnapshotRepository:      productVariationSnapshotRepository,
		MerchantVoucherRepository:               merchantVoucherRepository,
		MerchantPromotionRepository:             merchantPromotionRepository,
		MarketplaceVoucherRepository:            marketplaceVoucherRepository,
		RedisRepository:                         redisRepository,
		TransactionRepository:                   transactionRepository,
		TransactionMarketplaceVoucherRepository: transactionMarketplaceVoucherRepository,
		UserMerchantPromotionRepository:         userMerchantPromotionRepository,
		WalletRepository:                        walletRepository,
		WalletTransactionPaymentRepository:      walletTransactionPaymentRepository,
	})

	transactionService := services.NewTransactionService(&services.TransactionServiceConfig{
		TransactionRepository:                  transactionRepository,
		SealabsPayTransactionPaymentRepository: sealabsPayTransactionPaymentRepository,
	})

	cityService := services.NewCityService(&services.CityServiceConfig{
		CityRepository: cityRepository,
	})

	refundService := services.NewRefundService(&services.RefundServiceConfig{
		MerchantRepository:     merchantRepository,
		OrderRepository:        orderRepository,
		OrderDetailRepository:  orderDetailRepository,
		RefundRepository:       refundRepository,
		RefundDetailRepository: refundDetailRepository,
		RefundMediaRepository:  refundMediaRepository,
		UserRepository:         userRepository,
	})

	h := New(&HandlerConfig{
		RoleService:               roleService,
		WalletService:             walletService,
		AuthService:               authService,
		CategoryService:           categoryService,
		UserService:               userService,
		ProfileService:            profileService,
		SealabsPayService:         sealabsPayService,
		MerchantVoucherService:    merchantVoucherService,
		ProductService:            productService,
		AddressService:            addressService,
		MerchantService:           merchantService,
		OrderService:              orderService,
		ProductReviewService:      productReviewService,
		MerchantPromotionService:  merchantPromotionService,
		CartService:               cartService,
		SubDistrictService:        subDistrictService,
		RajaOngkirService:         rajaOngkirService,
		MarketplaceVoucherService: marketplaceVoucherService,
		CourierService:            courierService,
		LogisticService:           logisticService,
		TransactionService:        transactionService,
		CityService:               cityService,
		RefundService:             refundService,
	})

	return h
}
