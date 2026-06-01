package service

import (
	"cart-api/internal/entity"
	"context"
)

func (s *service) CreateCart(ctx context.Context) entity.Cart {
	return s.repo.CreateCart(ctx)
}