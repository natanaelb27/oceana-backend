package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
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
	newModel.Pin, err = utils.Hash(requestModel.Pin)
	if err != nil {
		return nil, e.ErrHashingFailed
	}
	return p.walletRepository.Activate(id, newModel)
}
