package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "crypto-prices-app/internal/prices"
)

type Handler struct {
    PricesService *prices.PricesService
}

func NewHandler(pricesService *prices.PricesService) *Handler {
    return &Handler{PricesService: pricesService}
}

func (h *Handler) GetPrices(c *gin.Context) {
    prices, err := h.PricesService.FetchPrices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, prices)
}

func (h *Handler) GetPriceByID(c *gin.Context) {
    id := c.Param("id")
    price, err := h.PricesService.GetPrice(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Price not found"})
        return
    }
    c.JSON(http.StatusOK, price)
}