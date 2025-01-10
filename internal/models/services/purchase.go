package services

import (
	"time"
)

type Purchase struct {
	ID     uint64    `json:"id"`
	UserID uint64    `json:"user_id"`
	Cost   float64   `json:"cost"`
	Date   time.Time `json:"date"`
}
