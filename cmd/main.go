package main

import (
	"cart-api/internal/config"
	"cart-api/internal/handlers"
	"cart-api/internal/repository"
	"cart-api/internal/service"
	"fmt"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Fatal error:", err)
		os.Exit(1)
	}

	fmt.Println(cfg)

	dbConn, err := repository.InitDB(cfg)
	if err != nil {
		fmt.Println("Fatal error:", err)
	}
	defer dbConn.Close()

	repo := repository.NewRepository(dbConn)
	svc := service.NewService(repo)
	handler := handlers.NewCartHandler(svc)

	mux := http.NewServeMux()

	handlers.SetupRoutes(mux, handler)

	serverAddr := ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		// ReadTimeout:  5 * time.Second,  
		// WriteTimeout: 10 * time.Second, 
		// IdleTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting HTTP Server on %s...\n", serverAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server stopped with error: %v\n", err)
		os.Exit(1)
	}
}