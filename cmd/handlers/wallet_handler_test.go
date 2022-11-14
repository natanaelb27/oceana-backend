package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/internal/mocks"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
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
		UserId:   uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		IsActive: true,
	}

	type fields struct {
		walletService *mocks.WalletService
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
					c.Set("id", model.UserId)
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
					c.Set("id", model.UserId)
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Bad Request","message":"Pin is required. PinConfirmation is required.","metadata":null}`,
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
					c.Set("id", model.UserId)
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
				mocks.NewWalletService(t),
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
