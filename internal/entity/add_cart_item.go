package entity

type AddCartItemRequest struct {
	Product string  `json:"product"`
	Price   float64 `json:"price"`
}