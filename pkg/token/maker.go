package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"aidanwoods.dev/go-paseto"
)

const defaultDuration = 24 * time.Hour

var ErrInvalidToken = errors.New("invalid token")

// Maker создаёт и верифицирует PASETO v4 local токены.
type Maker struct {
	key      paseto.V4SymmetricKey
	duration time.Duration
}

// New создаёт Maker с симметричным ключом в hex-формате.
// Если keyHex пустой — генерирует случайный ключ (для dev).
func New(keyHex string, duration time.Duration) (*Maker, error) {
	if duration == 0 {
		duration = defaultDuration
	}

	var key paseto.V4SymmetricKey
	if keyHex == "" {
		key = paseto.NewV4SymmetricKey()
	} else {
		var err error
		key, err = paseto.V4SymmetricKeyFromHex(keyHex)
		if err != nil {
			return nil, fmt.Errorf("invalid token key: %w", err)
		}
	}
	return &Maker{key: key, duration: duration}, nil
}

// CreateToken выпускает токен для пользователя.
func (m *Maker) CreateToken(userID int64, email string) (string, *Payload, error) {
	payload := &Payload{
		UserID:    userID,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(m.duration),
	}

	t := paseto.NewToken()
	t.SetIssuedAt(payload.IssuedAt)
	t.SetExpiration(payload.ExpiredAt)
	t.SetString("user_id", strconv.FormatInt(userID, 10))
	t.SetString("email", email)

	return t.V4Encrypt(m.key, nil), payload, nil
}

// VerifyToken проверяет токен и возвращает payload.
func (m *Maker) VerifyToken(tokenStr string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	t, err := parser.ParseV4Local(m.key, tokenStr, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userIDStr, err := t.GetString("user_id")
	if err != nil {
		return nil, ErrInvalidToken
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, ErrInvalidToken
	}
	email, _ := t.GetString("email")
	exp, _ := t.GetExpiration()

	return &Payload{
		UserID:    userID,
		Email:     email,
		ExpiredAt: exp,
	}, nil
}
