package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"github.com/google/uuid"
)

type WalletService interface {
	Activate(id uuid.UUID, requestModel *models.Wallet) (*models.Wallet, error)
}

type walletService struct {
	walletRepository r.WalletRepository
}

type WalletServiceConfig struct {
	WalletRepository r.WalletRepository
}

func NewWalletService(c *WalletServiceConfig) WalletService {
	return &walletService{
		walletRepository: c.WalletRepository,
	}
}

func (p *walletService) Activate(id uuid.UUID, requestModel *models.Wallet) (*models.Wallet, error) {
	newModel, err := p.walletRepository.GetByUserId(id)
	if err != nil {
		return newModel, err
	}
	newModel.IsActive = requestModel.IsActive
	return p.walletRepository.Activate(id, newModel)
}
