package main

import (
	"cart-api/internal/config"
	"cart-api/internal/handlers"
	"cart-api/internal/repository"
	"cart-api/internal/service"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"
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
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func () {
		fmt.Printf("Starting HTTP server on %s...\n", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server stopped with error: %v\n", err)
			stop <- syscall.SIGTERM
		}
	}()

	sig := <- stop

	fmt.Println(sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Server forced shutdown: ", err)
	} else {
		fmt.Println("Server good shutdown")
	}


}