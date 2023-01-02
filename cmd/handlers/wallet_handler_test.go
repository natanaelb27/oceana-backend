package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_WalletActivate(t *testing.T) {
	var body io.Reader

	model := &models.Wallet{
		ID:       uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:   uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		IsActive: true,
	}

	type fields struct {
		walletService *services.MockWalletService
	}

	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		{
			mock: func(f *fields) {
				f.walletService.On("Activate", mock.Anything, mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/wallet/activate"
				handler := h.WalletActivate
				body = NewBody(dtos.WalletActivateRequest{
					Pin:             "123456",
					PinConfirmation: "123456",
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","balance":0,"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/wallet/activate"
				handler := h.WalletActivate
				body = NewBody(dtos.WalletActivateRequest{})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"Pin is required. PinConfirmation is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.walletService.On("Activate", mock.Anything, mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/wallet/activate"
				handler := h.WalletActivate
				body = NewBody(dtos.WalletActivateRequest{
					Pin:             "123456",
					PinConfirmation: "123456",
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/wallet/activate"
				handler := h.WalletActivate
				body = NewBody(dtos.WalletActivateRequest{
					Pin:             "123456",
					PinConfirmation: "123456",
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockWalletService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				WalletService: tt.fields.walletService,
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

func TestHandler_WalletChangePin(t *testing.T) {
	var body io.Reader

	model := &models.Wallet{
		ID:       uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:   uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		IsActive: true,
	}

	type fields struct {
		walletService *services.MockWalletService
	}

	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		{
			mock: func(f *fields) {
				f.walletService.On("ChangePin", mock.Anything, mock.Anything, mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})

				endpoint := "/wallet/activate"
				handler := h.WalletChangePin
				r.PATCH(endpoint, handler)

				body = NewBody(dtos.WalletChangePinRequest{
					Pin:             "123456",
					PinConfirmation: "123456",
				})
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.walletService.On("ChangePin", mock.Anything, mock.Anything, mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})

				endpoint := "/wallet/activate"
				handler := h.WalletChangePin
				r.PATCH(endpoint, handler)

				body = NewBody(dtos.WalletChangePinRequest{
					Pin:             "123456",
					PinConfirmation: "123456",
				})
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})

				endpoint := "/wallet/activate"
				handler := h.WalletChangePin
				r.PATCH(endpoint, handler)

				body = NewBody(dtos.WalletChangePinRequest{
					Pin: "123456",
				})
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"PinConfirmation is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Set("id", nil)
					c.Next()
				})

				endpoint := "/wallet/activate"
				handler := h.WalletChangePin
				r.PATCH(endpoint, handler)

				body = NewBody(dtos.WalletChangePinRequest{
					Pin: "123456",
				})
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockWalletService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				WalletService: tt.fields.walletService,
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

func TestHandler_WalletGetDetails(t *testing.T) {
	var body io.Reader
	datetime, _ := time.Parse("2006-Jan-02", "2023-Jan-02")
	model := &dtos.WalletDetailsResponse{
		ID:      uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:  uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		Balance: 50000,
		WalletTransactionHistory: []*internalDtos.WalletTransactionHistory{
			{
				ID:                  uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
				WalletID:            uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b853"),
				Amount:              22,
				TransactionDatetime: datetime,
				TransactionType:     "TOPUP",
				Type:                "TEST",
			},
		},
		IsActive:              true,
		InvalidPinCount:       1,
		BlockedExpirationTime: nil,
		CreatedAt:             datetime,
		UpdatedAt:             datetime,
	}

	pagination := &utils.Pagination{
		Limit:      1,
		Page:       1,
		TotalRows:  1,
		TotalPages: 1,
	}

	type fields struct {
		walletService *services.MockWalletService
	}

	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		{
			mock: func(f *fields) {
				f.walletService.On("GetWalletDetails", mock.Anything, mock.Anything).Return(model, pagination, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})

				endpoint := "/wallet"
				handler := h.WalletGetDetails

				body = NewBody(dtos.WalletDetailsQuery{})

				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","user_id":"e56a013f-58c3-4dd3-93d1-190d6a46b852","balance":50000,"wallet_transaction_history":[{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","wallet_id":"e56a013f-58c3-4dd3-93d1-190d6a46b853","amount":22,"transaction_datetime":"2023-01-02T00:00:00Z","transaction_type":"TOPUP","type":"TEST"}],"is_active":true,"invalid_pin_count":1,"created_at":"2023-01-02T00:00:00Z","updated_at":"2023-01-02T00:00:00Z"},"error":null,"message":"success","metadata":{"limit":1,"page":1,"total_rows":1,"total_pages":1}}`,
		},
		{
			mock: func(f *fields) {
				f.walletService.On("GetWalletDetails", mock.Anything, mock.Anything).Return(nil, nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})

				endpoint := "/wallet"
				handler := h.WalletGetDetails

				body = NewBody(dtos.WalletDetailsQuery{})

				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Next()
				})

				endpoint := "/wallet"
				handler := h.WalletGetDetails

				body = NewBody(dtos.WalletDetailsQuery{})

				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockWalletService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				WalletService: tt.fields.walletService,
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

// func TestHandler_WalletTopup(t *testing.T) {
// 	var body io.Reader

// 	userID := uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851")

// 	type fields struct {
// 		walletService *services.MockWalletService
// 	}

// 	tests := []struct {
// 		name            string
// 		mock            func(f *fields)
// 		configureRouter func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder)
// 		fields          fields
// 		want            string
// 	}{
// 		{
// 			mock: func(f *fields) {
// 				// f.walletService.On("WithTx", mock.Anything).Return().
// 				// f.walletService.On("CreateWalletTopup", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			},
// 			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
// 				r.Use(func(c *gin.Context) {
// 					c.Set("id", userID)
// 					c.Set("tx", &gorm.DB{})
// 					c.Next()
// 				})

// 				endpoint := "/wallet/activate"
// 				handler := h.WalletTopup
// 				r.PATCH(endpoint, handler)

// 				body = NewBody(dtos.TopupWalletRequest{
// 					CardNumber:   "weewwe",
// 					TxID:         "werwerewr",
// 					Amount:       "3000",
// 					Merchantcode: "jjoewr",
// 					Status:       "OK",
// 					Message:      "Success",
// 					Signature:    "JU9R82U932U9",
// 				})
// 				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)

// 				r.ServeHTTP(w, req)
// 			},
// 			want: `{"code":200,"data":{"id":"e56a013f-58c3-4dd3-93d1-190d6a46b851","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.fields = fields{
// 				services.NewMockWalletService(t),
// 			}
// 			h := handlers.New(&handlers.HandlerConfig{
// 				WalletService: tt.fields.walletService,
// 			})

// 			tt.mock(&tt.fields)

// 			w := httptest.NewRecorder()
// 			r := gin.Default()
// 			tt.configureRouter(r, h, w)

// 			response := w.Body.String()

// 			assert.Equal(t, tt.want, response)
// 		})
// 	}
// }
