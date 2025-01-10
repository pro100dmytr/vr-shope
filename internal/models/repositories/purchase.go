package repositories

import (
	"github.com/google/uuid"
	"time"
)

type Purchase struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Cost   float64   `json:"cost"`
	Date   time.Time `json:"date"`
}
