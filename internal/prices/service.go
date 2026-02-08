package prices

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PricesService struct {
	client *http.Client
}

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Response struct {
	Prices []Price `json:"prices"`
}

func NewPricesService(client *http.Client) *PricesService {
	return &PricesService{client: client}
}

func (s *PricesService) FetchPrices(apiURL string) (*Response, error) {
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch prices: %s", resp.Status)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *PricesService) GetPrice(symbol string, apiURL string) (float64, error) {
	response, err := s.FetchPrices(apiURL)
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
