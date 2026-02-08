package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"crypto-prices-app/internal/prices"

	"github.com/gin-gonic/gin"
)

func setupRouter(ps *prices.PricesService) *gin.Engine {
	h := NewHandler(ps)
	r := gin.New()
	r.GET("/prices", h.GetPrices)
	r.GET("/prices/:id", h.GetPriceByID)
	return r
}

func newUpstreamServer() *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]float64{
			"BTC": 123.45,
		})
	})

	handler.HandleFunc("/prices/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/prices/")
		if strings.EqualFold(id, "btc") {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":    "btc",
				"price": 123.45,
			})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})

	return httptest.NewServer(handler)
}

func TestGetPrices_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := newUpstreamServer()
	defer upstream.Close()

	ps := prices.NewPricesService(upstream.URL, upstream.Client())
	router := setupRouter(ps)

	req := httptest.NewRequest(http.MethodGet, "/prices", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(strings.ToUpper(rec.Body.String()), "BTC") {
		t.Fatalf("expected response to contain BTC, got %q", rec.Body.String())
	}
}

func TestGetPriceByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := newUpstreamServer()
	defer upstream.Close()

	ps := prices.NewPricesService(upstream.URL, upstream.Client())
	router := setupRouter(ps)

	req := httptest.NewRequest(http.MethodGet, "/prices/btc", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(strings.ToLower(rec.Body.String()), "btc") {
		t.Fatalf("expected response to contain btc, got %q", rec.Body.String())
	}
}

func TestGetPriceByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := newUpstreamServer()
	defer upstream.Close()

	ps := prices.NewPricesService(upstream.URL, upstream.Client())
	router := setupRouter(ps)

	req := httptest.NewRequest(http.MethodGet, "/prices/unknown", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	if !strings.Contains(strings.ToLower(rec.Body.String()), "price not found") {
		t.Fatalf("expected response to contain 'price not found', got %q", rec.Body.String())
	}
}
