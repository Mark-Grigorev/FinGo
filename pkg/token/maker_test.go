package token

import (
	"testing"
	"time"
)

func TestNew_RandomKey(t *testing.T) {
	m, err := New("", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil maker")
	}
}

func TestNew_ValidHexKey(t *testing.T) {
	// 32 bytes → 64 hex chars
	key := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	m, err := New(key, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil maker")
	}
}	

func TestNew_InvalidHexKey(t *testing.T) {
	_, err := New("not-valid-hex", time.Hour)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestNew_ZeroDuration_UsesDefault(t *testing.T) {
	m, err := New("", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.duration != defaultDuration {
		t.Errorf("expected duration %v, got %v", defaultDuration, m.duration)
	}
}

func TestCreateAndVerifyToken(t *testing.T) {
	m, _ := New("", time.Hour)

	tokenStr, payload, err := m.CreateToken(42, "user@example.com")
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}
	if payload.UserID != 42 {
		t.Errorf("payload.UserID = %d, want 42", payload.UserID)
	}
	if payload.Email != "user@example.com" {
		t.Errorf("payload.Email = %q, want %q", payload.Email, "user@example.com")
	}

	verified, err := m.VerifyToken(tokenStr)
	if err != nil {
		t.Fatalf("VerifyToken: %v", err)
	}
	if verified.UserID != 42 {
		t.Errorf("verified.UserID = %d, want 42", verified.UserID)
	}
	if verified.Email != "user@example.com" {
		t.Errorf("verified.Email = %q, want %q", verified.Email, "user@example.com")
	}
}

func TestVerifyToken_Expired(t *testing.T) {
	m, _ := New("", -time.Second)

	tokenStr, _, err := m.CreateToken(1, "user@example.com")
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	_, err = m.VerifyToken(tokenStr)
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestVerifyToken_InvalidString(t *testing.T) {
	m, _ := New("", time.Hour)

	_, err := m.VerifyToken("not-a-valid-token")
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestVerifyToken_WrongKey(t *testing.T) {
	m1, _ := New("", time.Hour)
	m2, _ := New("", time.Hour)

	tokenStr, _, _ := m1.CreateToken(1, "user@example.com")

	_, err := m2.VerifyToken(tokenStr)
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
