package handlers

import (
	"cart-api/internal/service"
	"net/http"
)

type CartHandler interface {
	CreateCart(w http.ResponseWriter, r *http.Request)
	AddToCart(w http.ResponseWriter, r *http.Request)
	RemoveFromCart(w http.ResponseWriter, r *http.Request)
	ViewCart(w http.ResponseWriter, r *http.Request)
	CalculateDiscount(w http.ResponseWriter, r *http.Request)
}

type cartHandler struct {
	service service.Service
}

func NewCartHandler(service service.Service) CartHandler {
	return &cartHandler{
		service: service,
	}
}

func SetupRoutes(r *http.ServeMux, h CartHandler) {
	api := "/api/v1/carts"
	r.HandleFunc("POST " + api, h.CreateCart)
	r.HandleFunc("GET " + api + "/{id}", h.ViewCart)
	r.HandleFunc("POST " + api + "/{id}/items", h.AddToCart)
	r.HandleFunc("DELETE " + api + "/{id}/items/{item_id}", h.RemoveFromCart)
	r.HandleFunc("GET " + api + "/{id}/price", h.CalculateDiscount)
}