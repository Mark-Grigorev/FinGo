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

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
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
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	AccountID    int64     `json:"account_id"`
	CategoryID   *int64    `json:"category_id,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
	Icon         string    `json:"icon,omitempty"`
	Type         string    `json:"type"` // income, expense
	Amount       float64   `json:"amount"`
	Name         string    `json:"name"`
	Date         time.Time `json:"date"`
	CreatedAt    time.Time `json:"created_at"`
}

type Category struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Type   string `json:"type"` // income, expense
}
