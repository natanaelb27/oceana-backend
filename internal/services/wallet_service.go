package services

import (
	"net/http"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletService interface {
	Activate(id uuid.UUID, requestModel *models.Wallet) (*models.Wallet, error)
	ChangePin(c *gin.Context, userId uuid.UUID, request *externalDtos.WalletChangePinRequest) (*models.Wallet, error)
	CreateWalletTopUp(c *gin.Context, userID uuid.UUID, request *externalDtos.TopUpWalletRequest) error
	GetWalletDetails(id uuid.UUID) (*externalDtos.WalletDetailsResponse, error)
}

type walletService struct {
	walletRepository      r.WalletRepository
	redisRepository       r.RedisRepository
	walletTopUpRepository r.WalletTopUpRepository
	sealabsPayRepository  r.SealabsPayRepository
}

type WalletServiceConfig struct {
	WalletRepository      r.WalletRepository
	RedisRepository       r.RedisRepository
	WalletTopUpRepository r.WalletTopUpRepository
	SealabsPayRepository  r.SealabsPayRepository
}

func NewWalletService(c *WalletServiceConfig) WalletService {
	return &walletService{
		walletRepository:      c.WalletRepository,
		redisRepository:       c.RedisRepository,
		walletTopUpRepository: c.WalletTopUpRepository,
		sealabsPayRepository:  c.SealabsPayRepository,
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

func (p *walletService) ChangePin(c *gin.Context, userId uuid.UUID, request *externalDtos.WalletChangePinRequest) (wallet *models.Wallet, err error) {
	redisKey := "upgrade-token-password:" + userId.String()

	token, err := p.redisRepository.Get(c, redisKey)
	if err != nil {
		return nil, e.ErrStepUpNotFound
	}

	if token != request.UpgradeToken {
		return nil, e.ErrStepUpNotFound
	}

	jwtWrapper, err := utils.NewJwtWrapper()
	if err != nil {
		return nil, err
	}
	claims, err := jwtWrapper.JwtStepUpValidateToken(token)
	if err != nil {
		utils.Response(c, http.StatusUnauthorized, nil, nil, http.StatusText(http.StatusUnauthorized), err.Error())
		c.Abort()
		return
	}
	if claims.Data.Scope != "level1" {
		err = e.ErrStepUpScopeIncorrect
		return
	}

	newModel, err := p.walletRepository.GetByUserId(userId)
	if err != nil {
		return newModel, err
	}
	if !newModel.IsActive {
		return nil, e.ErrWalletIsNotActive
	}

	newModel.Pin, err = utils.Hash(request.Pin)
	if err != nil {
		return nil, e.ErrHashingFailed
	}

	wallet, err = p.walletRepository.Update(newModel.ID, newModel)
	if err != nil {
		return
	}

	err = p.redisRepository.Delete(c, redisKey)
	if err != nil {
		return
	}

	return
}

func (p *walletService) GetWalletDetails(id uuid.UUID) (*externalDtos.WalletDetailsResponse, error) {
	wallet, err := p.walletRepository.GetByUserId(id)
	if err != nil {
		return nil, err
	}
	walletTransactionHistory, err := p.walletRepository.GetWalletTransactionHistory(wallet.ID)
	if err != nil {
		return nil, err
	}
	responseData := externalDtos.NewWalletDetailsResponse(wallet, walletTransactionHistory)
	return responseData, nil
}

func (p *walletService) CreateWalletTopUp(c *gin.Context, userID uuid.UUID, request *externalDtos.TopUpWalletRequest) error {
	// validate card number length
	if request.CardNumber == "" || len(request.CardNumber) != constants.AccountNumberMaxLength {
		return e.ErrInvalidSealabsPayAccountNumber
	}

	// check if signature matches with the one in redis
	_, err := p.redisRepository.Get(c, request.Signature)
	if err != nil {
		return e.ErrInvalidSignature
	}
	// clean up signature from redis
	err = p.redisRepository.Delete(c, request.Signature)
	if err != nil {
		return err
	}

	// get wallet_id by user id
	wallet, err := p.walletRepository.GetByUserId(userID)
	if err != nil {
		return err
	}

	// get sealabs_pay_id by account number
	sealabsPayAccount, err := p.sealabsPayRepository.GetByUserIDAndAccountNumber(userID, request.CardNumber)
	if err != nil {
		return err
	}

	// create wallet_topup entry
	walletTopup, err := request.ToModel()
	if err != nil {
		return err
	}
	walletTopup.WalletID = wallet.ID
	walletTopup.SealabsPayID = sealabsPayAccount.ID
	_, err = p.walletTopUpRepository.Create(walletTopup)
	if err != nil {
		return err
	}

	// update wallet balance
	wallet.Balance += walletTopup.Amount
	_, err = p.walletRepository.Update(wallet.ID, wallet)
	if err != nil {
		return err
	}

	return nil
}
