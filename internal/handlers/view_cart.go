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

	var cartViewResponse CartViewResponse
	cartViewResponse.ID = cart.ID
	for _, item := range cart.Items {
		cartViewResponse.Items = append(cartViewResponse.Items, CartItemResponse{
			ID:      item.ID,
			CartID:  item.CartID,
			Product: item.Product,
			Price:   item.Price,
		})
	}
	
	
	json.NewEncoder(w).Encode(cartViewResponse)
	w.WriteHeader(http.StatusOK)
}