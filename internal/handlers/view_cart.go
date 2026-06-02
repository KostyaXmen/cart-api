package handlers

import (
	"net/http"
	"strconv"
	"encoding/json"
)

func (c *cartHandler) ViewCart(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format: must be an integer", http.StatusBadRequest)
		return
	}

	cart, err := c.service.ViewCart(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cart.ID = int64(id)
	
	json.NewEncoder(w).Encode(cart)
	w.WriteHeader(http.StatusOK)
}