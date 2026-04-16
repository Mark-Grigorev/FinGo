package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/pkg/token"
)

func newTestMaker(t *testing.T) *token.Maker {
	t.Helper()
	m, err := token.New("", time.Hour)
	require.NoError(t, err)
	return m
}

func TestAuthLogin_EmptyEmail(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Login(context.Background(), "", "password")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAuthLogin_EmptyPassword(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Login(context.Background(), "user@example.com", "")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAuthLogin_UserNotFound(t *testing.T) {
	store := &mockStore{
		getUserByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Login(context.Background(), "notexist@example.com", "password")
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAuthLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	store := &mockStore{
		getUserByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: "user@example.com", PasswordHash: string(hash)}, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Login(context.Background(), "user@example.com", "wrong")
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAuthLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	want := &domain.User{ID: 7, Email: "user@example.com", PasswordHash: string(hash)}
	store := &mockStore{
		getUserByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return want, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	user, tok, payload, err := svc.Login(context.Background(), "USER@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, int64(7), user.ID)
	assert.NotEmpty(t, tok)
	assert.Equal(t, int64(7), payload.UserID)
}

func TestAuthRegister_InvalidInput(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	cases := []struct {
		email, name, password string
	}{
		{"", "Alice", "password123"},
		{"alice@example.com", "", "password123"},
		{"alice@example.com", "Alice", "short"},
	}
	for _, c := range cases {
		_, _, _, err := svc.Register(context.Background(), c.email, c.name, c.password)
		assert.ErrorIsf(t, err, domain.ErrInvalidInput,
			"Register(%q,%q,%q)", c.email, c.name, c.password)
	}
}

func TestAuthRegister_AlreadyExists(t *testing.T) {
	store := &mockStore{
		createUserFn: func(_ context.Context, _, _, _ string) (*domain.User, error) {
			return nil, domain.ErrAlreadyExists
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Register(context.Background(), "alice@example.com", "Alice", "password123")
	require.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestAuthRegister_Success(t *testing.T) {
	want := &domain.User{ID: 5, Email: "alice@example.com", Name: "Alice"}
	store := &mockStore{
		createUserFn: func(_ context.Context, email, name, _ string) (*domain.User, error) {
			return want, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	user, tok, payload, err := svc.Register(context.Background(), "alice@example.com", "Alice", "password123")
	require.NoError(t, err)
	assert.Equal(t, int64(5), user.ID)
	assert.NotEmpty(t, tok)
	assert.Equal(t, int64(5), payload.UserID)
}

func TestAuthGetUser(t *testing.T) {
	want := &domain.User{ID: 3, Email: "bob@example.com", Name: "Bob"}
	store := &mockStore{
		getUserByIDFn: func(_ context.Context, id int64) (*domain.User, error) {
			if id != 3 {
				return nil, domain.ErrNotFound
			}
			return want, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	user, err := svc.GetUser(context.Background(), 3)
	require.NoError(t, err)
	assert.Equal(t, int64(3), user.ID)
}
