package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CartItemDTO struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type CartDTO struct {
	UserID string        `json:"user_id"`
	Items  []CartItemDTO `json:"items"`
}

type CartResponse struct {
	Success bool    `json:"success"`
	Data    CartDTO `json:"data"`
}

type CartClient interface {
	GetCart(ctx context.Context, userID string, token string) (*CartDTO, error)
	ClearCart(ctx context.Context, userID string, token string) error
}

type cartClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCartClient(baseURL string) CartClient {
	return &cartClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *cartClient) GetCart(ctx context.Context, userID string, token string) (*CartDTO, error) {
	url := fmt.Sprintf("%s/api/v1/cart", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("NewRequestWithContext failed: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http Do failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cart service returned status %d", resp.StatusCode)
	}

	var res CartResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decode cart response failed: %w", err)
	}

	return &res.Data, nil
}

func (c *cartClient) ClearCart(ctx context.Context, userID string, token string) error {
	url := fmt.Sprintf("%s/api/v1/cart", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("NewRequestWithContext failed: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http Do failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cart service returned status %d", resp.StatusCode)
	}

	return nil
}
