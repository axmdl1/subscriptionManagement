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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "subscriptionManagement/docs"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("server started on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
