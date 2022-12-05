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
	productItemService := services.NewProductItemService(&services.ProductItemServiceConfig{
		ProductItemRepository: productItemRepository,
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

	productService := services.NewProductService(&services.ProductServiceConfig{
		ProductRepository:          productRepository,
		MerchantRepository:         merchantRepository,
		ProductItemRepository:      productItemRepository,
		ProductVariationRepository: productVariationRepository,
		ProductViewRepository:      productViewRepository,
		ProductFavoriteRepository:  productFavoriteRepository,
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

	h := New(&HandlerConfig{
		RoleService:              roleService,
		WalletService:            walletService,
		AuthService:              authService,
		CategoryService:          categoryService,
		UserService:              userService,
		ProfileService:           profileService,
		SealabsPayService:        sealabsPayService,
		MerchantVoucherService:   merchantVoucherService,
		ProductService:           productService,
		ProductItemService:       productItemService,
		AddressService:           addressService,
		MerchantService:          merchantService,
		OrderService:             orderService,
		ProductReviewService:     productReviewService,
		MerchantPromotionService: merchantPromotionService,
		CartService:              cartService,
		SubDistrictService:       subDistrictService,
	})

	return h
}
