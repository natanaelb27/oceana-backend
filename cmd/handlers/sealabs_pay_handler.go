package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

func (h *Handler) SealabsPayCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.SealabsPayCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	sealabsPayCreateResponse, err := h.sealabsPayService.Create(id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, sealabsPayCreateResponse, nil)
}

func (h *Handler) SealabsPayGetAll(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	results, err := h.sealabsPayService.GetAll(id)
	if err != nil {
		utils.ResponseError(c, err)
	}

	var responseData []*dtos.SealabsPayCreateResponse = make([]*dtos.SealabsPayCreateResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewSealabsPayCreateResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)

}

func (h *Handler) SealabsPayGetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("sealabs_pay_id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.sealabsPayService.GetById(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewSealabsPayCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) SealabsPayDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("sealabs_pay_id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.sealabsPayService.Delete(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewSealabsPayDeleteResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) SealabsPayTopup(c *gin.Context) {
	var inboundRequest *dtos.SealabsPayPaymentInboundRequest
	err := c.ShouldBindJSON(&inboundRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	client := &http.Client{
		// turn off auto-redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	apiURL := os.Getenv("SEALABSPAY_API_URL")
	apiKey := os.Getenv("SEALABSPAY_API_KEY")
	merchantCode := os.Getenv("SEALABSPAY_MERCHANT_CODE")
	redirectBaseURL := os.Getenv("FRONTEND_URL")
	redirectURL := fmt.Sprintf("%s/sealabs-pay/redirect?card_number=%s", redirectBaseURL, inboundRequest.CardNumber)
	callbackBaseURL := os.Getenv("SEALABSPAY_CALLBACK_BASE_URL")
	callbackURL := fmt.Sprintf("%s/top-up", callbackBaseURL)
	// hashing signature
	hasher := hmac.New(sha256.New, []byte(apiKey))
	signatureStringToEncrypt := fmt.Sprintf("%s:%s:%s", inboundRequest.CardNumber, inboundRequest.Amount, merchantCode)
	hasher.Write([]byte(signatureStringToEncrypt))
	signatureBytes := hasher.Sum(nil)
	signatureString := hex.EncodeToString(signatureBytes)

	outboundRequest := url.Values{}
	outboundRequest.Set("card_number", inboundRequest.CardNumber)
	outboundRequest.Set("amount", inboundRequest.Amount)
	outboundRequest.Set("merchant_code", merchantCode)
	outboundRequest.Set("redirect_url", redirectURL)
	outboundRequest.Set("callback_url", callbackURL)
	outboundRequest.Set("signature", signatureString)

	encodedOutboundRequest := outboundRequest.Encode()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/pay", apiURL), strings.NewReader(encodedOutboundRequest))
	if err != nil {
		utils.ResponseError(c, e.ErrUnknown)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedOutboundRequest)))

	res, err := client.Do(req)
	if err != nil {
		utils.ResponseError(c, e.ErrUnknown)
		return
	}

	var errResponse *dtos.SealabsPayRedirectResponseError
	responseBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.ResponseError(c, e.ErrInternalOperationFailed)
		return
	}

	json.Unmarshal(responseBodyBytes, &errResponse)
	if res.StatusCode != http.StatusSeeOther {
		utils.Response(c, res.StatusCode, errResponse, nil, nil, res.Status)
		return
	}

	otpURL, ok := res.Header["Location"]
	if !ok {
		utils.ResponseError(c, e.ErrUnknown)
		return
	}

	utils.ResponseSuccess(c, otpURL[0], nil)
}

func (h *Handler) SealabsPayTopupCallback(c *gin.Context) {
	var request *dtos.SealabsPayCallbackRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	err = h.sealabsPayService.StoreSignature(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, request, nil)
}
