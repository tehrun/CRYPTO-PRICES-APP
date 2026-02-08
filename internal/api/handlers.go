package api

import (
	"net/http"
	"strings"

	"crypto-prices-app/internal/prices"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	prices *prices.PricesService
}

func NewHandler(ps *prices.PricesService) *Handler {
	return &Handler{prices: ps}
}

func (h *Handler) GetPrices(c *gin.Context) {
	response, err := h.prices.FetchPrices()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch prices"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetPriceByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing price id"})
		return
	}

	price, err := h.prices.GetPrice(strings.ToUpper(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "price not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    strings.ToLower(id),
		"price": price,
	})
}
