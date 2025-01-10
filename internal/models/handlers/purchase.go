package handlers

import (
	"time"
)

type Purchase struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Cost   float64   `json:"cost"`
	Date   time.Time `json:"date"`
}
