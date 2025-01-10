package handlers

import (
	"time"
)

type User struct {
	ID              int       `json:"id"`
	Login           string    `json:"login"`
	Name            string    `json:"name"`
	LastName        string    `json:"lastName"`
	PhoneNumber     string    `json:"phoneNumber"`
	Password        string    `json:"password"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	WalletUSDT      float64   `json:"wallet_usdt"`
	NumberPurchases int       `json:"number_purchases"`
}
