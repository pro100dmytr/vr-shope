package dto

import (
	"time"
)

type PurchaseRequest struct {
	UserID int     `json:"user_id"`
	Cost   float64 `json:"cost"`
}

type PurchaseResponse struct {
	Message string    `json:"message"`
	ID      uint64    `json:"id"`
	UserID  uint64    `json:"user_id"`
	Cost    float64   `json:"cost"`
	Date    time.Time `json:"date"`
}
