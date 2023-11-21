package accrual

import (
	"encoding/json"
	"fmt"
	"github.com/DavidGQK/go-loyalty-system/internal/logger"
	"github.com/sethgrid/pester"
	"net/http"
)

type Client struct {
	HTTPClient *pester.Client
}

func NewClient() *Client {
	client := pester.New()
	client.MaxRetries = 3
	client.KeepLog = true

	return &Client{
		HTTPClient: client,
	}
}

func (c *Client) GetOrderRequest(url string) (*OrderResponse, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		logger.Log.Errorf("accrual request error: %v", c.HTTPClient.LogString())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status in accrual response: %v", resp.Status)
	}

	r := OrderResponse{
		Accrual: 0,
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
