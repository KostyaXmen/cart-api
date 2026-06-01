package repository

import (
	"context"
	"cart-api/internal/entity"
)

func (r *Repository) CreateCart(ctx context.Context) entity.Cart {
	var insertedCart entity.Cart

	query := `INSERT INTO carts DEFAULT VALUES RETURNING id`
	err := r.db.QueryRowxContext(ctx, query).StructScan(&insertedCart)
	if err != nil {
		return insertedCart
	}

	return insertedCart
}