package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAuthVerifyToken_Valid(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	tok, payload, err := svc.tokenMaker.CreateToken(42, "user@example.com")
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	got, err := svc.VerifyToken(tok)
	if err != nil {
		t.Fatalf("VerifyToken: %v", err)
	}
	if got.UserID != payload.UserID {
		t.Errorf("UserID = %d, want %d", got.UserID, payload.UserID)
	}
}

func TestAuthUpdateProfile_InvalidInput(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	cases := []struct{ name, email string }{
		{"", "user@example.com"},
		{"Alice", ""},
		{"  ", "user@example.com"},
	}
	for _, c := range cases {
		_, err := svc.UpdateProfile(context.Background(), 1, c.name, c.email)
		if err != domain.ErrInvalidInput {
			t.Errorf("UpdateProfile(%q,%q): expected ErrInvalidInput, got %v", c.name, c.email, err)
		}
	}
}

func TestAuthUpdateProfile_Success(t *testing.T) {
	want := &domain.User{ID: 1, Name: "Bob", Email: "bob@example.com"}
	store := &mockStore{
		updateUserFn: func(_ context.Context, id int64, name, email string) (*domain.User, error) {
			return want, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	got, err := svc.UpdateProfile(context.Background(), 1, "Bob", "bob@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Bob" {
		t.Errorf("Name = %q, want Bob", got.Name)
	}
}

func TestAuthForgotPassword_EmailNotFound(t *testing.T) {
	store := &mockStore{
		getUserByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ForgotPassword(context.Background(), "nobody@example.com")
	if err == nil {
		t.Error("expected error for unknown email, got nil")
	}
}

func TestAuthForgotPassword_Success(t *testing.T) {
	store := &mockStore{
		getUserByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: "user@example.com"}, nil
		},
		createPasswordResetFn: func(_ context.Context, _ string, _ int64, _ time.Time) error {
			return nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	if err := svc.ForgotPassword(context.Background(), "user@example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthResetPassword_ShortPassword(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	err := svc.ResetPassword(context.Background(), "token", "abc")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAuthResetPassword_TokenNotFound(t *testing.T) {
	store := &mockStore{
		getPasswordResetFn: func(_ context.Context, _ string) (*domain.PasswordReset, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ResetPassword(context.Background(), "badtoken", "password123")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAuthResetPassword_Expired(t *testing.T) {
	store := &mockStore{
		getPasswordResetFn: func(_ context.Context, _ string) (*domain.PasswordReset, error) {
			return &domain.PasswordReset{
				Token:     "token",
				UserID:    1,
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			}, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ResetPassword(context.Background(), "token", "password123")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput for expired token, got %v", err)
	}
}

func TestAuthResetPassword_AlreadyUsed(t *testing.T) {
	usedAt := time.Now().Add(-5 * time.Minute)
	store := &mockStore{
		getPasswordResetFn: func(_ context.Context, _ string) (*domain.PasswordReset, error) {
			return &domain.PasswordReset{
				Token:     "token",
				UserID:    1,
				ExpiresAt: time.Now().Add(10 * time.Minute),
				UsedAt:    &usedAt,
			}, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ResetPassword(context.Background(), "token", "password123")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput for used token, got %v", err)
	}
}

func TestAuthResetPassword_Success(t *testing.T) {
	store := &mockStore{
		getPasswordResetFn: func(_ context.Context, _ string) (*domain.PasswordReset, error) {
			return &domain.PasswordReset{
				Token:     "token",
				UserID:    1,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}, nil
		},
		updatePasswordFn: func(_ context.Context, _ int64, _ string) error {
			return nil
		},
		markPasswordResetUsedFn: func(_ context.Context, _ string) error {
			return nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	if err := svc.ResetPassword(context.Background(), "token", "newpassword"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthChangePassword_ShortNewPassword(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	err := svc.ChangePassword(context.Background(), 1, "oldpassword", "abc")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAuthChangePassword_WrongOldPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	store := &mockStore{
		getUserByIDFn: func(_ context.Context, _ int64) (*domain.User, error) {
			return &domain.User{ID: 1, PasswordHash: string(hash)}, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ChangePassword(context.Background(), 1, "wrong", "newpassword")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthChangePassword_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.MinCost)
	store := &mockStore{
		getUserByIDFn: func(_ context.Context, _ int64) (*domain.User, error) {
			return &domain.User{ID: 1, PasswordHash: string(hash)}, nil
		},
		updatePasswordFn: func(_ context.Context, _ int64, _ string) error {
			return nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	if err := svc.ChangePassword(context.Background(), 1, "oldpassword", "newpassword"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
