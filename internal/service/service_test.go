package service

import (
	"context"
	"errors"
	"testing"

	"cart-api/internal/entity"
	"cart-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateCart(t *testing.T) {
	ctx := context.Background()
	expectedCart := entity.Cart{ID: 1}

	var mockRepo *mocks.Repository

	tests := []struct {
		name         string
		mockBehavior func()
		expectedCart entity.Cart
		expectedErr  error
	}{
		{
			name: "Success - Cart created successfully",
			mockBehavior: func() {
				mockRepo.On("CreateCart", ctx).Return(expectedCart, nil)
			},
			expectedCart: expectedCart,
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = mocks.NewRepository(t)
			svc := NewService(mockRepo)

			tt.mockBehavior()

			res, err := svc.CreateCart(ctx)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCart, res)
			}
		})
	}
}

func TestAddToCart(t *testing.T) {
	ctx := context.Background()
	req := entity.AddCartItemRequest{Product: "item", Price: 100}
	expectedItem := entity.CartItem{ID: 1, CartID: 1, Product: "item", Price: 100}

	var mockRepo *mocks.Repository

	tests := []struct {
		name         string
		mockBehavior func()
		expectedItem entity.CartItem
		expectedErr  error
	}{
		{
			name: "Success - Item added successfully",
			mockBehavior: func() {
				mockRepo.On("AddCartItem", ctx, int64(1), req).Return(expectedItem, nil)
			},
			expectedItem: expectedItem,
			expectedErr:  nil,
		},
	}

	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = mocks.NewRepository(t)
			svc := NewService(mockRepo)

			tt.mockBehavior()

			res, err := svc.AddToCart(ctx, 1, req)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItem, res)
			}
		})
	}
}

func TestRemoveFromCart(t *testing.T) {
	ctx := context.Background()

	var mockRepo *mocks.Repository

	tests := []struct {
		name         string
		mockBehavior func()
		expectedErr  error
	}{
		{
			name: "Success - Item removed successfully",
			mockBehavior: func() {
				mockRepo.On("RemoveCartItem", ctx, int64(1), int64(2)).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = mocks.NewRepository(t)
			svc := NewService(mockRepo)

			tt.mockBehavior()

			err := svc.RemoveFromCart(ctx, 1, 2)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestViewCart(t *testing.T) {
	ctx := context.Background()
	expectedCart := entity.Cart{ID: 1}

	var mockRepo *mocks.Repository

	tests := []struct {
		name         string
		mockBehavior func()
		expectedCart entity.Cart
		expectedErr  error
	}{
		{
			name: "Success - Cart viewed successfully",
			mockBehavior: func() {
				mockRepo.On("ViewCart", ctx, int64(1)).Return(expectedCart, nil)
			},
			expectedCart: expectedCart,
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = mocks.NewRepository(t)
			svc := NewService(mockRepo)

			tt.mockBehavior()

			res, err := svc.ViewCart(ctx, 1)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCart, res)
			}
		})
	}
}

func TestCalculateDiscount(t *testing.T) {
	ctx := context.Background()

	var mockRepo *mocks.Repository

	tests := []struct {
		name             string
		mockBehavior     func()
		expectedTotal    float64
		expectedDiscount float64
		expectedFinal    float64
		expectedErr      error
	}{
		{
			name: "Error - Repository failure",
			mockBehavior: func() {
				mockRepo.On("ViewCart", ctx, int64(1)).Return(entity.Cart{}, errors.New("db error"))
			},
			expectedTotal:    0.0,
			expectedDiscount: 0.0,
			expectedFinal:    0.0,
			expectedErr:      errors.New("db error"),
		},
		{
			name: "Discount 10 percent - High total price",
			mockBehavior: func() {
				cart := entity.Cart{
					Items: []entity.CartItem{
						{Price: 6000},
					},
				}
				mockRepo.On("ViewCart", ctx, int64(1)).Return(cart, nil)
			},
			expectedTotal:    6000.0,
			expectedDiscount: 10.0,
			expectedFinal:    5400.0,
			expectedErr:      nil,
		},
		{
			name: "Discount 5 percent - More than 3 items",
			mockBehavior: func() {
				cart := entity.Cart{
					Items: []entity.CartItem{
						{Price: 100},
						{Price: 100},
						{Price: 100},
						{Price: 100},
						{Price: 100},
					},
				}
				mockRepo.On("ViewCart", ctx, int64(1)).Return(cart, nil)
			},
			expectedTotal:    500.0,
			expectedDiscount: 5.0,
			expectedFinal:    475.0,
			expectedErr:      nil,
		},
		{
			name: "No discount - Low price and few items",
			mockBehavior: func() {
				cart := entity.Cart{
					Items: []entity.CartItem{
						{Price: 100},
					},
				}
				mockRepo.On("ViewCart", ctx, int64(1)).Return(cart, nil)
			},
			expectedTotal:    100.0,
			expectedDiscount: 0.0,
			expectedFinal:    100.0,
			expectedErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = mocks.NewRepository(t)
			svc := NewService(mockRepo)

			tt.mockBehavior()

			total, disc, final, err := svc.CalculateDiscount(ctx, 1)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, total)
				assert.Equal(t, tt.expectedDiscount, disc)
				assert.Equal(t, tt.expectedFinal, final)
			}
		})
	}
}