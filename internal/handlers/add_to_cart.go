package handlers

import (
	"cart-api/internal/entity"

	"net/http"
	"strconv"
	"encoding/json"
)

func (c *cartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format: must be an integer", http.StatusBadRequest)
		return
	}

	var item entity.AddCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addedItem, err := c.service.AddToCart(r.Context(), int64(id), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(addedItem)
	w.WriteHeader(http.StatusCreated)
}