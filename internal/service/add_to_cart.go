package service

import (
	// "fmt"
	"cart-api/internal/entity"
	"context"
)

func (s *service) AddToCart(ctx context.Context, cartID int64, item entity.AddCartItemRequest) (entity.CartItem, error) {
	return s.repo.AddCartItem(ctx, cartID, item)
}

// func (s *service) AddToCart(ctx context.Context, cartID int64, item entity.AddCartItemRequest) (entity.CartItem, error) {
// 	cartItem, err := s.repo.AddCartItem(ctx, cartID, item)
// 	if err != nil {
// 		return entity.CartItem{}, fmt.Errorf("failed to add item to the cart")
// 	}
// 	return cartItem, nil
// }

