package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
)

func Init() *Handler {
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
	})
	userRepository := repositories.NewUserRepository(&repositories.UserRepositoryConfig{
		DB: db.Get(),
	})
	emailRepository := repositories.NewEmailRepository(&repositories.EmailRepositoryConfig{
		DB: db.Get(),
	})
	refreshTokenRepository := repositories.NewRefreshTokenRepository(&repositories.RefreshTokenRepositoryConfig{
		DB: db.Get(),
	})
	authService := services.NewAuthService(&services.AuthServiceConfig{
		UserRepository:         userRepository,
		EmailRepository:        emailRepository,
		WalletRepository:       walletRepository,
		RefreshTokenRepository: refreshTokenRepository,
	})
	categoryRepository := repositories.NewCategoryRepository(&repositories.CategoryRepositoryConfig{
		DB: db.Get(),
	})
	categoryService := services.NewCategoryService(&services.CategoryServiceConfig{
		CategoryRepository: categoryRepository,
	})

	userService := services.NewUserService(&services.UserServiceConfig{
		UserRepository: userRepository,
	})

	merchantRepository := repositories.NewMerchantRepository(&repositories.MerchantRepositoryConfig{
		DB: db.Get(),
	})
	merchantVoucherRepository := repositories.NewMerchantVoucherRepository(&repositories.MerchantVoucherRepositoryConfig{
		DB: db.Get(),
	})
	merchantVoucherService := services.NewMerchantVoucherService(&services.MerchantVoucherServiceConfig{
		MerchantRepository:        merchantRepository,
		MerchantVoucherRepository: merchantVoucherRepository,
	})

	h := New(&HandlerConfig{
		RoleService:            roleService,
		WalletService:          walletService,
		AuthService:            authService,
		CategoryService:        categoryService,
		UserService:            userService,
		MerchantVoucherService: merchantVoucherService,
	})

	return h
}
