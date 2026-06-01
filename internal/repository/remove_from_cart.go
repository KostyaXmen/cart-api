package repository

import (
	"context"
	"cart-api/internal/entity"
	"cart-api/internal/errorsx"
)

func (r *Repository)RemoveCartItem(ctx context.Context, cartID int64, itemID int64) error {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)", cartID)
	if err != nil {
		return err
	}
	if !exists {
		return errorsx.ErrCartNotFound
	}

	err = r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM cart_items WHERE id = $1 AND cart_id = $2)", itemID, cartID)
	if err != nil {
		return err
	}
	if !exists {
		return errorsx.ErrCartItemNotFound
	}

	query := `DELETE FROM cart_items WHERE id = $1 AND cart_id = $2 RETURNING id, cart_id, product, price`
	err = r.db.QueryRowxContext(ctx, query, itemID, cartID).StructScan(&entity.CartItem{})
	if err != nil {
		return err
	}

	return nil

}