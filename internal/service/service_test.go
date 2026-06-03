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
	mockRepo := new(mocks.Repository)
	svc := NewService(mockRepo)

	expectedCart := entity.Cart{ID: 1}
	mockRepo.On("CreateCart", ctx).Return(expectedCart, nil)

	res, err := svc.CreateCart(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedCart, res)
}

func TestAddToCart(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.Repository)
	svc := NewService(mockRepo)

	req := entity.AddCartItemRequest{Product: "item", Price: 100}
	expectedItem := entity.CartItem{ID: 1, CartID: 1, Product: "item", Price: 100}
	mockRepo.On("AddCartItem", ctx, int64(1), req).Return(expectedItem, nil)

	res, err := svc.AddToCart(ctx, 1, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem, res)
}

func TestRemoveFromCart(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.Repository)
	svc := NewService(mockRepo)

	mockRepo.On("RemoveCartItem", ctx, int64(1), int64(2)).Return(nil)

	err := svc.RemoveFromCart(ctx, 1, 2)
	assert.NoError(t, err)
}

func TestViewCart(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.Repository)
	svc := NewService(mockRepo)

	expectedCart := entity.Cart{ID: 1}
	mockRepo.On("ViewCart", ctx, int64(1)).Return(expectedCart, nil)

	res, err := svc.ViewCart(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCart, res)
}

func TestCalculateDiscount(t *testing.T) {
	ctx := context.Background()

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(mocks.Repository)
		svc := NewService(mockRepo)

		mockRepo.On("ViewCart", ctx, int64(1)).Return(entity.Cart{}, errors.New("db error"))

		total, disc, final, err := svc.CalculateDiscount(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, 0.0, total)
		assert.Equal(t, 0.0, disc)
		assert.Equal(t, 0.0, final)
	})

	t.Run("Discount 10 percent", func(t *testing.T) {
		mockRepo := new(mocks.Repository)
		svc := NewService(mockRepo)

		cart := entity.Cart{
			Items: []entity.CartItem{
				{Price: 6000},
			},
		}
		mockRepo.On("ViewCart", ctx, int64(1)).Return(cart, nil)

		total, disc, final, err := svc.CalculateDiscount(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, 6000.0, total)
		assert.Equal(t, 10.0, disc)
		assert.Equal(t, 5400.0, final)
	})

	t.Run("Discount 5 percent", func(t *testing.T) {
		mockRepo := new(mocks.Repository)
		svc := NewService(mockRepo)

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

		total, disc, final, err := svc.CalculateDiscount(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, 500.0, total)
		assert.Equal(t, 5.0, disc)
		assert.Equal(t, 475.0, final)
	})

	t.Run("No discount", func(t *testing.T) {
		mockRepo := new(mocks.Repository)
		svc := NewService(mockRepo)

		cart := entity.Cart{
			Items: []entity.CartItem{
				{Price: 100},
			},
		}
		mockRepo.On("ViewCart", ctx, int64(1)).Return(cart, nil)

		total, disc, final, err := svc.CalculateDiscount(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, 100.0, total)
		assert.Equal(t, 0.0, disc)
		assert.Equal(t, 100.0, final)
	})
}