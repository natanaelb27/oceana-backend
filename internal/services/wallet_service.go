package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletService interface {
	Activate(id uuid.UUID, requestModel *models.Wallet) (*models.Wallet, error)
	ChangePin(c *gin.Context, userId uuid.UUID, requestModel *models.Wallet) (*models.Wallet, error)
}

type walletService struct {
	walletRepository r.WalletRepository
	redisRepository  r.RedisRepository
}

type WalletServiceConfig struct {
	WalletRepository r.WalletRepository
	RedisRepository  r.RedisRepository
}

func NewWalletService(c *WalletServiceConfig) WalletService {
	return &walletService{
		walletRepository: c.WalletRepository,
		redisRepository:  c.RedisRepository,
	}
}

func (p *walletService) Activate(id uuid.UUID, requestModel *models.Wallet) (*models.Wallet, error) {
	newModel, err := p.walletRepository.GetByUserId(id)
	if err != nil {
		return newModel, err
	}
	if newModel.IsActive {
		return nil, e.ErrWalletAlreadyActive
	}
	newModel.IsActive = true
	newModel.Pin, err = utils.Hash(requestModel.Pin)
	if err != nil {
		return nil, e.ErrHashingFailed
	}
	return p.walletRepository.Activate(newModel.ID, newModel)
}

func (p *walletService) ChangePin(c *gin.Context, userId uuid.UUID, requestModel *models.Wallet) (*models.Wallet, error) {
	_, err := p.redisRepository.Get(c, "step-up-token:"+userId.String()+":"+c.ClientIP())
	if err != nil {
		return nil, e.ErrStepUpNotFound
	}

	// Do Blacklist (delete stepup token) here

	newModel, err := p.walletRepository.GetByUserId(userId)
	if err != nil {
		return newModel, err
	}
	if !newModel.IsActive {
		return nil, e.ErrWalletIsNotActive
	}

	newModel.Pin, err = utils.Hash(requestModel.Pin)
	if err != nil {
		return nil, e.ErrHashingFailed
	}

	return p.walletRepository.Update(newModel.ID, newModel)
}
