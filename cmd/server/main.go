package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "crypto-prices-app/internal/api"
    "crypto-prices-app/internal/config"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("could not load config: %v", err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/prices", api.GetPrices).Methods("GET")
    r.HandleFunc("/prices/{id}", api.GetPriceByID).Methods("GET")

    log.Printf("Starting server on port %s", cfg.ServerPort)
    if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
        log.Fatalf("could not start server: %v", err)
    }
}