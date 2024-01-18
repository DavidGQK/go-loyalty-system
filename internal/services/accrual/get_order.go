package accrual

import "net/url"

const (
	OrderStatusRegistered = "REGISTERED"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
)

func CheckOrder(host, orderNumber string) (*OrderResponse, error) {
	fullURL, err := url.JoinPath(host, "/api/orders/", orderNumber)
	if err != nil {
		return nil, err
	}
	client := NewClient()

	return client.GetOrderRequest(fullURL)
}
