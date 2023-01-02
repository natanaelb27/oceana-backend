package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	pkgDTO "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_AddressCreate(t *testing.T) {
	var body io.Reader

	model := &models.Address{
		ID:            uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"),
		UserID:        uuid.MustParse("a110dddd-29a1-49a4-bd49-7594d024c354"),
		Name:          "Home",
		Detail:        "Jl. Jenderal Sudirman No. 1",
		SubDistrictID: 1,
		IsDefault:     true,
	}

	type fields struct {
		addressService *services.MockAddressService
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
				f.addressService.On("Create", mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressCreate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Name:          "Home",
					Detail:        "Jl. Jenderal Sudirman No. 1",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"a66de3fd-34ca-44c2-b283-3a25bd835685","user_id":"a110dddd-29a1-49a4-bd49-7594d024c354","name":"Home","detail":"Jl. Jenderal Sudirman No. 1","sub_district_id":1,"is_default":true},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressCreate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Detail:        "Jl. Jenderal Sudirman No. 1",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"Name is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressCreate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Detail:        "Jl. Jenderal Sudirman No. 1",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("Create", mock.Anything).Return(nil, e.ErrUnauthorized)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressCreate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Name:          "Home",
					Detail:        "Jl. Jenderal Sudirman No. 1",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":401,"data":null,"error":"Unauthorized","message":"access is denied due to invalid credentials","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("Create", mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressCreate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Name:          "Home",
					Detail:        "Jl. Jenderal Sudirman No. 1",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockAddressService(t),
			}
			h := handlers.New(&handlers.HandlerConfig{
				AddressService: tt.fields.addressService,
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

func TestHandler_AddressGetAll(t *testing.T) {

	type fields struct {
		addressService *services.MockAddressService
	}

	addresses := []*dtos.AddressDetails{
		{
			ID:            uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"),
			UserID:        uuid.MustParse("a110dddd-29a1-49a4-bd49-7594d024c354"),
			PhoneNumber:   "081234567789",
			Name:          "Home",
			Detail:        "Jl. Jenderal Sudirman No. 1",
			SubDistrictID: 1,
			SubDistrict:   "Kalibata",
			ZipCode:       "12740",
			DistrictID:    1,
			District:      "Pancoran",
			CityID:        1,
			City:          "Jakarta Selatan",
			ProvinceID:    1,
			Province:      "DKI Jakarta",
			IsDefault:     true,
			IsShopAddress: false,
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
				f.addressService.On("GetProfileAddressess", mock.Anything).Return(addresses, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", addresses[0].UserID)
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"id":"a66de3fd-34ca-44c2-b283-3a25bd835685","user_id":"a110dddd-29a1-49a4-bd49-7594d024c354","phone_number":"081234567789","name":"Home","detail":"Jl. Jenderal Sudirman No. 1","sub_district_id":1,"sub_district":"Kalibata","zip_code":"12740","district_id":1,"district":"Pancoran","city_id":1,"city":"Jakarta Selatan","province_id":1,"province":"DKI Jakarta","is_default":true,"is_shop_address":false}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("GetProfileAddressess", mock.Anything).Return(nil, e.ErrUnauthorized)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses"
				handler := h.AddressGetAll
				r.Use(func(c *gin.Context) {
					c.Set("id", addresses[0].UserID)
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":401,"data":null,"error":"Unauthorized","message":"access is denied due to invalid credentials","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockAddressService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				AddressService: tt.fields.addressService,
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

func TestHandler_AddressGetById(t *testing.T) {
	type fields struct {
		addressService *services.MockAddressService
	}

	address := &dtos.AddressDetails{
		ID:            uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"),
		UserID:        uuid.MustParse("a110dddd-29a1-49a4-bd49-7594d024c354"),
		PhoneNumber:   "081234567789",
		Name:          "Home",
		Detail:        "Jl. Jenderal Sudirman No. 1",
		SubDistrictID: 1,
		SubDistrict:   "Kalibata",
		ZipCode:       "12740",
		DistrictID:    1,
		District:      "Pancoran",
		CityID:        1,
		City:          "Jakarta Selatan",
		ProvinceID:    1,
		Province:      "DKI Jakarta",
		IsDefault:     true,
		IsShopAddress: false,
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
				f.addressService.On("GetProfileAddressById", mock.Anything).Return(address, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + address.ID.String()
				handler := h.AddressGetById
				r.Use(func(c *gin.Context) {
					c.Set("id", address.UserID)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: address.ID.String(),
						},
					}
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"a66de3fd-34ca-44c2-b283-3a25bd835685","user_id":"a110dddd-29a1-49a4-bd49-7594d024c354","phone_number":"081234567789","name":"Home","detail":"Jl. Jenderal Sudirman No. 1","sub_district_id":1,"sub_district":"Kalibata","zip_code":"12740","district_id":1,"district":"Pancoran","city_id":1,"city":"Jakarta Selatan","province_id":1,"province":"DKI Jakarta","is_default":true,"is_shop_address":false},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/"
				handler := h.AddressGetById
				r.Use(func(c *gin.Context) {
					c.Set("id", address.UserID)
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("GetProfileAddressById", mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/c1fd4939-6b73-41a0-b847-fc3902332d72"
				handler := h.AddressGetById
				r.Use(func(c *gin.Context) {
					c.Set("id", address.UserID)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: "c1fd4939-6b73-41a0-b847-fc3902332d72",
						},
					}
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
				f.addressService.On("GetProfileAddressById", mock.Anything).Return(nil, e.ErrUnauthorized)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + address.ID.String()
				handler := h.AddressGetById
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: address.ID.String(),
						},
					}
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":401,"data":null,"error":"Unauthorized","message":"access is denied due to invalid credentials","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockAddressService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				AddressService: tt.fields.addressService,
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

func TestHandler_AddressUpdate(t *testing.T) {
	var body io.Reader

	model := &models.Address{
		ID:            uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"),
		UserID:        uuid.MustParse("a110dddd-29a1-49a4-bd49-7594d024c354"),
		Name:          "Home",
		Detail:        "Jl. Abubakar Ibrahim No. 29E",
		SubDistrictID: 1,
		IsDefault:     true,
	}

	type fields struct {
		addressService *services.MockAddressService
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
				f.addressService.On("UpdateProfileAddress", mock.Anything, mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String()
				handler := h.AddressUpdate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Name:          "Home",
					Detail:        "Jl. Abubakar Ibrahim No. 29E",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.ID.String(),
						},
					}
					c.Next()
				})
				r.PUT(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPut, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"a66de3fd-34ca-44c2-b283-3a25bd835685","user_id":"a110dddd-29a1-49a4-bd49-7594d024c354","name":"Home","detail":"Jl. Abubakar Ibrahim No. 29E","sub_district_id":1,"is_default":true},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/"
				handler := h.AddressUpdate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Name:          "Home",
					Detail:        "Jl. Abubakar Ibrahim No. 29E",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.PUT(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPut, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("UpdateProfileAddress", mock.Anything, mock.Anything).Return(nil, e.ErrInvalidIdentifier)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String()
				handler := h.AddressUpdate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Name:          "Home",
					Detail:        "Jl. Abubakar Ibrahim No. 29E",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.ID.String(),
						},
					}
					c.Next()
				})
				r.PUT(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPut, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String()
				handler := h.AddressUpdate
				body = NewBody(pkgDTO.AddressCreateRequest{
					Detail:        "Jl. Abubakar Ibrahim No. 29E",
					SubDistrictID: 1,
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.ID.String(),
						},
					}
					c.Next()
				})
				r.PUT(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPut, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"Name is required.","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockAddressService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				AddressService: tt.fields.addressService,
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

func TestHandler_AddressDelete(t *testing.T) {

	model := &models.Address{
		ID: uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"),
	}

	type fields struct {
		addressService *services.MockAddressService
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
				f.addressService.On("DeleteProfileAddress", mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String()
				handler := h.AddressDelete
				r.Use(func(c *gin.Context) {
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.ID.String(),
						},
					}
					c.Next()
				})
				r.DELETE(endpoint, handler)
				req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"a66de3fd-34ca-44c2-b283-3a25bd835685"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {

			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/"
				handler := h.AddressDelete
				r.DELETE(endpoint, handler)
				req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("DeleteProfileAddress", mock.Anything).Return(nil, e.ErrDeleteDefaultAddressConflict)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String()
				handler := h.AddressDelete
				r.Use(func(c *gin.Context) {
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.ID.String(),
						},
					}
					c.Next()
				})
				r.DELETE(endpoint, handler)
				req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":409,"data":null,"error":"Conflict","message":"cannot delete default address","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockAddressService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				AddressService: tt.fields.addressService,
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

func TestHandler_AddressUpdateDefault(t *testing.T) {
	model := &models.Address{
		ID:            uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"),
		UserID:        uuid.MustParse("a110dddd-29a1-49a4-bd49-7594d024c354"),
		Name:          "Home",
		Detail:        "Jl. Jenderal Sudirman No. 1",
		SubDistrictID: 1,
		IsDefault:     true,
	}

	type fields struct {
		addressService *services.MockAddressService
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
				f.addressService.On("UpdateProfileAddressDefault", mock.Anything, mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String() + "/default"
				handler := h.AddressUpdateDefault
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.ID.String(),
						},
					}
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":"a66de3fd-34ca-44c2-b283-3a25bd835685","user_id":"a110dddd-29a1-49a4-bd49-7594d024c354","name":"Home","detail":"Jl. Jenderal Sudirman No. 1","sub_district_id":1,"is_default":true},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String() + "/default"
				handler := h.AddressUpdateDefault
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("UpdateProfileAddressDefault", mock.Anything, mock.Anything).Return(nil, e.ErrInvalidIdentifier)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.ID.String() + "/default"
				handler := h.AddressUpdateDefault
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.ID.String(),
						},
					}
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockAddressService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				AddressService: tt.fields.addressService,
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

func TestHandler_AddressShopUpdateDefault(t *testing.T) {
	model := &models.Merchant{
		ID:        1,
		UserID:    uuid.MustParse("a110dddd-29a1-49a4-bd49-7594d024c354"),
		Name:      "E-Shop",
		AddressID: uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"),
	}

	type fields struct {
		addressService *services.MockAddressService
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
				f.addressService.On("UpdateProfileAddressShopDefault", mock.Anything, mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.AddressID.String() + "/default-shop"
				handler := h.AddressShopUpdateDefault
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.AddressID.String(),
						},
					}
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"id":1,"user_id":"a110dddd-29a1-49a4-bd49-7594d024c354","name":"E-Shop","address_id":"a66de3fd-34ca-44c2-b283-3a25bd835685","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.AddressID.String() + "/default-shop"
				handler := h.AddressShopUpdateDefault
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("UpdateProfileAddressShopDefault", mock.Anything, mock.Anything).Return(model, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.AddressID.String() + "/default-shop"
				handler := h.AddressShopUpdateDefault
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.AddressID.String(),
						},
					}
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}{"code":200,"data":{"id":1,"user_id":"a110dddd-29a1-49a4-bd49-7594d024c354","name":"E-Shop","address_id":"a66de3fd-34ca-44c2-b283-3a25bd835685","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.addressService.On("UpdateProfileAddressShopDefault", mock.Anything, mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/profile/addresses/" + model.AddressID.String() + "/default-shop"
				handler := h.AddressShopUpdateDefault
				r.Use(func(c *gin.Context) {
					c.Set("id", model.UserID)
					c.Params = []gin.Param{
						{
							Key:   "address_id",
							Value: model.AddressID.String(),
						},
					}
					c.Next()
				})
				r.PATCH(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPatch, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":500,"data":null,"error":"Internal Server Error","message":"database operation failed","metadata":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				services.NewMockAddressService(t),
			}

			h := handlers.New(&handlers.HandlerConfig{
				AddressService: tt.fields.addressService,
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
