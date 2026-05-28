package repository

import (
	"context"
	"cart-api/internal/entity"
)

func (r *Repository) AddCart(ctx context.Context) (entity.Cart, error) {
	var insertedCart entity.Cart

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return insertedCart, err
	}
	defer tx.Rollback()

	query := `INSERT INTO carts (id) VALUES ($1) RETURNING id`
	err  = tx.QueryRowxContext(ctx, query).StructScan(&insertedCart)
	if err != nil {
		return insertedCart, err
	}

	return insertedCart, tx.Commit()
}