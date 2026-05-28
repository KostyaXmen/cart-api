package entity

type Cart struct {
	ID    int64       `json:"id" db:"id"`
	Items []CartItem  `json:"items" db:"-"`
}

type CartPriceCalculation struct {
	CartID          int64   `json:"cart_id"`
	TotalPrice      float64 `json:"total_price"`
	DiscountPercent int     `json:"discount_percent"`
	FinalPrice      float64 `json:"final_price"`
}