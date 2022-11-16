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
	sealabsPayRepository := repositories.NewSealabsPayRepository(&repositories.SealabsPayRepositoryConfig{
		DB: db.Get(),
	})
	sealabsPayService := services.NewSealabsPayService(&services.SealabsPayServiceConfig{
		SealabsPayRepository: sealabsPayRepository,
	})

	h := New(&HandlerConfig{
		RoleService:     roleService,
		WalletService:   walletService,
		AuthService:     authService,
		CategoryService: categoryService,
		SealabsPay:      sealabsPayService,
	})

	return h
}
