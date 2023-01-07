package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/services"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CourierGetAll(t *testing.T) {
	type fields struct {
		courierService *services.MockCourierService
	}

	couriers := []*dtos.Courier{
		{
			ID:   1,
			Code: "jne",
			Name: "JNE",
			CourierService: []*dtos.CourierService{
				{
					ID:          2,
					Code:        "CTCYES",
					Description: "JNE City Courier",
				},
				{
					ID:          3,
					Code:        "OKE",
					Description: "Ongkos Kirim Ekonomis",
				},
				{
					ID:          1,
					Code:        "CTC",
					Description: "JNE City Courier",
				},
				{
					ID:          4,
					Code:        "REG",
					Description: "Layanan Reguler",
				},
			},
		},
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
				f.courierService.On("GetAll").Return(couriers, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/couriers"
				handler := h.CourierGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"id":1,"code":"jne","name":"JNE","courier_services":[{"id":2,"code":"CTCYES","description":"JNE City Courier"},{"id":3,"code":"OKE","description":"Ongkos Kirim Ekonomis"},{"id":1,"code":"CTC","description":"JNE City Courier"},{"id":4,"code":"REG","description":"Layanan Reguler"}]}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.courierService.On("GetAll").Return(nil, e.ErrUnauthorized)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/couriers"
				handler := h.CourierGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":401,"data":null,"error":"Unauthorized","message":"access is denied due to invalid credentials","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockCourierService(t),
				}

				h := handlers.New(&handlers.HandlerConfig{
					CourierService: tt.fields.courierService,
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
