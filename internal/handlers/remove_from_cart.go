package handlers

import (
	"net/http"
	"strconv"
)

func (c *cartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format: must be an integer", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(r.PathValue("item_id"))
	if err != nil {
		http.Error(w, "Invalid ID format: must be an integer", http.StatusBadRequest)
		return
	}

	err = c.service.RemoveFromCart(r.Context(), int64(id), int64(itemID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}