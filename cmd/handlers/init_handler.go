package handlers

import (
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
	roleService := services.NewRoleService(&services.RoleServiceConfig{
		RoleRepository: roleRepository,
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
	walletService := services.NewWalletService(&services.WalletServiceConfig{
		WalletRepository: walletRepository,
		RedisRepository:  redisRepository,
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
	authService := services.NewAuthService(&services.AuthServiceConfig{
		UserRepository:   userRepository,
		EmailRepository:  emailRepository,
		WalletRepository: walletRepository,
		RedisRepository:  redisRepository,
	})
	categoryRepository := repositories.NewCategoryRepository(&repositories.CategoryRepositoryConfig{
		DB: db.Get(),
	})
	categoryService := services.NewCategoryService(&services.CategoryServiceConfig{
		CategoryRepository: categoryRepository,
	})
	merchantRepository := repositories.NewMerchantRepository(&repositories.MerchantRepositoryConfig{
		DB: db.Get(),
	})
	addressService := services.NewAddressService(&services.AddressServiceConfig{
		AddressRepository:  addressRepository,
		MerchantRepository: merchantRepository,
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
	productReviewService := services.NewProductReviewService(&services.ProductReviewServiceConfig{
		ProductReviewRepository:      productReviewRepository,
		ProductReviewMediaRepository: productReviewMediaRepository,
	})
	orderService := services.NewOrderService(&services.OrderServiceConfig{
		OrderRepository:           orderRepository,
		OrderDetailRepository:     orderDetailRepository,
		AddressSnapshotRepository: addressSnapshotRepository,
		OrderTrackingRepository:   orderTrackingRepository,
	})

	userService := services.NewUserService(&services.UserServiceConfig{
		UserRepository:  userRepository,
		EmailRepository: emailRepository,
	})
	sealabsPayRepository := repositories.NewSealabsPayRepository(&repositories.SealabsPayRepositoryConfig{
		DB: db.Get(),
	})
	sealabsPayService := services.NewSealabsPayService(&services.SealabsPayServiceConfig{
		SealabsPayRepository: sealabsPayRepository,
	})

	merchantVoucherRepository := repositories.NewMerchantVoucherRepository(&repositories.MerchantVoucherRepositoryConfig{
		DB: db.Get(),
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

	cartRepository := repositories.NewCartRepository(&repositories.CartRepositoryConfig{
		DB: db.Get(),
	})
	cartService := services.NewCartService(&services.CartServiceConfig{
		CartRepository:        cartRepository,
		ProductItemRepository: productItemRepository,
	})

	courierServiceRepository := repositories.NewCourierServiceRepository(&repositories.CourierServiceRepositoryConfig{
		DB: db.Get(),
	})
	courierRepository := repositories.NewCourierRepository(&repositories.CourierRepositoryConfig{
		DB: db.Get(),
	})
	courierService := services.NewCourierService(&services.CourierServiceConfig{
		CourierRepository:        courierRepository,
		CourierServiceRepository: courierServiceRepository,
	})

	merchantConfigurationRepository := repositories.NewMerchantConfigurationRepository(&repositories.MerchantConfigurationRepositoryConfig{
		DB: db.Get(),
	})
	merchantCourierServiceRepository := repositories.NewMerchantCourierServiceRepository(&repositories.MerchantCourierServiceRepositoryConfig{
		DB: db.Get(),
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

	productService := services.NewProductService(&services.ProductServiceConfig{
		ProductRepository:                productRepository,
		MerchantRepository:               merchantRepository,
		ProductItemRepository:            productItemRepository,
		ProductVariationRepository:       productVariationRepository,
		ProductViewRepository:            productViewRepository,
		ProductFavoriteRepository:        productFavoriteRepository,
		CategoryRepository:               categoryRepository,
		MerchantCourierServiceRepository: merchantCourierServiceRepository,
		ProductVariantTypeRepository:     productVariantTypeRepository,
		ProductVariantMediaRepository:    productVariantMediaRepository,
		ProductVariantRepository:         productVariantRepository,
		ProductCourierServiceRepository:  productCourierServiceRepository,
	})

	merchantPromotionRepository := repositories.NewMerchantPromotionRepository(&repositories.MerchantPromotionRepositoryConfig{
		DB: db.Get(),
	})
	productMerchantPromotionRepository := repositories.NewProductMerchantPromotionRepositoryConfig(&repositories.ProductMerchantPromotionRepositoryConfig{
		DB: db.Get(),
	})
	merchantPromotionService := services.NewMerchantPromotionService(&services.MerchantPromotionServiceConfig{
		MerchantRepository:                 merchantRepository,
		MerchantPromotionRepository:        merchantPromotionRepository,
		ProductMerchantPromotionRepository: productMerchantPromotionRepository,
	})
	provinceRepository := repositories.NewProvinceRepository(&repositories.ProvinceRepositoryConfig{
		DB: db.Get(),
	})
	citiesRepository := repositories.NewCityRepository(&repositories.CityRepositoryConfig{
		DB: db.Get(),
	})
	districtRepository := repositories.NewDistrictRepository(&repositories.DistrictRepositoryConfig{
		DB: db.Get(),
	})
	subDistrictRepository := repositories.NewSubDistrictRepository(&repositories.SubDistrictRepositoryConfig{
		DB: db.Get(),
	})
	subDistrictService := services.NewSubDistrictService(&services.SubDistrictServiceConfig{
		ProvinceRepository:    provinceRepository,
		CityRepository:        citiesRepository,
		DistrictRepository:    districtRepository,
		SubDistrictRepository: subDistrictRepository,
	})

	marketplaceVoucherRepository := repositories.NewMarketplaceVoucherRepository(&repositories.MarketplaceVoucherRepositoryConfig{
		DB: db.Get(),
	})
	marketplaceVoucherService := services.NewMarketplaceVoucherService(&services.MarketplaceVoucherServiceConfig{
		MarketplaceVoucherRepository: marketplaceVoucherRepository,
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
		MarketplaceVoucherService: marketplaceVoucherService,
		CourierService:            courierService,
		LogisticService:           logisticService,
	})

	return h
}
