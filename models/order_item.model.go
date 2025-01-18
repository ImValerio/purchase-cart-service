package models

import "time"

type OrderItem struct {
	Id        int       `json:"id"`
	OrderId   int       `json:"order_id"`
	ProductId int       `json:"product_id"`
	Price     float64   `json:"price"`
	VAT       float64   `json:"vat"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemDto struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type ItemDetail struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	VAT       float64 `json:"vat"`
}
