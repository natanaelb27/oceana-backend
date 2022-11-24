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
	productRepository := repositories.NewProductRepository(&repositories.ProductRepositoryConfig{
		DB: db.Get(),
	})
	productService := services.NewProductService(&services.ProductServiceConfig{
		ProductRepository: productRepository,
	})

	h := New(&HandlerConfig{
		RoleService:            roleService,
		WalletService:          walletService,
		AuthService:            authService,
		CategoryService:        categoryService,
		UserService:            userService,
		ProfileService:         profileService,
		SealabsPayService:      sealabsPayService,
		MerchantVoucherService: merchantVoucherService,
		AddressService:         addressService,
		MerchantService:        merchantService,
		ProductService:         productService,
		OrderService:           orderService,
		ProductReviewService:   productReviewService,
	})

	return h
}
