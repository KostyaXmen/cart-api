package service

import (
	// "fmt"
	"cart-api/internal/entity"
	"context"
)

func (s *service) AddToCart(ctx context.Context, cartID int64, item entity.AddCartItemRequest) (entity.CartItem, error) {
	return s.repo.AddCartItem(ctx, cartID, item)
}