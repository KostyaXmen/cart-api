package service

import (
	"cart-api/internal/entity"
	"context"
)

func (s *service) ViewCart(ctx context.Context, cartID int64) (entity.Cart, error) {
	return s.repo.ViewCart(ctx, cartID)
}
