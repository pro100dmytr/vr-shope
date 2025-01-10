package repositories

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Cost          float64   `json:"cost"`
	QuantityStock int       `json:"quantity_stock"`
	Guarantees    time.Time `json:"guarantees"`
	Country       string    `json:"country"`
	Like          int       `json:"like"`
}
