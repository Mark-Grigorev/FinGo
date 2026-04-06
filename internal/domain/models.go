package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrForbidden     = errors.New("forbidden")
)

// InsufficientFundsError is returned when an expense would make a non-credit account go negative.
type InsufficientFundsError struct {
	AccountName  string    `json:"account_name"`
	Balance      float64   `json:"balance"`
	Amount       float64   `json:"amount"`
	Alternatives []Account `json:"alternatives"`
}

func (e *InsufficientFundsError) Error() string { return "insufficient funds" }

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	BaseCurrency string    `json:"base_currency"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Session struct {
	Token     string
	UserID    int64
	ExpiresAt time.Time
}

type Account struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // cash, card, savings, investment
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	AccountID     int64     `json:"account_id"`
	AccountName   string    `json:"account_name,omitempty"`
	CategoryID    *int64    `json:"category_id,omitempty"`
	CategoryName  string    `json:"category_name,omitempty"`
	CategoryColor string    `json:"category_color,omitempty"`
	Icon          string    `json:"icon,omitempty"`
	Type          string    `json:"type"` // income, expense
	Amount        float64   `json:"amount"`
	Name          string    `json:"name"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
}

type Category struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Color    string `json:"color"`
	Type     string `json:"type"` // income, expense
	IsSystem bool   `json:"is_system"`
}

type Budget struct {
	ID         int64   `json:"id"`
	UserID     int64   `json:"user_id"`
	CategoryID int64   `json:"category_id"`
	Name       string  `json:"name"`   // category name
	Color      string  `json:"color"`  // category color
	Month      string  `json:"month"`  // YYYY-MM-01
	Limit      float64 `json:"limit"`
	Spent      float64 `json:"spent"`
	Pct        float64 `json:"pct"`
}

type RecurringPayment struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	AccountID       int64     `json:"account_id"`
	CategoryID      *int64    `json:"category_id,omitempty"`
	CategoryName    string    `json:"category_name,omitempty"`
	Name            string    `json:"name"`
	Amount          float64   `json:"amount"`
	Frequency       string    `json:"frequency"` // monthly, weekly, yearly
	NextPaymentDate time.Time `json:"next_payment_date"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type PasswordReset struct {
	Token     string
	UserID    int64
	ExpiresAt time.Time
	UsedAt    *time.Time
}

type ExchangeRate struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Currency  string    `json:"currency"`
	Rate      float64   `json:"rate"`
	UpdatedAt time.Time `json:"updated_at"`
}
