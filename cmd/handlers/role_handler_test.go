package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestHandler_RoleCreateFailingWithTransaction(t *testing.T) {
	var body io.Reader

	// userID := uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851")

	type fields struct {
		mockService *services.MockRoleService
	}

	model := &models.Role{
		Name: "TEST",
	}

	tests := []struct {
		name            string
		mock            func(f *fields)
		configureRouter func(r *gin.Engine, h *Handler, w *httptest.ResponseRecorder)
		fields          fields
		want            string
	}{
		{
			mock: func(f *fields) {
				f.mockService.On("ActivateTx", mock.Anything)
				f.mockService.On("DeactivateTx", mock.Anything)
				// f.mockService.On("DeactivateTx")
				f.mockService.On("MultipleOperationFailingWithTransaction", mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *Handler, w *httptest.ResponseRecorder) {
				r.Use(func(c *gin.Context) {
					c.Set("tx", &gorm.DB{})
					c.Next()
				})

				endpoint := "/roles/failing"
				handler := h.RoleCreateFailingWithTransaction
				r.POST(endpoint, handler)

				body = utils.NewBody(dtos.RoleCreateRequest{
					Name: "Azmi",
				})
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)

				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"00000000-0000-0000-0000-000000000000","name":"TEST","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockRoleService(t),
			}
			h := New(&HandlerConfig{
				RoleService: tt.fields.mockService,
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
