package types

type Price struct {
    Currency string  `json:"currency"`
    Amount   float64 `json:"amount"`
}

type Response struct {
    Success bool   `json:"success"`
    Data    []Price `json:"data"`
    Error   string  `json:"error,omitempty"`
}