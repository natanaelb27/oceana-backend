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
	"gorm.io/gorm"
)

type SealabsPayService interface {
	WithTx(tx *gorm.DB) *sealabsPayService
	Create(c *gin.Context, userID uuid.UUID, request *externalDtos.SealabsPayCreateRequest) (*externalDtos.SealabsPayCreateResponse, error)
	GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error)
	GetById(id uuid.UUID) (model *models.SealabsPay, err error)
	GetByUserIDAndCardNumber(userID uuid.UUID, cardNumber string) (*internalDtos.SealabsPayDTO, error)
	Delete(id uuid.UUID) (*models.SealabsPay, error)
	SetDefault(id uuid.UUID, userID uuid.UUID) error
	StoreSignature(c *gin.Context, request *externalDtos.SealabsPayCallbackRequest) error
	TopupWalletOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error)
	RegisterOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error)
	PaymentOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error)
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

func (u sealabsPayService) WithTx(tx *gorm.DB) *sealabsPayService {
	u.sealabsPayRepository = u.sealabsPayRepository.WithTx(tx)

	return &u
}

func (p *sealabsPayService) Create(c *gin.Context, userID uuid.UUID, request *externalDtos.SealabsPayCreateRequest) (*externalDtos.SealabsPayCreateResponse, error) {
	fmt.Println("test 5")
	// refund user
	refundRequest, err := externalDtos.NewSealabsPayRefundRequest(request.TxnID)
	if err != nil {
		fmt.Println("SERVICE", err)
		return nil, err
	}
	fmt.Println("test 6")

	err = p.Refund(refundRequest)
	if err != nil {
		fmt.Println("SERVICE", err)
		return nil, err
	}

	fmt.Println("test 7")
	model := request.ToModel(userID)

	if !utils.ValidateSealabsPayCardNumber(model.AccountNumber) {
		fmt.Println("SERVICE", e.ErrInvalidSealabsPayCardNumber)
		return nil, e.ErrInvalidSealabsPayCardNumber
	}
	fmt.Println("test 8")

	// set current card as default if none currently exists
	_, err = p.sealabsPayRepository.GetAll(userID)
	if err == e.ErrDataNotFound {
		model.IsDefault = true
	}

	fmt.Println("test 9")
	_, err = p.redisRepository.Get(c, request.Signature)
	if err != nil {
		fmt.Println("SERVICE", err)
		return nil, e.ErrInvalidSignature
	}
	fmt.Println("test 10")
	err = p.redisRepository.Delete(c, request.Signature)
	if err != nil {
		fmt.Println("SERVICE", err)
		return nil, err
	}

	fmt.Println("test 11")
	newModel, err := p.sealabsPayRepository.Create(model)
	if err != nil {
		fmt.Println("SERVICE", err)
		return nil, err
	}

	fmt.Println("test 12")
	response := externalDtos.NewSealabsPayCreateResponse(newModel)

	fmt.Println("SERVICE SUCCESS")
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

func (p *sealabsPayService) Delete(id uuid.UUID) (*models.SealabsPay, error) {
	model, err := p.sealabsPayRepository.Delete(id)
	if err != nil {
		return nil, err
	}
	if !model.IsDefault {
		return model, nil
	}

	// if default was deleted, set the first in GetAll() as default
	models, _ := p.sealabsPayRepository.GetAll(model.UserID)
	if len(models) == 0 {
		return model, nil
	}
	firstCard := models[0]
	firstCard.IsDefault = true
	err = p.sealabsPayRepository.Update(firstCard.ID, firstCard)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (p *sealabsPayService) SetDefault(id uuid.UUID, userID uuid.UUID) error {
	// unset current default
	currentDefault, _ := p.sealabsPayRepository.GetDefaultByUserID(userID)
	if currentDefault.ID != uuid.Nil {
		err := p.sealabsPayRepository.SetDefaultFalse(currentDefault.ID, currentDefault)
		if err != nil {
			return err
		}
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
	signatureMinutesDuration, err := strconv.Atoi(os.Getenv("SEALABSPAY_SIGNATURE_MINUTES_DURATION"))
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

func (p *sealabsPayService) PaymentOTP(inboundRequest *externalDtos.SealabsPayPaymentInboundRequest) (*string, error) {
	sealabsPayInfo, err := p.generateSealabsPayInfo(constants.SealabsPayPayment, inboundRequest)
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
	fmt.Println("refund 1")
	sealabsPayInfo, err := p.generateSealabsPayInfo(constants.SealabsPayRegister, &externalDtos.SealabsPayPaymentInboundRequest{})
	if err != nil {
		return err
	}

	fmt.Println("refund 2")
	outboundRequest, err := p.generateRefundRequest(refundRequest, sealabsPayInfo)
	if err != nil {
		return err
	}

	fmt.Println("refund 3")
	err = p.callSealabsPayRefundApi(outboundRequest)
	fmt.Println(err)
	if err != nil {
		return err
	}

	fmt.Println("refund 4")
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
		fmt.Println(errResponse)
		return nil, fmt.Errorf("%v", errResponse.Message)
	}

	headerLocation, ok := response.Header["Location"]
	if !ok {
		return nil, e.ErrInternalOperationFailed
	}
	return &headerLocation[0], nil
}

func (p *sealabsPayService) callSealabsPayRefundApi(request *http.Request) error {
	fmt.Println("callSealabsPayRefundApi 1")
	client := &http.Client{
		// turn off auto-redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	fmt.Println("callSealabsPayRefundApi 2")
	response, err := client.Do(request)
	fmt.Println(err)
	if err != nil {
		return err
	}

	fmt.Println("callSealabsPayRefundApi 3")
	if response.StatusCode != http.StatusOK {
		var errResponse *externalDtos.SealabsPayPaymentResponseError
		responseBodyBytes, err := ioutil.ReadAll(response.Body)
		fmt.Println("callSealabsPayRefundApi 4")
		fmt.Println(responseBodyBytes)
		fmt.Println(err)
		if err != nil {
			return err
		}
		json.Unmarshal(responseBodyBytes, &errResponse)
		return fmt.Errorf("%s", errResponse.Message)
	}
	fmt.Println("callSealabsPayRefundApi 5")

	return nil
}
