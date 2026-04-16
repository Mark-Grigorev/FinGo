package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestRecurringList_Success(t *testing.T) {
	want := []domain.RecurringPayment{{ID: 1, Name: "Netflix", Amount: 500}}
	store := &mockStore{
		listRecurringFn: func(_ context.Context, _ int64) ([]domain.RecurringPayment, error) {
			return want, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	got, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Netflix" {
		t.Errorf("unexpected list: %v", got)
	}
}

func TestRecurringCreate_InvalidInput(t *testing.T) {
	svc := NewRecurring(&mockStore{}, slog.Default())
	cases := []domain.RecurringPayment{
		{Name: "", Amount: 100, AccountID: 1},       // empty name
		{Name: "  ", Amount: 100, AccountID: 1},     // whitespace name
		{Name: "Sub", Amount: 0, AccountID: 1},      // zero amount
		{Name: "Sub", Amount: -10, AccountID: 1},    // negative amount
		{Name: "Sub", Amount: 100, AccountID: 0},    // zero account
	}
	for _, r := range cases {
		_, err := svc.Create(context.Background(), 1, r)
		if err != domain.ErrInvalidInput {
			t.Errorf("Create(%+v): expected ErrInvalidInput, got %v", r, err)
		}
	}
}

func TestRecurringCreate_DefaultFrequency(t *testing.T) {
	var gotFreq string
	store := &mockStore{
		createRecurringFn: func(_ context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error) {
			gotFreq = r.Frequency
			return r, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	_, err := svc.Create(context.Background(), 1, domain.RecurringPayment{
		Name: "Netflix", Amount: 500, AccountID: 1, Frequency: "biweekly",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFreq != "monthly" {
		t.Errorf("frequency = %q, want %q", gotFreq, "monthly")
	}
}

func TestRecurringCreate_Success(t *testing.T) {
	want := &domain.RecurringPayment{ID: 7, Name: "Rent", Amount: 30000, Frequency: "monthly"}
	store := &mockStore{
		createRecurringFn: func(_ context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error) {
			return want, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	got, err := svc.Create(context.Background(), 1, domain.RecurringPayment{
		Name: "Rent", Amount: 30000, AccountID: 1, Frequency: "monthly",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 7 {
		t.Errorf("got.ID = %d, want 7", got.ID)
	}
}

func TestRecurringUpdate_InvalidInput(t *testing.T) {
	svc := NewRecurring(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, domain.RecurringPayment{
		Name: "", Amount: 100, AccountID: 1,
	})
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestRecurringUpdate_Success(t *testing.T) {
	want := &domain.RecurringPayment{ID: 3, Name: "Gym", Amount: 1500}
	store := &mockStore{
		updateRecurringFn: func(_ context.Context, id, userID int64, name string, amount float64, freq string, nextDate time.Time, accountID int64, categoryID *int64) (*domain.RecurringPayment, error) {
			return want, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	got, err := svc.Update(context.Background(), 3, 1, domain.RecurringPayment{
		Name: "Gym", Amount: 1500, AccountID: 2, Frequency: "monthly",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("got.ID = %d, want 3", got.ID)
	}
}

func TestRecurringDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteRecurringFn: func(_ context.Context, id, userID int64) error {
			called = true
			return nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	if err := svc.Delete(context.Background(), 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("DeleteRecurring was not called")
	}
}
