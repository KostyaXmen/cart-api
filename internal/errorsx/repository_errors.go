package errorsx

import (
	"errors"
)

var (
	ErrCartNotFound     = errors.New("cart not found")
	ErrCartItemNotFound = errors.New("cart item not found")
	ErrCartLimitReached = errors.New("cart has reached the maximum capacity of 5 distinct products")
	ErrInvalidProduct   = errors.New("product name cannot be blank")
	ErrInvalidPrice     = errors.New("price must be greater than zero")
)