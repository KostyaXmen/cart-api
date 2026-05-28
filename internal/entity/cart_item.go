package entity

type CartItem struct {
	ID      int64   `json:"id" db:"id"`
	CartID  int64   `json:"cart_id" db:"cart_id"`
	Product string  `json:"product" db:"product"`
	Price   float64 `json:"price" db:"price"`
}
