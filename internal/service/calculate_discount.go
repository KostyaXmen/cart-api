package service

import (
	"context"
)

func (s *service) CalculateDiscount(ctx context.Context, cartID int64) (float64, float64, float64, error) {
	viewedCart, err := s.repo.ViewCart(ctx, cartID)
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	var total float64
	for _, cartItem := range viewedCart.Items {
		total += cartItem.Price
	}
	
	if total > 5000 {
		return total, 10, total * 0.9, nil
	}

	if len(viewedCart.Items) > 3 {
		return total, 5, total * 0.95, nil
	}

	return total, 0.0, total, nil
}