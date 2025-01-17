package models

import "time"

type Purchase struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	ProductID  uint64    `json:"product_id"`
	Date       time.Time `json:"date"`
	WalletUSDT float32   `json:"wallet_usdt"`
	Cost       float32   `json:"cost"`
}

type PurchaseRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type PurchaseResponse struct {
	Message    string    `json:"message"`
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	ProductID  uint64    `json:"product_id"`
	Date       time.Time `json:"date"`
	WalletUSDT float32   `json:"wallet_usdt"`
	Cost       float32   `json:"cost"`
}
