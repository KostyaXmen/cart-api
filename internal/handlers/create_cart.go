package handlers

import (
	"encoding/json"
	"net/http"
)

func (c *cartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	cart, err := c.service.CreateCart(r.Context())
	if err != nil {
		http.Error(w, "Failed to create cart", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cart)
	w.WriteHeader(http.StatusCreated)
}