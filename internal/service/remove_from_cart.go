package service

import "context"

func (s *service) RemoveFromCart(ctx context.Context, cartID int64, itemID int64) error {
	return s.repo.RemoveCartItem(ctx, cartID, itemID)
}