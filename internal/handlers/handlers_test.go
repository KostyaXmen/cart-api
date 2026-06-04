package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cart-api/internal/entity"
	"cart-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCart(t *testing.T) {
	var mockSvc *mocks.Service

	tests := []struct {
		name           string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Cart created",
			mockBehavior: func() {
				mockSvc.On("CreateCart", mock.Anything).
					Return(entity.Cart{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"items":null}`,
		},
		{
			name: "Error - Service failure",
			mockBehavior: func() {
				mockSvc.On("CreateCart", mock.Anything).
					Return(entity.Cart{}, errors.New("internal service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to create cart\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc = mocks.NewService(t)
			h := NewCartHandler(mockSvc)

			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodPost, "/api/v1/carts", nil)
			w := httptest.NewRecorder()

			h.CreateCart(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if strings.HasSuffix(tt.expectedBody, "\n") {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestAddToCart(t *testing.T) {
	var mockSvc *mocks.Service

	tests := []struct {
		name           string
		idPath         string
		body           string
		mockBehavior   func()
		expectedStatus int
	}{
		{
			name:   "Success - Item added",
			idPath: "1",
			body:   `{"product":"milk","price":150}`,
			mockBehavior: func() {
				reqItem := entity.AddCartItemRequest{Product: "milk", Price: 150}
				resItem := entity.CartItem{ID: 10, CartID: 1, Product: "milk", Price: 150}
				mockSvc.On("AddToCart", mock.Anything, int64(1), reqItem).
					Return(resItem, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Invalid ID format",
			idPath:         "invalid-id",
			body:           `{"product":"crocodile","price":150}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - Invalid JSON body",
			idPath:         "1",
			body:           `{invalid-json}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error - Service failure",
			idPath: "1",
			body:   `{"product":"elephant","price":150}`,
			mockBehavior: func() {
				mockSvc.On("AddToCart", mock.Anything, mock.Anything, mock.Anything).
					Return(entity.CartItem{}, errors.New("limit reached"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc = mocks.NewService(t)
			h := NewCartHandler(mockSvc)

			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodPost, "/api/v1/carts/"+tt.idPath+"/items", bytes.NewBufferString(tt.body))
			req.SetPathValue("id", tt.idPath)
			w := httptest.NewRecorder()

			h.AddToCart(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRemoveFromCart(t *testing.T) {
	var mockSvc *mocks.Service

	tests := []struct {
		name           string
		idPath         string
		itemIdPath     string
		mockBehavior   func()
		expectedStatus int
	}{
		{
			name:       "Success - Item removed",
			idPath:     "1",
			itemIdPath: "42",
			mockBehavior: func() {
				mockSvc.On("RemoveFromCart", mock.Anything, int64(1), int64(42)).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Invalid Cart ID",
			idPath:         "bimbimbim",
			itemIdPath:     "42",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - Invalid Item ID",
			idPath:         "1",
			itemIdPath:     "bambambam",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "Error - Service failure",
			idPath:     "1",
			itemIdPath: "42",
			mockBehavior: func() {
				mockSvc.On("RemoveFromCart", mock.Anything, int64(1), int64(42)).
					Return(errors.New("not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc = mocks.NewService(t)
			h := NewCartHandler(mockSvc)

			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/carts/"+tt.idPath+"/items/"+tt.itemIdPath, nil)
			req.SetPathValue("id", tt.idPath)
			req.SetPathValue("item_id", tt.itemIdPath)
			w := httptest.NewRecorder()

			h.RemoveFromCart(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestViewCart(t *testing.T) {
	var mockSvc *mocks.Service

	tests := []struct {
		name           string
		idPath         string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "Success - Cart found",
			idPath: 	  "1",
			mockBehavior: func() {
				cartData := entity.Cart{
					ID: 1,
					Items: []entity.CartItem{
						{ID: 10, CartID: 1, Product: "water", Price: 50},
					},
				}
				mockSvc.On("ViewCart", mock.Anything, int64(1)).
					Return(cartData, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"items":[{"id":10,"cart_id":1,"product":"water","price":50}]}`,
		},
		{
			name:           "Error - Invalid ID format",
			idPath:         "lol",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error - Service failure",
			idPath: "1",
			mockBehavior: func() {
				mockSvc.On("ViewCart", mock.Anything, int64(1)).
					Return(entity.Cart{}, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc = mocks.NewService(t)
			h := NewCartHandler(mockSvc)

			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/carts/"+tt.idPath, nil)
			req.SetPathValue("id", tt.idPath)
			w := httptest.NewRecorder()

			h.ViewCart(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestCalculateDiscount(t *testing.T) {
	var mockSvc *mocks.Service

	tests := []struct {
		name           string
		idPath         string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success - Discount calculated",
			idPath: "1",
			mockBehavior: func() {
				mockSvc.On("CalculateDiscount", mock.Anything, int64(1)).
					Return(100.0, 10.0, 90.0, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"total":100,"discount":10,"discounted":90}`,
		},
		{
			name:           "Error - Invalid ID format",
			idPath:         "asdasd",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error - Service failure",
			idPath: "1",
			mockBehavior: func() {
				mockSvc.On("CalculateDiscount", mock.Anything, int64(1)).
					Return(0.0, 0.0, 0.0, errors.New("calculation error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc = mocks.NewService(t)
			h := NewCartHandler(mockSvc)

			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/carts/"+tt.idPath+"/price", nil)
			req.SetPathValue("id", tt.idPath)
			w := httptest.NewRecorder()

			h.CalculateDiscount(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
