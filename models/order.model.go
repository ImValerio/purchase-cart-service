package models

import "time"

type Order struct {
	Id        int
	CreatedAt time.Time
}
