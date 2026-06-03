package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"cart-api/internal/entity"
	"cart-api/internal/errorsx"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAddCartItem(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "postgres")
	repo := &Repository{db: sqlxDB}

	ctx := context.Background()
	cartID := int64(1)
	itemID := int64(100)
	inputItem := entity.AddCartItemRequest{
		Product: "sprite",
		Price:   1500,
	}

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedItem  entity.CartItem
		expectedError error
	}{
		{
			name: "Success - Item added successfully",
			mockBehavior: func() {
				mock.ExpectBegin()
				
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
				
				mock.ExpectQuery("SELECT COUNT(*) FROM cart_items WHERE cart_id = $1").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
				
				mock.ExpectQuery("INSERT INTO cart_items (cart_id, product, price) VALUES ($1, $2, $3) RETURNING id, cart_id, product, price").
					WithArgs(cartID, inputItem.Product, inputItem.Price).
					WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product", "price"}).AddRow(itemID, cartID, "sprite", 1500))
				
				mock.ExpectCommit()
			},
			expectedItem: entity.CartItem{
				ID:      itemID,
				CartID:  cartID,
				Product: "sprite",
				Price:   1500,
			},
			expectedError: nil,
		},
		{
			name: "Error - Cart hasn't found",
			mockBehavior: func() {
				mock.ExpectBegin()
				
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
				
				mock.ExpectRollback()
			},
			expectedItem:  entity.CartItem{},
			expectedError: errorsx.ErrCartNotFound,
		},
		{
			name: "Error - Cart item limit reached",
			mockBehavior: func() {
				mock.ExpectBegin()
				
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
				
				mock.ExpectQuery("SELECT COUNT(*) FROM cart_items WHERE cart_id = $1").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
				
				mock.ExpectRollback()
			},
			expectedItem:  entity.CartItem{},
			expectedError: errorsx.ErrCartLimitReached,
		},
		{
			name: "Error - Insert statement failed (triggers rollback)",
			mockBehavior: func() {
				mock.ExpectBegin()
				
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
				
				mock.ExpectQuery("SELECT COUNT(*) FROM cart_items WHERE cart_id = $1").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				
				mock.ExpectQuery("INSERT INTO cart_items (cart_id, product, price) VALUES ($1, $2, $3) RETURNING id, cart_id, product, price").
					WithArgs(cartID, inputItem.Product, inputItem.Price).
					WillReturnError(sql.ErrConnDone)
				
				mock.ExpectRollback()
			},
			expectedItem:  entity.CartItem{},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := repo.AddCartItem(ctx, cartID, inputItem)

			if tt.expectedError != nil {
				assert.Error(t, err)

				if errors.Is(tt.expectedError, errorsx.ErrCartNotFound) || 
				   errors.Is(tt.expectedError, errorsx.ErrCartLimitReached) {
					assert.ErrorIs(t, err, tt.expectedError)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItem, result)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err, "there are unfulfilled expectations in sqlmock")
		})
	}

}

func TestCreateCart(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "postgres")
	repo := &Repository{db: sqlxDB}

	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedCart  entity.Cart
		expectedError error
	}{
		{
			name: "Success - Cart created successfully",
			mockBehavior: func() {
				mock.ExpectQuery("INSERT INTO carts DEFAULT VALUES RETURNING id").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(42)))
			},
			expectedCart: entity.Cart{
				ID: 42,
			},
			expectedError: nil,
		},
		{
			name: "Error - Database query failed",
			mockBehavior: func() {
				mock.ExpectQuery("INSERT INTO carts DEFAULT VALUES RETURNING id").
					WillReturnError(sql.ErrConnDone)
			},
			expectedCart:  entity.Cart{},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := repo.CreateCart(ctx)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCart, result)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err, "there are unfulfilled expectations in sqlmock")
		})
	}
}

func TestRemoveCartItem(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "postgres")
	repo := &Repository{db: sqlxDB}

	ctx := context.Background()
	cartID := int64(1)
	itemID := int64(5)

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Success - Item removed from cart",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM cart_items WHERE id = $1 AND cart_id = $2)").
					WithArgs(itemID, cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("DELETE FROM cart_items WHERE id = $1 AND cart_id = $2 RETURNING id, cart_id, product, price").
					WithArgs(itemID, cartID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product", "price"}).AddRow(itemID, cartID, "potato", 1500))
			},
			expectedError: nil,
		},
		{
			name: "Error - Cart not found",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expectedError: errorsx.ErrCartNotFound,
		},
		{
			name: "Error - Cart item not found",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM cart_items WHERE id = $1 AND cart_id = $2)").
					WithArgs(itemID, cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expectedError: errorsx.ErrCartItemNotFound,
		},
		{
			name: "Error - Database failure on cart check",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Error - Database failure on delete execution",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM cart_items WHERE id = $1 AND cart_id = $2)").
					WithArgs(itemID, cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("DELETE FROM cart_items WHERE id = $1 AND cart_id = $2 RETURNING id, cart_id, product, price").
					WithArgs(itemID, cartID).
					WillReturnError(sql.ErrTxDone)
			},
			expectedError: sql.ErrTxDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			err := repo.RemoveCartItem(ctx, cartID, itemID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				
				if errors.Is(tt.expectedError, errorsx.ErrCartNotFound) || 
				   errors.Is(tt.expectedError, errorsx.ErrCartItemNotFound) {
					assert.ErrorIs(t, err, tt.expectedError)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err, "there are unfulfilled expectations in sqlmock")
		})
	}
}

func TestViewCart(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "postgres")
	repo := &Repository{db: sqlxDB}

	ctx := context.Background()
	cartID := int64(1)

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedCart  entity.Cart
		expectedError error
	}{
		{
			name: "Success - Cart with items found",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				columns := []string{"id", "cart_id", "product", "price"}
				mock.ExpectQuery("SELECT id, cart_id, product, price FROM cart_items WHERE cart_id = $1").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(int64(10), cartID, "potato", 1500).
						AddRow(int64(11), cartID, "cucumber", 800))
			},
			expectedCart: entity.Cart{
				Items: []entity.CartItem{
					{ID: 10, CartID: cartID, Product: "potato", Price: 1500},
					{ID: 11, CartID: cartID, Product: "cucumber", Price: 800},
				},
			},
			expectedError: nil,
		},
		{
			name: "Success - Empty cart found",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				columns := []string{"id", "cart_id", "product", "price"}
				mock.ExpectQuery("SELECT id, cart_id, product, price FROM cart_items WHERE cart_id = $1").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows(columns))
			},
			expectedCart:  entity.Cart{Items: nil},
			expectedError: nil,
		},
		{
			name: "Error - Cart hasn't found",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expectedCart:  entity.Cart{},
			expectedError: errorsx.ErrCartNotFound,
		},
		{
			name: "Error - Database failure on fetching items",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)").
					WithArgs(cartID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectQuery("SELECT id, cart_id, product, price FROM cart_items WHERE cart_id = $1").
					WithArgs(cartID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCart:  entity.Cart{},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := repo.ViewCart(ctx, cartID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, errorsx.ErrCartNotFound) {
					assert.ErrorIs(t, err, tt.expectedError)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCart, result)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err, "there are unfulfilled expectations in sqlmock")
		})
	}
}


