package handlers

type CartPriceResponse struct {
	ID              int64   `json:"id"`
	Total           float64 `json:"total"`
	Discount        float64 `json:"discount"`
	DiscountedTotal float64 `json:"discounted"`
}

type CartItemResponse struct {
	ID      int64   `json:"id" db:"id"`
	CartID  int64   `json:"cart_id" db:"cart_id"`
	Product string  `json:"product" db:"product"`
	Price   float64 `json:"price" db:"price"`
}

type CartViewResponse struct {
	ID    int64 	         `json:"id"`
	Items []CartItemResponse `json:"items"`
}