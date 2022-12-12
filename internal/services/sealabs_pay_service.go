package services

import (
	"bytes"
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
	Create(c *gin.Context, userID uuid.UUID, request *externalDtos.SealabsPayCreateRequest) (*externalDtos.SealabsPayCreateResponse, error)
	GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error)
	GetById(id uuid.UUID) (model *models.SealabsPay, err error)
	GetByUserIDAndCardNumber(userID uuid.UUID, cardNumber string) (*internalDtos.SealabsPayDTO, error)
	Delete(id uuid.UUID) (model *models.SealabsPay, err error)
	SetDefault(id uuid.UUID, userID uuid.UUID) error
	StoreSignature(c *gin.Context, request *externalDtos.SealabsPayCallbackRequest) error
	TopupWalletOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error)
	RegisterOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error)
	Refund(request *externalDtos.SealabsPayRefundRequest) error
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

func (p *sealabsPayService) Create(c *gin.Context, userID uuid.UUID, request *externalDtos.SealabsPayCreateRequest) (*externalDtos.SealabsPayCreateResponse, error) {
	// refund user
	refundRequest, err := externalDtos.NewSealabsPayRefundRequest(request.TxnID)
	if err != nil {
		return nil, err
	}

	err = p.Refund(refundRequest)
	if err != nil {
		return nil, err
	}

	model := request.ToModel(userID)

	if !utils.ValidateSealabsPayCardNumber(model.AccountNumber) {
		return nil, e.ErrInvalidSealabsPayCardNumber
	}

	// set current card as default if none currently exists
	_, err = p.sealabsPayRepository.GetAll(userID)
	if err == e.ErrDataNotFound {
		model.IsDefault = true
	}

	_, err = p.redisRepository.Get(c, request.Signature)
	if err != nil {
		return nil, e.ErrInvalidSignature
	}
	err = p.redisRepository.Delete(c, request.Signature)
	if err != nil {
		return nil, err
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

func (p *sealabsPayService) GetByUserIDAndCardNumber(userID uuid.UUID, cardNumber string) (*internalDtos.SealabsPayDTO, error) {
	model, err := p.sealabsPayRepository.GetByUserIDAndCardNumber(userID, cardNumber)
	if err != nil {
		return nil, err
	}

	sealabsPayDTO := internalDtos.NewSealabsPayDTO(model)

	return sealabsPayDTO, nil
}

func (p *sealabsPayService) Delete(id uuid.UUID) (model *models.SealabsPay, err error) {
	return p.sealabsPayRepository.Delete(id)
}

func (p *sealabsPayService) SetDefault(id uuid.UUID, userID uuid.UUID) error {
	// unset current default
	currentDefault, err := p.sealabsPayRepository.GetDefaultByUserID(userID)
	if err != nil {
		return err
	}

	currentDefault.IsDefault = false
	err = p.sealabsPayRepository.Update(currentDefault.ID, currentDefault)
	if err != nil {
		return err
	}

	// set new default
	newDefault, err := p.sealabsPayRepository.GetById(id)
	if err != nil {
		return e.ErrInvalidSealabsPayID
	}

	newDefault.IsDefault = true
	err = p.sealabsPayRepository.Update(id, newDefault)
	if err != nil {
		return err
	}

	return nil
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

func (p *sealabsPayService) TopupWalletOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error) {
	sealabsPayInfo, err := p.generateSealabsPayInfo(constants.SealabsPayTopup, inboundRequest)
	if err != nil {
		return nil, err
	}

	req, err := p.generatePaymentRequest(inboundRequest, sealabsPayInfo)
	if err != nil {
		return nil, err
	}

	otpURL, err := p.callSealabsPayPaymentApi(req)
	if err != nil {
		return nil, err
	}

	return otpURL, nil
}

func (p *sealabsPayService) RegisterOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error) {
	_, err := p.sealabsPayRepository.GetByCardNumber(inboundRequest.CardNumber)
	if err == nil {
		return nil, e.ErrSealabsPayCardNumberAlreadyExist
	}

	sealabsPayInfo, err := p.generateSealabsPayInfo(constants.SealabsPayRegister, inboundRequest)
	if err != nil {
		return nil, err
	}

	req, err := p.generatePaymentRequest(inboundRequest, sealabsPayInfo)
	if err != nil {
		return nil, err
	}

	otpURL, err := p.callSealabsPayPaymentApi(req)
	if err != nil {
		return nil, err
	}

	return otpURL, nil

}

func (p *sealabsPayService) Refund(refundRequest *externalDtos.SealabsPayRefundRequest) error {
	sealabsPayInfo, err := p.generateSealabsPayInfo(constants.SealabsPayRegister, &externalDtos.SealabsPayPaymentInboundRequest{})
	if err != nil {
		return err
	}

	outboundRequest, err := p.generateRefundRequest(refundRequest, sealabsPayInfo)
	if err != nil {
		return err
	}

	err = p.callSealabsPayRefundApi(outboundRequest)
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
	var redirectURL string
	switch action {
	case constants.SealabsPayRegister:
		redirectURL = fmt.Sprintf("%s/sealabs-pay/redirect/register?card_number=%s&name_on_card=%s", redirectBaseURL, inboundRequest.CardNumber, inboundRequest.NameOnCard)
	case constants.SealabsPayPayment:
		redirectURL = fmt.Sprintf("%s/sealabs-pay/redirect/payment?card_number=%s", redirectBaseURL, inboundRequest.CardNumber)
	case constants.SealabsPayTopup:
		redirectURL = fmt.Sprintf("%s/sealabs-pay/redirect/topup?card_number=%s", redirectBaseURL, inboundRequest.CardNumber)
	default:
		return nil, e.ErrInvalidSealabsPayAction
	}

	callbackURL := os.Getenv("SEALABSPAY_CALLBACK_URL")
	if len(callbackURL) == 0 {
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

func (p *sealabsPayService) generatePaymentRequest(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest, sealabsPayInfo *internalDtos.SealabsPayInfo) (*http.Request, error) {

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

func (p *sealabsPayService) generateRefundRequest(refundRequest *externalDtos.SealabsPayRefundRequest, sealabsPayInfo *internalDtos.SealabsPayInfo) (*http.Request, error) {
	data, err := json.Marshal(&refundRequest)
	if err != nil {
		return nil, e.ErrInternalOperationFailed
	}

	outboundRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/cancel", sealabsPayInfo.ApiURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, e.ErrInternalOperationFailed
	}

	outboundRequest.Header.Set("Content-Type", "application/json")
	outboundRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sealabsPayInfo.ApiKey))

	return outboundRequest, nil
}

func (p *sealabsPayService) callSealabsPayPaymentApi(request *http.Request) (*string, error) {
	client := &http.Client{
		// turn off auto-redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var errResponse *externalDtos.SealabsPayPaymentResponseError
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(responseBodyBytes, &errResponse)
	if response.StatusCode != http.StatusSeeOther {
		return nil, e.ErrSealabsPayOperationFailed
	}

	headerLocation, ok := response.Header["Location"]
	if !ok {
		return nil, e.ErrInternalOperationFailed
	}
	return &headerLocation[0], nil
}

func (p *sealabsPayService) callSealabsPayRefundApi(request *http.Request) error {
	client := &http.Client{
		// turn off auto-redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		var errResponse *externalDtos.SealabsPayPaymentResponseError
		responseBodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		json.Unmarshal(responseBodyBytes, &errResponse)
		return fmt.Errorf("%s", errResponse.Message)
	}

	return nil
}
