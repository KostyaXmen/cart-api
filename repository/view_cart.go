package repository

import (
	"cart-api/internal/entity"
	"cart-api/internal/errorsx"
	"context"
)

func (r *Repository) ViewCart(ctx context.Context, cartID int64) (entity.Cart, error) {
	var viewedCart entity.Cart

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return viewedCart, err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)", cartID)
	if err != nil {
		return viewedCart, err
	}
	if !exists {
		return viewedCart, errorsx.ErrCartNotFound
	}

	query := `SELECT id, cart_id, product, price FROM cart_items WHERE cart_id = $1`
	rows, err := tx.QueryxContext(ctx, query, cartID)
	if err != nil {
		return viewedCart, err
	}
	defer rows.Close()

	var cartItems []entity.CartItem
	for rows.Next() {
		var cartItem entity.CartItem
		err := rows.Scan(&cartItem.ID, &cartItem.CartID, &cartItem.Product, &cartItem.Price)
		if err != nil {
			return viewedCart, err
		}
		cartItems = append(cartItems, cartItem)
	}

	viewedCart.Items = cartItems

	return viewedCart, nil
}