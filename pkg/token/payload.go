package token

import "time"

// Payload содержит данные внутри PASETO токена.
type Payload struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
