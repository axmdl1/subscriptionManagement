package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"subscriptionManagement/internal/config"
	"subscriptionManagement/internal/handler"
	repopg "subscriptionManagement/internal/repository/postgres"
	"subscriptionManagement/internal/service"
	"subscriptionManagement/internal/storage/postgres"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	pool, err := postgres.NewPool(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	subRepo := repopg.NewSubscriptionRepository(pool)
	subService := service.NewSubscriptionService(subRepo)
	subHandler := handler.NewSubscriptionHandler(subService)

	router := gin.Default()

	handler.RegisterRoutes(router, subHandler)

	log.Println("server started on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
