package service

import (
	"context"
)

const (
	tenPercentDiscountValue  = 5000
	fivePercentDiscountValue = 3
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
	
	if total > tenPercentDiscountValue {
		return total, 10, total * 0.9, nil
	}

	if len(viewedCart.Items) > fivePercentDiscountValue {
		return total, 5, total * 0.95, nil
	}

	return total, 0.0, total, nil
}