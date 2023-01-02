package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CategoryGetAll(t *testing.T) {
	type fields struct {
		categoryService *services.MockCategoryService
	}

	models := []*internalDtos.Category{
		{
			ID:    1,
			Image: "https://ecs7.tokopedia.net/img/WgKiGm/2021/6/23/a4e1aa64-556c-4838-808e-8c008da1225d.png",
			Name:  "Animal Care",
			Level: 1,
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
				f.categoryService.On("GetAll", mock.Anything).Return(models, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories"
				handler := h.CategoryGetAll
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"id":1,"image":"https://ecs7.tokopedia.net/img/WgKiGm/2021/6/23/a4e1aa64-556c-4838-808e-8c008da1225d.png","name":"Animal Care","level":1,"category":null}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.categoryService.On("GetAll", mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories"
				handler := h.CategoryGetAll
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockCategoryService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				CategoryService: tt.fields.categoryService,
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

func TestHandler_CategoryAncestorsGetByID(t *testing.T) {
	type fields struct {
		categoryService *services.MockCategoryService
	}

	models := []*dtos.GetCategoriesResponse{
		{
			ID:       23,
			MediaURL: "https://ecs7.tokopedia.net/img/WgKiGm/2021/6/23/a4e1aa64-556c-4838-808e-8c008da1225d.png",
			Name:     "Animal Care",
			ParentID: 0,
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
				f.categoryService.On("GetAncestorsByID", mock.Anything).Return(models, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories/ancestors/337"
				handler := h.CategoryAncestorsGetByID
				r.Use(func(c *gin.Context) {
					c.Params = []gin.Param{
						{
							Key:   "category-id",
							Value: "337",
						},
					}
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"id":23,"media_url":"https://ecs7.tokopedia.net/img/WgKiGm/2021/6/23/a4e1aa64-556c-4838-808e-8c008da1225d.png","name":"Animal Care","parent_id":0}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories/ancestors/337"
				handler := h.CategoryAncestorsGetByID
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.categoryService.On("GetAncestorsByID", mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories/ancestors/337"
				handler := h.CategoryAncestorsGetByID
				r.Use(func(c *gin.Context) {
					c.Params = []gin.Param{
						{
							Key:   "category-id",
							Value: "337",
						},
					}
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockCategoryService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				CategoryService: tt.fields.categoryService,
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

func TestHandler_GetMerchantCategories(t *testing.T) {
	type fields struct {
		categoryService *services.MockCategoryService
	}

	models := []*internalDtos.Category{
		{
			ID:    1,
			Image: "https://ecs7.tokopedia.net/img/WgKiGm/2021/6/23/a4e1aa64-556c-4838-808e-8c008da1225d.png",
			Name:  "Animal Care",
			Level: 1,
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
				f.categoryService.On("GetMerchantCategories", mock.Anything).Return(models, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories/merchants/1"
				handler := h.GetMerchantCategories
				r.Use(func(c *gin.Context) {
					c.Params = []gin.Param{
						{
							Key:   "merchant-id",
							Value: "1",
						},
					}
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"id":1,"image":"https://ecs7.tokopedia.net/img/WgKiGm/2021/6/23/a4e1aa64-556c-4838-808e-8c008da1225d.png","name":"Animal Care","level":1,"category":null}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories/merchants/1"
				handler := h.GetMerchantCategories
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.categoryService.On("GetMerchantCategories", mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/categories/merchants/1"
				handler := h.GetMerchantCategories
				r.Use(func(c *gin.Context) {
					c.Params = []gin.Param{
						{
							Key:   "merchant-id",
							Value: "1",
						},
					}
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockCategoryService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				CategoryService: tt.fields.categoryService,
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
