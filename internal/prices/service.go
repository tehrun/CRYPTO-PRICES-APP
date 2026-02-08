package prices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PricesService struct {
	baseURL string
	client  *http.Client
}

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Response struct {
	Prices []Price `json:"prices"`
}

func NewPricesServiceWithBaseURL(baseURL string, client *http.Client) *PricesService {
	if client == nil {
		client = &http.Client{}
	}

	return &PricesService{baseURL: strings.TrimRight(baseURL, "/"), client: client}
}

func NewPricesService(baseURL string, client *http.Client) *PricesService {
	return NewPricesServiceWithBaseURL(baseURL, client)
}

func (s *PricesService) FetchPrices() (*Response, error) {
	if s.baseURL == "" {
		return nil, fmt.Errorf("missing base URL")
	}

	resp, err := s.client.Get(s.baseURL + "/prices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch prices: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err == nil {
		if response.Prices != nil {
			return &response, nil
		}
	}

	var priceMap map[string]float64
	if err := json.Unmarshal(body, &priceMap); err != nil {
		return nil, err
	}

	converted := Response{Prices: make([]Price, 0, len(priceMap))}
	for symbol, price := range priceMap {
		converted.Prices = append(converted.Prices, Price{Symbol: symbol, Price: price})
	}

	return &converted, nil
}

func (s *PricesService) GetPrice(symbol string) (float64, error) {
	response, err := s.FetchPrices()
	if err != nil {
		return 0, err
	}

	for _, price := range response.Prices {
		if price.Symbol == symbol {
			return price.Price, nil
		}
	}

	return 0, fmt.Errorf("price for symbol %s not found", symbol)
}
