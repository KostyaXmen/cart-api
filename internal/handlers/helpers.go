package handlers

type CartPrice struct {
	ID              int64   `json:"id"`
	Total           float64 `json:"total"`
	Discount        float64 `json:"discount"`
	DiscountedTotal float64 `json:"discounted"`
}