package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rusinov-artem/gophermart/app/dto"
)

type Client struct {
	ctx     context.Context
	address string
}

func New(ctx context.Context, address string) *Client {
	return &Client{
		ctx:     ctx,
		address: address,
	}
}

func (c *Client) GetSingleOrder(orderNr string) (dto.OrderListItem, error) {
	order := dto.OrderListItem{}
	url := fmt.Sprintf("%s/api/orders/%s", c.address, orderNr)
	req, _ := http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return order, fmt.Errorf("unable to fetch order info: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return order, fmt.Errorf("unable to fetch order info: http.code %d", resp.StatusCode)
	}

	jsonOrder := struct {
		OrderNr string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float32 `json:"accrual"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&jsonOrder)
	if err != nil {
		return order, fmt.Errorf("unable to fetch order info: %w", err)
	}

	order.OrderNr = jsonOrder.OrderNr
	order.Status = jsonOrder.Status
	order.Accrual = jsonOrder.Accrual

	return order, nil
}
