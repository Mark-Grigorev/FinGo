package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/pkg/token"
)

func newTestMaker(t *testing.T) *token.Maker {
	t.Helper()
	m, err := token.New("", time.Hour)
	if err != nil {
		t.Fatalf("token.New: %v", err)
	}
	return m
}

func TestAuthLogin_EmptyEmail(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Login(context.Background(), "", "password")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAuthLogin_EmptyPassword(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Login(context.Background(), "user@example.com", "")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAuthLogin_UserNotFound(t *testing.T) {
	store := &mockStore{
		getUserByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	_, _, _, err := svc.Login(context.Background(), "notexist@example.com", "password")
	if err != domain.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
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
	if err != domain.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 7 {
		t.Errorf("user.ID = %d, want 7", user.ID)
	}
	if tok == "" {
		t.Error("expected non-empty token string")
	}
	if payload.UserID != 7 {
		t.Errorf("payload.UserID = %d, want 7", payload.UserID)
	}
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
		if err != domain.ErrInvalidInput {
			t.Errorf("Register(%q,%q,%q): expected ErrInvalidInput, got %v", c.email, c.name, c.password, err)
		}
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
	if err != domain.ErrAlreadyExists {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 5 {
		t.Errorf("user.ID = %d, want 5", user.ID)
	}
	if tok == "" {
		t.Error("expected non-empty token")
	}
	if payload.UserID != 5 {
		t.Errorf("payload.UserID = %d, want 5", payload.UserID)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 3 {
		t.Errorf("user.ID = %d, want 3", user.ID)
	}
}
