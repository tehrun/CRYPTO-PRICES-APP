package client

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.httpClient.Get(url)
}

func (c *Client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return c.httpClient.Post(url, contentType, body)
}
