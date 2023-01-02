package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ChatGetAllChatList(t *testing.T) {
	chats := []*internalDtos.ChatList{
		{
			ToUserID:        "f5677e10-81a1-4b23-9190-a72388070c15",
			ToUserPhoto:     "",
			ToUserFullName:  "Budi Suherman",
			ToMerchantID:    2,
			ToMerchantName:  "Toko Budi",
			ToMerchantPhoto: "https://res.cloudinary.com/dnjy6s5li/image/upload/v1672655864/oceana/dnlutfjcotlesfquvzos.jpg",
			LastMessage:     "pro player",
		},
	}

	type fields struct {
		chatService *services.MockChatService
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
				f.chatService.On("GetAllChatList", mock.Anything).Return(chats, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats"
				handler := h.ChatGetAllChatList
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":[{"to_user_id":"f5677e10-81a1-4b23-9190-a72388070c15","to_user_photo":"","to_user_full_name":"Budi Suherman","to_merchant_id":2,"to_merchant_name":"Toko Budi","to_merchant_photo":"https://res.cloudinary.com/dnjy6s5li/image/upload/v1672655864/oceana/dnlutfjcotlesfquvzos.jpg","last_message":"pro player","created_at":"0001-01-01T00:00:00Z"}],"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats"
				handler := h.ChatGetAllChatList
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
				f.chatService.On("GetAllChatList", mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats"
				handler := h.ChatGetAllChatList
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
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
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockChatService(t),
				}

				h := handlers.New(&handlers.HandlerConfig{
					ChatService: tt.fields.chatService,
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

func TestHandler_ChatGetChat(t *testing.T) {
	chat := &internalDtos.ChatResponse{
		ToUserID:        "f5677e10-81a1-4b23-9190-a72388070c15",
		ToUserPhoto:     "",
		ToUserFullName:  "Budi Suherman",
		ToMerchantID:    2,
		ToMerchantName:  "Toko Budi",
		ToMerchantPhoto: "https://res.cloudinary.com/dnjy6s5li/image/upload/v1672655864/oceana/dnlutfjcotlesfquvzos.jpg",
		Chats: []*internalDtos.Chat{{
			ID:         37,
			FromUserID: uuid.MustParse("76747c91-6b8b-42b2-bbac-d76a8e2090db"),
			ToUserID:   uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
			Message:    "pro player",
			MediaURL:   "",
		}},
	}

	type fields struct {
		chatService *services.MockChatService
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
				f.chatService.On("GetChat", mock.Anything, mock.Anything).Return(chat, nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats/f5677e10-81a1-4b23-9190-a72388070c15"
				handler := h.ChatGetChat
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
					c.Params = []gin.Param{
						{
							Key:   "to-user-id",
							Value: "f5677e10-81a1-4b23-9190-a72388070c15",
						},
					}
					c.Next()
				})
				r.GET(endpoint, handler)
				req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":{"to_user_id":"f5677e10-81a1-4b23-9190-a72388070c15","to_user_photo":"","to_user_full_name":"Budi Suherman","to_merchant_id":2,"to_merchant_name":"Toko Budi","to_merchant_photo":"https://res.cloudinary.com/dnjy6s5li/image/upload/v1672655864/oceana/dnlutfjcotlesfquvzos.jpg","chats":[{"id":37,"from_user_id":"76747c91-6b8b-42b2-bbac-d76a8e2090db","to_user_id":"f5677e10-81a1-4b23-9190-a72388070c15","message":"pro player","media_url":"","created_at":"0001-01-01T00:00:00Z"}]},"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats/f5677e10-81a1-4b23-9190-a72388070c15"
				handler := h.ChatGetChat
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
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
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats/f5677e10-81a1-4b23-9190-a72388070c15"
				handler := h.ChatGetChat
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.Nil)
					c.Params = []gin.Param{
						{
							Key:   "to-user-id",
							Value: "f5677e10-81a1-4b23-9190-a72388070c15",
						},
					}
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
				f.chatService.On("GetChat", mock.Anything, mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats/f5677e10-81a1-4b23-9190-a72388070c15"
				handler := h.ChatGetChat
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
					c.Params = []gin.Param{
						{
							Key:   "to-user-id",
							Value: "f5677e10-81a1-4b23-9190-a72388070c15",
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
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockChatService(t),
				}

				h := handlers.New(&handlers.HandlerConfig{
					ChatService: tt.fields.chatService,
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

func TestHandler_ChatCreate(t *testing.T) {
	type fields struct {
		chatService *services.MockChatService
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
				f.chatService.On("Create", mock.Anything).Return(nil)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats"
				handler := h.ChatCreate
				body := NewBody(internalDtos.ChatRequest{
					FromUserID: uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
					ToUserID:   uuid.MustParse("adea8ba3-9418-4124-b5ae-809e35e5986f"),
					Message:    "Permisi, apakah barangnya ready stock?",
					MediaURL:   "",
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":200,"data":true,"error":null,"message":"success","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats"
				handler := h.ChatCreate
				body := NewBody(internalDtos.ChatRequest{
					FromUserID: uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
					ToUserID:   uuid.MustParse("adea8ba3-9418-4124-b5ae-809e35e5986f"),
					Message:    "Permisi, apakah barangnya ready stock?",
					MediaURL:   "",
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"invalid identifier","metadata":null}`,
		},
		{
			mock: func(f *fields) {
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats"
				handler := h.ChatCreate
				body := NewBody(internalDtos.ChatRequest{
					FromUserID: uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
					ToUserID:   uuid.MustParse("adea8ba3-9418-4124-b5ae-809e35e5986f"),
					MediaURL:   "",
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
					c.Next()
				})
				r.POST(endpoint, handler)
				req, _ := http.NewRequest(http.MethodPost, endpoint, body)
				r.ServeHTTP(w, req)
			},
			want: `{"code":400,"data":null,"error":"Bad Request","message":"Message is required.","metadata":null}`,
		},
		{
			mock: func(f *fields) {
				f.chatService.On("Create", mock.Anything).Return(e.ErrDatabaseOperationFailed)
			},
			configureRouter: func(r *gin.Engine, h *handlers.Handler, w *httptest.ResponseRecorder) {
				endpoint := "/chats"
				handler := h.ChatCreate
				body := NewBody(internalDtos.ChatRequest{
					FromUserID: uuid.MustParse("f5677e10-81a1-4b23-9190-a72388070c15"),
					ToUserID:   uuid.MustParse("adea8ba3-9418-4124-b5ae-809e35e5986f"),
					Message:    "Permisi, apakah barangnya ready stock?",
					MediaURL:   "",
				})
				r.Use(func(c *gin.Context) {
					c.Set("id", uuid.MustParse("a66de3fd-34ca-44c2-b283-3a25bd835685"))
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
			t.Run(tt.name, func(t *testing.T) {
				tt.fields = fields{
					services.NewMockChatService(t),
				}

				h := handlers.New(&handlers.HandlerConfig{
					ChatService: tt.fields.chatService,
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
