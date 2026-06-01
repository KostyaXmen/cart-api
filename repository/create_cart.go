package repository

import (
	"context"
	"cart-api/internal/entity"
)

func (r *Repository) AddCart(ctx context.Context) (entity.Cart, error) {
	var insertedCart entity.Cart

	query := `INSERT INTO carts (id) VALUES ($1) RETURNING id`
	err := r.db.QueryRowxContext(ctx, query).StructScan(&insertedCart)
	if err != nil {
		return insertedCart, err
	}

	return insertedCart, nil
}