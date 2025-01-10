package repositories

import (
	"github.com/google/uuid"
	"time"
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
