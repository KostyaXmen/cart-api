package service

import (
	"cart-api/internal/entity"
	"context"
)

func (s *service) CreateCart(ctx context.Context) (entity.Cart, error) {
	return s.repo.CreateCart(ctx)
}