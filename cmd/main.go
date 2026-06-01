package main

import (
	"cart-api/internal/config"
	"cart-api/internal/repository"
	"cart-api/internal/entity"
	"fmt"
	"os"
	"time"
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cart := repo.AddCart(ctx)

	fmt.Println(cart)

	insertedItem, err := repo.AddCartItem(ctx, 1, entity.AddCartItemRequest{"product2", 10.0})

	if err != nil {
		fmt.Println("Error adding item in cart:", err)
	} else {
		fmt.Println(insertedItem)
		fmt.Println("List: ")
		fmt.Println(repo.ViewCart(ctx, 1))
	}
	


}