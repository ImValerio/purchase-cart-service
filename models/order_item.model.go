package models

import "time"

type OrderItem struct {
	Id        int
	OrderId   int
	ProductId int
	Price     float64
	VAT       float64
	Quantity  int
	CreatedAt time.Time
}
