package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SealabsPayService interface {
	Create(id uuid.UUID, request *externalDtos.SealabsPayCreateRequest) (*externalDtos.SealabsPayCreateResponse, error)
	GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error)
	GetById(id uuid.UUID) (model *models.SealabsPay, err error)
	GetByUserIDAndAccountNumber(userID uuid.UUID, accountNumber string) (*internalDtos.SealabsPayDTO, error)
	Delete(id uuid.UUID) (model *models.SealabsPay, err error)
	StoreSignature(c *gin.Context, request *externalDtos.SealabsPayCallbackRequest) error
	TopupWallet(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error)
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

func (p *sealabsPayService) Create(id uuid.UUID, request *externalDtos.SealabsPayCreateRequest) (*externalDtos.SealabsPayCreateResponse, error) {
	model := request.ToModel(id)

	if !utils.ValidateSealabsPayAccountNumber(model.AccountNumber) {
		return nil, e.ErrInvalidSealabsPayAccountNumber
	}

	newModel, err := p.sealabsPayRepository.Create(model)

	if err != nil {
		return nil, err
	}

	response := externalDtos.NewSealabsPayCreateResponse(newModel)

	return response, nil
}

func (p *sealabsPayService) GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error) {
	return p.sealabsPayRepository.GetAll(userID)
}

func (p *sealabsPayService) GetById(id uuid.UUID) (model *models.SealabsPay, err error) {
	return p.sealabsPayRepository.GetById(id)
}

func (p *sealabsPayService) GetByUserIDAndAccountNumber(userID uuid.UUID, accountNumber string) (*internalDtos.SealabsPayDTO, error) {
	model, err := p.sealabsPayRepository.GetByUserIDAndAccountNumber(userID, accountNumber)
	if err != nil {
		return nil, err
	}

	sealabsPayDTO := internalDtos.NewSealabsPayDTO(model)

	return sealabsPayDTO, nil
}

func (p *sealabsPayService) Delete(id uuid.UUID) (model *models.SealabsPay, err error) {
	return p.sealabsPayRepository.Delete(id)
}

func (p *sealabsPayService) StoreSignature(c *gin.Context, request *externalDtos.SealabsPayCallbackRequest) error {
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

func (p *sealabsPayService) generateSealabsPayInfo(action constants.SealabsPayAction, inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*internalDtos.SealabsPayInfo, error) {
	apiURL := os.Getenv("SEALABSPAY_API_URL")
	if len(apiURL) == 0 {
		return nil, e.ErrInvalidSealabsPayAction
	}
	apiKey := os.Getenv("SEALABSPAY_API_KEY")
	if len(apiKey) == 0 {
		return nil, e.ErrInvalidSealabsPayAction
	}
	merchantCode := os.Getenv("SEALABSPAY_MERCHANT_CODE")
	if len(merchantCode) == 0 {
		return nil, e.ErrInvalidSealabsPayAction
	}
	redirectBaseURL := os.Getenv("FRONTEND_URL")
	if len(redirectBaseURL) == 0 {
		return nil, e.ErrInvalidSealabsPayAction
	}
	redirectURL := fmt.Sprintf("%s/sealabs-pay/redirect?card_number=%s", redirectBaseURL, inboundRequest.CardNumber)
	callbackBaseURL := os.Getenv("SEALABSPAY_CALLBACK_BASE_URL")
	if len(callbackBaseURL) == 0 {
		return nil, e.ErrInvalidSealabsPayAction
	}

	var callbackURL string
	switch action {
	case constants.SealabsPayRegister:
		callbackURL = fmt.Sprintf("%s/register", callbackBaseURL)
	case constants.SealabsPayPayment:
		callbackURL = fmt.Sprintf("%s/payment", callbackBaseURL)
	case constants.SealabsPayTopup:
		callbackURL = fmt.Sprintf("%s/topup", callbackBaseURL)
	default:
		return nil, e.ErrInvalidSealabsPayAction
	}

	signature, err := utils.HashPaymentInfoSHA256HMAC(inboundRequest.CardNumber, inboundRequest.Amount, merchantCode, apiKey)
	if err != nil {
		return nil, err
	}

	return &internalDtos.SealabsPayInfo{
		ApiURL:       apiURL,
		ApiKey:       apiKey,
		MerchantCode: merchantCode,
		RedirectURL:  redirectURL,
		CallbackURL:  callbackURL,
		Signature:    *signature,
	}, nil
}

func (p *sealabsPayService) generateOutboundRequest(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest, sealabsPayInfo *internalDtos.SealabsPayInfo) (*http.Request, error) {

	outboundRequestPayload := url.Values{}
	outboundRequestPayload.Set("card_number", inboundRequest.CardNumber)
	outboundRequestPayload.Set("amount", inboundRequest.Amount)
	outboundRequestPayload.Set("merchant_code", sealabsPayInfo.MerchantCode)
	outboundRequestPayload.Set("redirect_url", sealabsPayInfo.RedirectURL)
	outboundRequestPayload.Set("callback_url", sealabsPayInfo.CallbackURL)
	outboundRequestPayload.Set("signature", sealabsPayInfo.Signature)

	encodedOutboundRequest := outboundRequestPayload.Encode()
	if len(encodedOutboundRequest) == 0 {
		return nil, e.ErrSealabsPayOperationFailed
	}

	outboundRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/pay", sealabsPayInfo.ApiURL), strings.NewReader(encodedOutboundRequest))
	if err != nil {
		return nil, err
	}

	outboundRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	outboundRequest.Header.Set("Content-Length", strconv.Itoa(len(encodedOutboundRequest)))

	return outboundRequest, nil
}

func (p *sealabsPayService) handleErrorSealabsPayReidrect(response *http.Response) error {
	var errResponse *externalDtos.SealabsPayRedirectResponseError
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(responseBodyBytes, &errResponse)
	if response.StatusCode != http.StatusSeeOther {
		return e.ErrSealabsPayOperationFailed
	}

	return nil
}

func (p *sealabsPayService) retrieveOTPURL(response *http.Response) (*string, error) {
	headerLocation, ok := response.Header["Location"]
	if !ok {
		return nil, e.ErrInternalOperationFailed
	}
	return &headerLocation[0], nil
}

func (p *sealabsPayService) TopupWallet(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error) {
	client := &http.Client{
		// turn off auto-redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	sealabsPayInfo, err := p.generateSealabsPayInfo(constants.SealabsPayTopup, inboundRequest)
	if err != nil {
		return nil, err
	}

	req, err := p.generateOutboundRequest(inboundRequest, sealabsPayInfo)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	err = p.handleErrorSealabsPayReidrect(res)
	if err != nil {
		return nil, err
	}

	otpURL, err := p.retrieveOTPURL(res)
	if err != nil {
		return nil, err
	}
	return otpURL, nil
}
