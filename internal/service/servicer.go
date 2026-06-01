package service

import (
	"cart-api/internal/entity"
	"context"
)

type Service interface {
	CreateCart(ctx context.Context) entity.Cart
	AddToCart(ctx context.Context, cartID int64, item entity.AddCartItemRequest) (entity.CartItem, error)
	RemoveFromCart(ctx context.Context, cartID int64, itemID int64) error
	ViewCart(ctx context.Context, cartID int64) (entity.Cart, error)
	CalculateDiscount(ctx context.Context, cartID int64) (float64, error)
}
