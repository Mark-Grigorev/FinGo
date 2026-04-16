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
)

func TestAuthVerifyToken_Valid(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	tok, payload, err := svc.tokenMaker.CreateToken(42, "user@example.com")
	require.NoError(t, err)

	got, err := svc.VerifyToken(tok)
	require.NoError(t, err)
	assert.Equal(t, payload.UserID, got.UserID)
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
		assert.ErrorIsf(t, err, domain.ErrInvalidInput,
			"UpdateProfile(%q,%q)", c.name, c.email)
	}
}

func TestAuthUpdateProfile_Success(t *testing.T) {
	want := &domain.User{ID: 1, Name: "Bob", Email: "bob@example.com"}
	store := &mockStore{
		updateUserFn: func(_ context.Context, _ int64, _, _ string) (*domain.User, error) {
			return want, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	got, err := svc.UpdateProfile(context.Background(), 1, "Bob", "bob@example.com")
	require.NoError(t, err)
	assert.Equal(t, "Bob", got.Name)
}

func TestAuthForgotPassword_EmailNotFound(t *testing.T) {
	store := &mockStore{
		getUserByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ForgotPassword(context.Background(), "nobody@example.com")
	assert.Error(t, err)
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
	require.NoError(t, svc.ForgotPassword(context.Background(), "user@example.com"))
}

func TestAuthResetPassword_ShortPassword(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	err := svc.ResetPassword(context.Background(), "token", "abc")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAuthResetPassword_TokenNotFound(t *testing.T) {
	store := &mockStore{
		getPasswordResetFn: func(_ context.Context, _ string) (*domain.PasswordReset, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ResetPassword(context.Background(), "badtoken", "password123")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAuthResetPassword_Expired(t *testing.T) {
	store := &mockStore{
		getPasswordResetFn: func(_ context.Context, _ string) (*domain.PasswordReset, error) {
			return &domain.PasswordReset{
				Token:     "token",
				UserID:    1,
				ExpiresAt: time.Now().Add(-time.Hour),
			}, nil
		},
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	err := svc.ResetPassword(context.Background(), "token", "password123")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
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
	require.ErrorIs(t, err, domain.ErrInvalidInput)
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
		updatePasswordFn: func(_ context.Context, _ int64, _ string) error { return nil },
		markPasswordResetUsedFn: func(_ context.Context, _ string) error { return nil },
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	require.NoError(t, svc.ResetPassword(context.Background(), "token", "newpassword"))
}

func TestAuthChangePassword_ShortNewPassword(t *testing.T) {
	svc := NewAuth(&mockStore{}, newTestMaker(t), nil, "", slog.Default())
	err := svc.ChangePassword(context.Background(), 1, "oldpassword", "abc")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
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
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAuthChangePassword_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.MinCost)
	store := &mockStore{
		getUserByIDFn: func(_ context.Context, _ int64) (*domain.User, error) {
			return &domain.User{ID: 1, PasswordHash: string(hash)}, nil
		},
		updatePasswordFn: func(_ context.Context, _ int64, _ string) error { return nil },
	}
	svc := NewAuth(store, newTestMaker(t), nil, "", slog.Default())
	require.NoError(t, svc.ChangePassword(context.Background(), 1, "oldpassword", "newpassword"))
}
