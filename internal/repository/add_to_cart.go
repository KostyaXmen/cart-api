package repository

import (
	"context"
	"cart-api/internal/entity"
	"cart-api/internal/errorsx"
)

const (
	maxCartItems = 5
	minPrice     = 0
)

func (r *Repository) AddCartItem(ctx context.Context, cartID int64, item entity.AddCartItemRequest) (entity.CartItem, error) {
	var insertedItem entity.CartItem

	if item.Price < minPrice {
		return insertedItem, errorsx.ErrInvalidPrice
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return insertedItem, err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)", cartID)
	if err != nil {
		return insertedItem, err
	}
	if !exists {
		return insertedItem, errorsx.ErrCartNotFound
	}

	var count int
	err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		return insertedItem, err
	}
	if count >= maxCartItems {
		return insertedItem, errorsx.ErrCartLimitReached
	}

	query := `INSERT INTO cart_items (cart_id, product, price) VALUES ($1, $2, $3) RETURNING id, cart_id, product, price`
	err = tx.QueryRowxContext(ctx, query, cartID, item.Product, item.Price).StructScan(&insertedItem)
	if err != nil {
		return insertedItem, err
	}

	return insertedItem, tx.Commit()
}