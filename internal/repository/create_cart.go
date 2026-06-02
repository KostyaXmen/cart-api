package repository

import (
	"context"
	"cart-api/internal/entity"
)

func (r *Repository) CreateCart(ctx context.Context) (entity.Cart, error) {
	var insertedCart entity.Cart

	query := `INSERT INTO carts DEFAULT VALUES RETURNING id`
	err := r.db.QueryRowxContext(ctx, query).StructScan(&insertedCart)
	if err != nil {
		return insertedCart, err
	}

	return insertedCart, nil
}