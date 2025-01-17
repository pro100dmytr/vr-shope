package models

import "time"

type User struct {
	ID              uint64  `json:"id"`
	Login           string  `json:"login"`
	Name            string  `json:"name"`
	LastName        string  `json:"lastName"`
	PhoneNumber     string  `json:"phoneNumber"`
	Password        string  `json:"password"`
	Email           string  `json:"email"`
	WalletUSDT      float64 `json:"wallet_usdt"`
	NumberPurchases int     `json:"number_purchases"`
}

type UserRequest struct {
	Login       string  `json:"login"`
	Name        string  `json:"name"`
	LastName    string  `json:"lastName"`
	PhoneNumber string  `json:"phoneNumber"`
	Password    string  `json:"password"`
	Email       string  `json:"email"`
	WalletUSDT  float64 `json:"wallet_usdt"`
}

type UserResponse struct {
	Message         string    `json:"message"`
	ID              uint64    `json:"id"`
	Login           string    `json:"login"`
	Name            string    `json:"name"`
	LastName        string    `json:"lastName"`
	PhoneNumber     string    `json:"phoneNumber"`
	Password        string    `json:"password"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"createdAt"`
	WalletUSDT      float64   `json:"wallet_usdt"`
	NumberPurchases int       `json:"number_purchases"`
}
