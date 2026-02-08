package main

import (
	"log"

	"crypto-prices-app/internal/api"
	"crypto-prices-app/internal/config"
	"crypto-prices-app/internal/prices"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	pricesService := prices.NewPricesService(cfg.BaseURL, nil)
	h := api.NewHandler(pricesService)
	r := gin.New()
	r.GET("/prices", h.GetPrices)
	r.GET("/prices/:id", h.GetPriceByID)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
