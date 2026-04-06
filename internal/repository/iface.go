package repository

import (
	"context"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

// Storer abstracts all database operations used by the service layer.
type Storer interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreatePasswordReset(ctx context.Context, token string, userID int64, expiresAt time.Time) error
	GetPasswordReset(ctx context.Context, token string) (*domain.PasswordReset, error)
	MarkPasswordResetUsed(ctx context.Context, token string) error
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	CreateUser(ctx context.Context, email, name, hash string) (*domain.User, error)
	UpdateUser(ctx context.Context, id int64, name, email string) (*domain.User, error)
	UpdatePassword(ctx context.Context, id int64, hash string) error

	ListAccounts(ctx context.Context, userID int64) ([]domain.Account, error)
	GetAccount(ctx context.Context, id, userID int64) (*domain.Account, error)
	CreateAccount(ctx context.Context, a *domain.Account) (*domain.Account, error)
	UpdateAccount(ctx context.Context, id, userID int64, name, typ, currency string) (*domain.Account, error)
	DeleteAccount(ctx context.Context, id, userID int64) error

	ListTransactions(ctx context.Context, userID int64, f TransactionFilter) ([]domain.Transaction, int, error)
	CreateTransaction(ctx context.Context, t *domain.Transaction) (*domain.Transaction, error)
	DeleteTransaction(ctx context.Context, id, userID int64) error
	ExportTransactions(ctx context.Context, userID int64, from, to time.Time) ([]domain.Transaction, error)

	GetDashboardSummary(ctx context.Context, userID int64) (*DashboardSummary, error)
	GetReport(ctx context.Context, userID int64, from, to time.Time) (*ReportResult, error)

	ListCategories(ctx context.Context, userID int64) ([]domain.Category, error)
	CreateCategory(ctx context.Context, c *domain.Category) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id, userID int64, name, icon, color string) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id, userID int64) error

	ListBudgets(ctx context.Context, userID int64, month time.Time) ([]domain.Budget, error)
	CreateBudget(ctx context.Context, b *domain.Budget) (*domain.Budget, error)
	UpdateBudget(ctx context.Context, id, userID int64, limit float64) (*domain.Budget, error)
	DeleteBudget(ctx context.Context, id, userID int64) error

	ListRecurring(ctx context.Context, userID int64) ([]domain.RecurringPayment, error)
	CreateRecurring(ctx context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error)
	UpdateRecurring(ctx context.Context, id, userID int64, name string, amount float64, frequency string, nextDate time.Time, accountID int64, categoryID *int64) (*domain.RecurringPayment, error)
	DeleteRecurring(ctx context.Context, id, userID int64) error

	GetBaseCurrency(ctx context.Context, userID int64) (string, error)
	SetBaseCurrency(ctx context.Context, userID int64, currency string) error
	ListExchangeRates(ctx context.Context, userID int64) ([]domain.ExchangeRate, error)
	UpsertExchangeRate(ctx context.Context, userID int64, currency string, rate float64) (*domain.ExchangeRate, error)
	DeleteExchangeRate(ctx context.Context, userID int64, currency string) error
	GetRatesMap(ctx context.Context, userID int64) (map[string]float64, error)
}
