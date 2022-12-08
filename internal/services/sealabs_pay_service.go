package services

import (
	"os"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SealabsPayService interface {
	Create(id uuid.UUID, request *dtos.SealabsPayCreateRequest) (*dtos.SealabsPayCreateResponse, error)
	GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error)
	GetById(id uuid.UUID) (model *models.SealabsPay, err error)
	Delete(id uuid.UUID) (model *models.SealabsPay, err error)
	StoreSignature(c *gin.Context, request *dtos.SealabsPayCallbackRequest) error
}

type sealabsPayService struct {
	sealabsPayRepository r.SealabsPayRepository
	redisRepository      r.RedisRepository
}

type SealabsPayServiceConfig struct {
	SealabsPayRepository r.SealabsPayRepository
	RedisRepository      r.RedisRepository
}

func NewSealabsPayService(c *SealabsPayServiceConfig) SealabsPayService {
	return &sealabsPayService{
		sealabsPayRepository: c.SealabsPayRepository,
		redisRepository:      c.RedisRepository,
	}
}

func (p *sealabsPayService) Create(id uuid.UUID, request *dtos.SealabsPayCreateRequest) (*dtos.SealabsPayCreateResponse, error) {
	model := request.ToModel(id)

	if !utils.ValidateSealabsPayAccountNumber(model.AccountNumber) {
		return nil, e.ErrInvalidSealabsPayAccountNumber
	}

	newModel, err := p.sealabsPayRepository.Create(model)

	if err != nil {
		return nil, err
	}

	response := dtos.NewSealabsPayCreateResponse(newModel)

	return response, nil
}

func (p *sealabsPayService) GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error) {
	return p.sealabsPayRepository.GetAll(userID)
}

func (p *sealabsPayService) GetById(id uuid.UUID) (model *models.SealabsPay, err error) {
	return p.sealabsPayRepository.GetById(id)
}

func (p *sealabsPayService) Delete(id uuid.UUID) (model *models.SealabsPay, err error) {
	return p.sealabsPayRepository.Delete(id)
}

func (p *sealabsPayService) StoreSignature(c *gin.Context, request *dtos.SealabsPayCallbackRequest) error {
	signatureMinutesDuration, err := strconv.Atoi(os.Getenv("SEALABSPAY_TOPUP_SIGNATURE_MINUTES_DURATION"))
	if err != nil {
		return err
	}
	signatureDuration := time.Duration(signatureMinutesDuration) * time.Minute
	
	err = p.redisRepository.Set(c, request.Signature, request.Signature, signatureDuration)
	if err != nil {
		return err
	}
	
	return nil
}
