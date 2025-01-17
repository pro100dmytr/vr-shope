package repository

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	Login           string    `json:"login"`
	Name            string    `json:"name"`
	LastName        string    `json:"lastName"`
	PhoneNumber     string    `json:"phoneNumber"`
	Password        string    `json:"password"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	WalletUSDT      float64   `json:"wallet_usdt"`
	NumberPurchases int       `json:"number_purchases"`
	Salt            string    `json:"salt"`
}

type Purchase struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ProductID  uuid.UUID `json:"product_id"`
	Date       time.Time `json:"date"`
	WalletUSDT float32   `json:"wallet_usdt"`
	Cost       float32   `json:"cost"`
}

type Product struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Cost          float64   `json:"cost"`
	QuantityStock int       `json:"quantity_stock"`
	Guarantees    time.Time `json:"guarantees"`
	Country       string    `json:"country"`
	Like          int       `json:"like"`
}
