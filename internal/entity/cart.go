package entity

type Cart struct {
	ID    int64       `json:"id" db:"id"`
	Items []CartItem  `json:"items" db:"-"`
}