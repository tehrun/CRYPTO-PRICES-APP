package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "reflect"
    "strings"
    "testing"

    "crypto-prices-app/internal/prices"
    "github.com/gin-gonic/gin"
)

func configurePricesService(ps *prices.PricesService, baseURL string, client *http.Client) bool {
    v := reflect.ValueOf(ps)
    if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
        return false
    }

    elem := v.Elem()
    setURL := false
    setClient := false

    for i := 0; i < elem.NumField(); i++ {
        field := elem.Field(i)
        fieldType := elem.Type().Field(i)

        if !field.CanSet() {
            continue
        }

        // Set *http.Client
        if field.Type().AssignableTo(reflect.TypeOf(&http.Client{})) {
            field.Set(reflect.ValueOf(client))
            setClient = true
            continue
        }

        // Set base URL (string field containing "url" in name)
        if field.Kind() == reflect.String && strings.Contains(strings.ToLower(fieldType.Name), "url") {
            field.SetString(baseURL)
            setURL = true
            continue
        }
    }

    // If either is not found, still allow if at least URL is set.
    // Adjust if your PricesService requires client explicitly.
    return setURL || setClient
}

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

    ps := &prices.PricesService{}
    if ok := configurePricesService(ps, upstream.URL, upstream.Client()); !ok {
        t.Skip("unable to configure PricesService via reflection; adjust test setup to your service constructor/fields")
    }

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

    ps := &prices.PricesService{}
    if ok := configurePricesService(ps, upstream.URL, upstream.Client()); !ok {
        t.Skip("unable to configure PricesService via reflection; adjust test setup to your service constructor/fields")
    }

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

    ps := &prices.PricesService{}
    if ok := configurePricesService(ps, upstream.URL, upstream.Client()); !ok {
        t.Skip("unable to configure PricesService via reflection; adjust test setup to your service constructor/fields")
    }

    router := setupRouter(ps)

    req := httptest.NewRequest(http.MethodGet, "/prices/unknown", nil)
    rec := httptest.NewRecorder()

    router.ServeHTTP(rec, req)

    if rec.Code != http.StatusNotFound {
        t.Fatalf("expected status 404, got %d", rec.Code)
    }
    if !strings.Contains(strings.ToLower(rec.Body.String()), "price not found") {
        t.Fatalf("expected response to contain 'Price not found', got %q", rec.Body.String())
    }
}