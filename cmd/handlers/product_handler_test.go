package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ProductDetailsGetByID(t *testing.T) {
	response := &dtos.ProductDetailsGetByIDResponse{
		ID:             4,
		Name:           "Shirt D",
		MerchantID:     2,
		IsHazardous:    false,
		Weight:         100,
		PackageSize:    "1x2x3",
		Condition:      "new",
		Status:         "active",
		CategoryID:     3,
		Category:       "Kitchen",
		Images:         pq.StringArray{},
		Videos:         pq.StringArray{},
		Description:    "Baju D Gak Ada Variasi",
		TotalStock:     91,
		LikeCount:      1,
		TotalPurchased: 11,
		Price:          25000,
		PriceMax:       25000,
		PriceMin:       25000,
	}

	type fields struct {
		productService *services.MockProductService
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
				f.productService.On("ProductDetailsGetByID", mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products"
				handler := h.ProductDetailsGetByID
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", "76747c91-6b8b-42b2-bbac-d76a8e2090db")
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":4,"name":"Shirt D","merchant_id":2,"is_hazardous":false,"weight":100,"package_size":"1x2x3","condition":"new","status":"active","category_id":3,"category":"Kitchen","images":[],"videos":[],"description":"Baju D Gak Ada Variasi","total_stock":91,"like_count":1,"total_purchased":11,"product_review_details":null,"discount_type":null,"discount":null,"price":25000,"price_after_discount":null,"price_max":25000,"price_max_after_discount":null,"price_min":25000,"price_min_after_discount":null,"product_model":null,"product_variation":null},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products"
				handler := h.ProductDetailsGetByID
				r.Use(func(ctx *gin.Context) {
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.productService.On("ProductDetailsGetByID", mock.Anything).Return(nil, errors.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products"
				handler := h.ProductDetailsGetByID
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", "76747c91-6b8b-42b2-bbac-d76a8e2090db")
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
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
				services.NewMockProductService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				ProductService: tt.fields.productService,
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

func TestHandler_ViewProduct(t *testing.T) {

	response := &dtos.GetProductViewResponse{
		UserID:    uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
		ProductID: 1,
		Count:     1,
	}

	type fields struct {
		productService *services.MockProductService
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
				f.productService.On("ViewProduct", mock.Anything, mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/view"
				handler := h.ViewProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"user_id":"f5677e10-81a1-4b23-9190-a72388070c15","product_id":1,"last_viewed_at":"0001-01-01T00:00:00Z","count":1},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.productService.On("ViewProduct", mock.Anything, mock.Anything).Return(nil, errors.ErrDataNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/view"
				handler := h.ViewProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":404,"data":null,"error":"Not Found","message":"data not found","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/view"
				handler := h.ViewProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/view"
				handler := h.ViewProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "a",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockProductService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				ProductService: tt.fields.productService,
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

func TestHandler_GetFavoriteProduct(t *testing.T) {
	response := &dtos.GetProductFavoriteResponse{
		UserID:     uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
		ProductID:  1,
		IsFavorite: true,
	}

	type fields struct {
		productService *services.MockProductService
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
				f.productService.On("GetFavoriteProduct", mock.Anything, mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.GetFavoriteProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"user_id":"f5677e10-81a1-4b23-9190-a72388070c15","product_id":1,"is_favorite":true},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.productService.On("GetFavoriteProduct", mock.Anything, mock.Anything).Return(nil, errors.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.GetFavoriteProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.GetFavoriteProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.GetFavoriteProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Next()
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
				services.NewMockProductService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				ProductService: tt.fields.productService,
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

func TestHandler_FavoriteProduct(t *testing.T) {
	response := &dtos.GetProductFavoriteResponse{
		UserID:     uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
		ProductID:  1,
		IsFavorite: true,
	}
	type fields struct {
		productService *services.MockProductService
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
				f.productService.On("FavoriteProduct", mock.Anything, mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.FavoriteProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"user_id":"f5677e10-81a1-4b23-9190-a72388070c15","product_id":1,"is_favorite":true},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.productService.On("FavoriteProduct", mock.Anything, mock.Anything).Return(nil, errors.ErrDataNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.FavoriteProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":404,"data":null,"error":"Not Found","message":"data not found","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.FavoriteProduct
				r.Use(func(ctx *gin.Context) {
					// ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/favorite"
				handler := h.FavoriteProduct
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockProductService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				ProductService: tt.fields.productService,
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

func TestHandler_ProductCheckCourierService(t *testing.T) {
	var body io.Reader
	response := &product_dtos.ProductCheckCourierServiceResponse{
		AvailableProducts: pq.Int64Array{},
	}
	type fields struct {
		productService *services.MockProductService
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
				f.productService.On("CheckCourierService", mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/check-courier-service"
				handler := h.ProductCheckCourierService
				body = NewBody(&product_dtos.ProductCheckCourierService{
					CourierCode:        "JNE",
					CourierServiceCode: "EXPRESS",
					ProductIDs:         []uint{1},
				})
				r.Use(func(ctx *gin.Context) {
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"available_products":[]},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.productService.On("CheckCourierService", mock.Anything).Return(nil, errors.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/check-courier-service"
				handler := h.ProductCheckCourierService
				body = NewBody(&product_dtos.ProductCheckCourierService{
					CourierCode:        "JNE",
					CourierServiceCode: "EXPRESS",
					ProductIDs:         []uint{1},
				})
				r.Use(func(ctx *gin.Context) {
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/check-courier-service"
				handler := h.ProductCheckCourierService
				body = NewBody(&product_dtos.ProductCheckCourierService{})
				r.Use(func(ctx *gin.Context) {
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"CourierCode is required. CourierServiceCode is required. ProductIDs is required.","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockProductService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				ProductService: tt.fields.productService,
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

func TestHandler_ProductModelById(t *testing.T) {
	response := &product_dtos.ProductModel{}
	type fields struct {
		productService *services.MockProductService
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
				f.productService.On("GetUpdateModelByProductID", mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/model"
				handler := h.ProductModelById
				r.Use(func(ctx *gin.Context) {
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":0,"name":"","description":"","category_id":0,"merchant_id":0,"weight":0,"is_hazardous":null,"package_size":"","condition":"","status":"","process_day":0,"images":null,"videos":null,"product_items":null,"variants":null,"couriers":null},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.productService.On("GetUpdateModelByProductID", mock.Anything).Return(nil, errors.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/model"
				handler := h.ProductModelById
				r.Use(func(ctx *gin.Context) {
					ctx.Params = []gin.Param{
						{
							Key:   "product-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/products/model"
				handler := h.ProductModelById
				r.Use(func(ctx *gin.Context) {
					ctx.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockProductService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				ProductService: tt.fields.productService,
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
