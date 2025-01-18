package models

import "time"

type Product struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	VatRate   float64   `json:"vat_rate"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
