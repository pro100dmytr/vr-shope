package services

import (
	"time"
)

type Product struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	Cost          float64   `json:"cost"`
	QuantityStock int       `json:"quantity_stock"`
	Guarantees    time.Time `json:"guarantees"`
	Country       string    `json:"country"`
	Like          int       `json:"like"`
}
