package models

import "time"

type Product struct {
	Id        int
	Name      string
	VatRate   float64
	Price     float64
	CreatedAt time.Time
}
