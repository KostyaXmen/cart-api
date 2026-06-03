package service

import (
	"cart-api/internal/entity"
	"context"
)

//go:generate $HOME/go/bin/mockery --name=Repository --output=./mocks
type Repository interface {
	CreateCart(ctx context.Context) (entity.Cart, error)
	AddCartItem(ctx context.Context, cartID int64, item entity.AddCartItemRequest) (entity.CartItem, error)
	RemoveCartItem(ctx context.Context, cartID int64, itemID int64) error
	ViewCart(ctx context.Context, cartID int64) (entity.Cart, error)
}