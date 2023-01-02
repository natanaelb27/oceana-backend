package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_SealabsPayCreate(t *testing.T) {
	var body io.Reader

	result := &dtos.SealabsPayCreateResponse{
		ID:         uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		CardNumber: "2783993778282621",
		NameOnCard: "Timothy",
		IsDefault:  true,
	}

	type fields struct {
		sealabsPayService *services.MockSealabsPayService
	}
	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		// TODO: Add test cases.
		// {
		// 	mock: func(f *fields) {
		// 		f.sealabsPayService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(result, nil)
		// 	},
		// 	configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
		// 		endpoint := "/sealabs-pay"
		// 		handler := h.SealabsPayCreate
		// 		body = NewBody(&dtos.SealabsPayCreateRequest{
		// 			NameOnCard: "Timothy",
		// 			CardNumber: "2783993778282621",
		// 			TxnID:      "1",
		// 			Signature:  "aa0d71ab8ef998e7a40854646635fe2b5360fc655ae415476e17afbc066f127e",
		// 		})
		// 		r.Use(func(c *gin.Context) {
		// 			c.Set("id", result.ID)
		// 			c.Next()
		// 		})
		// 		r.POST(endpoint, handler)
		// 		req, _ := http.NewRequest(http.MethodPost, endpoint, body)
		// 		r.ServeHTTP(w, req)
		// 	},
		// 	want: `{"code":200,"data":{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","card_number":2783993778282621,"active_date":"0001-01-01T00:00:00Z","name_on_card":"Timothy","is_default":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		// },
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayCreate
				body = NewBody(&dtos.SealabsPayCreateRequest{
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", result.ID)
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"NameOnCard is required. CardNumber is required. TxnID is required. Signature is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayCreate
				body = NewBody(&dtos.SealabsPayCreateRequest{
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockSealabsPayService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				SealabsPayService: tt.fields.sealabsPayService,
			})

			tt.mock(&tt.fields)

			w := httptest.NewRecorder()
			r := gin.Default()
			tt.configureRouter(r, h, w)

			response := w.Body.String()

			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_SealabsPayGetAll(t *testing.T) {
	results := make([]*models.SealabsPay, 0)
	model := &models.SealabsPay{
		ID:            uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		AccountNumber: "2783993778282621",
		NameOnCard:    "Timothy",
		IsDefault:     true,
	}
	results = append(results, model)

	type fields struct {
		sealabsPayService *services.MockSealabsPayService
	}
	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		// TODO: Add test cases.
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("GetAll", mock.Anything).Return(results, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","card_number":"2783993778282621","active_date":"0001-01-01T00:00:00Z","name_on_card":"Timothy","is_default":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("GetAll", mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", "")
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockSealabsPayService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				SealabsPayService: tt.fields.sealabsPayService,
			})

			tt.mock(&tt.fields)

			w := httptest.NewRecorder()
			r := gin.Default()
			tt.configureRouter(r, h, w)

			response := w.Body.String()

			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_SealabsPayGetById(t *testing.T) {
	model := &models.SealabsPay{
		ID:            uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		AccountNumber: "2783993778282621",
		NameOnCard:    "Timothy",
		IsDefault:     true,
	}
	type fields struct {
		sealabsPayService *services.MockSealabsPayService
	}

	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		// TODO: Add test cases.
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("GetById", mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayGetById
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.AddParam("sealabs-pay-id", model.ID.String())
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","card_number":"2783993778282621","active_date":"0001-01-01T00:00:00Z","name_on_card":"Timothy","is_default":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("GetById", mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayGetById
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.AddParam("sealabs-pay-id", model.ID.String())
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":404,"data":null,"error":"Not Found","message":"data not found","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay"
				handler := h.SealabsPayGetById
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.AddParam("sealabs-pay-id", "")
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockSealabsPayService(t),
				}
				h := handlers.New(&handlers.HandlerConfig{
					SealabsPayService: tt.fields.sealabsPayService,
				})

				tt.mock(&tt.fields)

				w := httptest.NewRecorder()
				r := gin.Default()
				tt.configureRouter(r, h, w)

				response := w.Body.String()

				assert.Equal(t, tt.want, response)
			})
		})
	}
}

func TestHandler_SealabsPayRegisterOTP(t *testing.T) {
	var body io.Reader

	result := "test.com"

	type fields struct {
		sealabsPayService *services.MockSealabsPayService
	}
	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		// TODO: Add test cases.
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("RegisterOTP", mock.Anything).Return(&result, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/register-otp"
				handler := h.SealabsPayRegisterOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{
					NameOnCard: "Timothy",
					CardNumber: "2783993778282621",
					Amount:     "1",
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":"test.com","error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/register-otp"
				handler := h.SealabsPayRegisterOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"CardNumber is required. Amount is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("RegisterOTP", mock.Anything).Return(nil, e.ErrSealabsPayCardNumberAlreadyExist)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/register-otp"
				handler := h.SealabsPayRegisterOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{
					CardNumber: "2783993778282621",
					Amount:     "1",
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":409,"data":null,"error":"Conflict","message":"account number is already registered","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockSealabsPayService(t),
				}
				h := handlers.New(&handlers.HandlerConfig{
					SealabsPayService: tt.fields.sealabsPayService,
				})

				tt.mock(&tt.fields)

				w := httptest.NewRecorder()
				r := gin.Default()
				tt.configureRouter(r, h, w)

				response := w.Body.String()

				assert.Equal(t, tt.want, response)
			})
		})
	}
}

func TestHandler_SealabsPayTopupOTP(t *testing.T) {
	var body io.Reader

	result := "test.com"

	type fields struct {
		sealabsPayService *services.MockSealabsPayService
	}
	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		// TODO: Add test cases.
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("TopupWalletOTP", mock.Anything).Return(&result, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/topup-otp"
				handler := h.SealabsPayTopupOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{
					NameOnCard: "Timothy",
					CardNumber: "2783993778282621",
					Amount:     "1",
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":"test.com","error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/topup-otp"
				handler := h.SealabsPayTopupOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"CardNumber is required. Amount is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("TopupWalletOTP", mock.Anything).Return(nil, e.ErrInvalidSealabsPayAction)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/topup-otp"
				handler := h.SealabsPayTopupOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{
					CardNumber: "2783993778282621",
					Amount:     "1",
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"invalid sealabs pay action","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockSealabsPayService(t),
				}
				h := handlers.New(&handlers.HandlerConfig{
					SealabsPayService: tt.fields.sealabsPayService,
				})

				tt.mock(&tt.fields)

				w := httptest.NewRecorder()
				r := gin.Default()
				tt.configureRouter(r, h, w)

				response := w.Body.String()

				assert.Equal(t, tt.want, response)
			})
		})
	}
}

func TestHandler_SealabsPayPaymentOTP(t *testing.T) {
	var body io.Reader

	result := "test.com"

	type fields struct {
		sealabsPayService *services.MockSealabsPayService
	}
	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		// TODO: Add test cases.
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("PaymentOTP", mock.Anything).Return(&result, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/payment-otp"
				handler := h.SealabsPayPaymentOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{
					NameOnCard: "Timothy",
					CardNumber: "2783993778282621",
					Amount:     "1",
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":"test.com","error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/payment-otp"
				handler := h.SealabsPayPaymentOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"CardNumber is required. Amount is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("PaymentOTP", mock.Anything).Return(nil, e.ErrInvalidSealabsPayAction)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/sealabs-pay/payment-otp"
				handler := h.SealabsPayPaymentOTP
				body = NewBody(dtos.SealabsPayPaymentInboundRequest{
					CardNumber: "2783993778282621",
					Amount:     "1",
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"invalid sealabs pay action","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockSealabsPayService(t),
				}
				h := handlers.New(&handlers.HandlerConfig{
					SealabsPayService: tt.fields.sealabsPayService,
				})

				tt.mock(&tt.fields)

				w := httptest.NewRecorder()
				r := gin.Default()
				tt.configureRouter(r, h, w)

				response := w.Body.String()

				assert.Equal(t, tt.want, response)
			})
		})
	}
}

func TestHandler_SealabsPayCallback(t *testing.T) {
	var body io.Reader

	result := &dtos.SealabsPayCallbackRequest{
		TxnID:        "1",
		Amount:       "1",
		MerchantCode: "1",
		Status:       "1",
		Message:      "1",
		Signature:    "1",
	}

	type fields struct {
		sealabsPayService *services.MockSealabsPayService
	}
	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		// TODO: Add test cases.
		{
			mock: func(f *fields) {
				f.sealabsPayService.On("StoreSignature", mock.Anything, mock.Anything).Return(nil, e.ErrCacheOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/callback/sealabs-pay"
				handler := h.SealabsPayCallback
				body = NewBody(result)
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"txn_id":"1","amount":"1","merchant_code":"1","status":"1","message":"1","signature":"1"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/callback/sealabs-pay"
				handler := h.SealabsPayCallback
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"invalid request","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockSealabsPayService(t),
				}
				h := handlers.New(&handlers.HandlerConfig{
					SealabsPayService: tt.fields.sealabsPayService,
				})

				tt.mock(&tt.fields)

				w := httptest.NewRecorder()
				r := gin.Default()
				tt.configureRouter(r, h, w)

				response := w.Body.String()

				assert.Equal(t, tt.want, response)
			})
		})
	}
}
