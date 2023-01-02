package handlers_test

import (
	"bytes"
	// "io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_MerchantGetInformationByID(t *testing.T) {
	response := &dtos.MerchantGetInformationByIDResponse{
		ID:              1,
		Name:            "Toko Adi",
		UserID:          uuid.MustParse("76747c91-6b8b-42b2-bbac-d76a8e2090db"),
		Photo:           nil,
		TotalRating:     7,
		RatingStar:      4.428571428571429,
		ProductCount:    1,
		ShopAddress:     "Jalan Merpati",
		ShopSubDistrict: "ABAH LUENG",
		ShopDistrict:    "BANDAR BARU",
		ShopCity:        "KABUPATEN PIDIE JAYA",
		ShopProvince:    "NANGGROE ACEH DARUSSALAM (NAD)",
	}

	type fields struct {
		merchantService *services.MockMerchantService
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
				f.merchantService.On("GetInformationByID", mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/shop/information"
				handler := h.MerchantGetInformationByID
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "merchant-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":1,"name":"Toko Adi","photo":null,"total_rating":7,"rating_star":4.428571428571429,"product_count":1,"user_id":"76747c91-6b8b-42b2-bbac-d76a8e2090db","shop_address":"Jalan Merpati","shop_sub_district":"ABAH LUENG","shop_district":"BANDAR BARU","shop_city":"KABUPATEN PIDIE JAYA","shop_province":"NANGGROE ACEH DARUSSALAM (NAD)","joined_date":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/shop/information"
				handler := h.MerchantGetInformationByID
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "merchant-id",
							Value: "",
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
			mock: func(f *fields) {
				f.merchantService.On("GetInformationByID", mock.Anything).Return(nil, errors.ErrMerchantNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/shop/information"
				handler := h.MerchantGetInformationByID
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Params = []gin.Param{
						{
							Key:   "merchant-id",
							Value: "1",
						},
					}
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":404,"data":null,"error":"Not Found","message":"merchant not found","metadata":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockMerchantService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				MerchantService: tt.fields.merchantService,
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

func TestHandler_MerchantGetProfile(t *testing.T) {
	response := &dtos.MerchantGetProfileResponse{
		ID:        2,
		UserID:    uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
		Name:      "Toko Budi",
		AddressID: uuid.MustParse("97716f0e-8e31-4e10-bc50-f28682876820"),
		Photo:     nil,
	}

	type fields struct {
		merchantService *services.MockMerchantService
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
				f.merchantService.On("MerchantGetProfile", mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/my-shop/profile"
				handler := h.MerchantGetProfile
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":2,"user_id":"f5677e10-81a1-4b23-9190-a72388070c15","name":"Toko Budi","address_id":"97716f0e-8e31-4e10-bc50-f28682876820","photo":null},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/my-shop/profile"
				handler := h.MerchantGetProfile
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
				f.merchantService.On("MerchantGetProfile", mock.Anything).Return(nil, errors.ErrMerchantNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/my-shop/profile"
				handler := h.MerchantGetProfile
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":404,"data":null,"error":"Not Found","message":"merchant not found","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockMerchantService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				MerchantService: tt.fields.merchantService,
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

func TestHandler_MerchantUpdatePhoto(t *testing.T) {
	photo := "https://res.cloudinary.com/dnjy6s5li/image/upload/v1672649544/oceana/dwilc5qks6ehd9ohpy6a.jpg"
	response := &dtos.MerchantGetProfileResponse{
		ID:        2,
		UserID:    uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
		Name:      "Toko Budi",
		AddressID: uuid.MustParse("97716f0e-8e31-4e10-bc50-f28682876820"),
		Photo:     &photo,
	}

	type fields struct {
		merchantService *services.MockMerchantService
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
				f.merchantService.On("MerchantUpdatePhoto", mock.Anything).Return(response, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/my-shop/profile/photo"
				handler := h.MerchantUpdatePhoto
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				writer.WriteField("id", "f5677e10-81a1-4b23-9190-a72388070c15")
				part, _ := writer.CreateFormFile("photo", "file.jpg")
				part.Write([]byte(`sample`))
				writer.Close()
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":2,"user_id":"f5677e10-81a1-4b23-9190-a72388070c15","name":"Toko Budi","address_id":"97716f0e-8e31-4e10-bc50-f28682876820","photo":"https://res.cloudinary.com/dnjy6s5li/image/upload/v1672649544/oceana/dwilc5qks6ehd9ohpy6a.jpg"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/my-shop/profile/photo"
				handler := h.MerchantUpdatePhoto
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				r.Use(func(ctx *gin.Context) {
					ctx.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/my-shop/profile/photo"
				handler := h.MerchantUpdatePhoto
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				r.Use(func(ctx *gin.Context) {
					ctx.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.merchantService.On("MerchantUpdatePhoto", mock.Anything).Return(nil, errors.ErrMerchantNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/my-shop/profile/photo"
				handler := h.MerchantUpdatePhoto
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				writer.WriteField("id", "f5677e10-81a1-4b23-9190-a72388070c15")
				part, _ := writer.CreateFormFile("photo", "file.jpg")
				part.Write([]byte(`sample`))
				writer.Close()
				r.Use(func(ctx *gin.Context) {
					ctx.Set("id", response.UserID)
					ctx.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				r.ServeHTTP(w, req)
			},
			want: `{"code":404,"data":null,"error":"Not Found","message":"merchant not found","metadata":null}`,
		},
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockMerchantService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				MerchantService: tt.fields.merchantService,
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
