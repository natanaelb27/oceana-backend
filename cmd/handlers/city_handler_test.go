package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/internal/services"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CitiesGetAll(t *testing.T) {
	cities := []*dtos.CityResponse{
		{
			ID:   1,
			Name: "KABUPATEN ACEH BARAT",
		},
	}

	type fields struct {
		cityService *services.MockCityService
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
				f.cityService.On("GetAll", mock.Anything).Return(cities, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/cities"
				handler := h.CitiesGetAll

				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"id":1,"name":"KABUPATEN ACEH BARAT"}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.cityService.On("GetAll", mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/cities"
				handler := h.CitiesGetAll

				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockCityService(t),
				}

				h := handlers.New(&handlers.HandlerConfig{
					CityService: tt.fields.cityService,
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
