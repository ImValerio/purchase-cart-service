package models

import "time"

type Order struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderRequest struct {
	Order struct {
		Items []OrderItemDto `json:"items"`
	} `json:"order"`
}

type OrderResponse struct {
	OrderID    int          `json:"order_id"`
	OrderPrice float64      `json:"order_price"`
	OrderVAT   float64      `json:"order_vat"`
	Items      []ItemDetail `json:"items"`
}
