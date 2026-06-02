package handlers

import (
	"net/http"
	"strconv"
	"encoding/json"
)

func (c *cartHandler) CalculateDiscount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format: must be an integer", http.StatusBadRequest)
		return
	}

	total, discount, discountedTotal, err := c.service.CalculateDiscount(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CartPriceResponse{
		ID:              int64(id),
		Total:           total,
		Discount:        discount,
		DiscountedTotal: discountedTotal,
	})
	w.WriteHeader(http.StatusOK)
}