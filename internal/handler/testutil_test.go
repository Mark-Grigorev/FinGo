package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
	"github.com/Mark-Grigorev/FinGo/internal/service"
	"github.com/Mark-Grigorev/FinGo/pkg/token"
)

// mockStore implements repository.Storer without a real database.
type mockStore struct {
	getUserByEmailFn    func(ctx context.Context, email string) (*domain.User, error)
	getUserByIDFn       func(ctx context.Context, id int64) (*domain.User, error)
	createUserFn        func(ctx context.Context, email, name, hash string) (*domain.User, error)
	listAccountsFn      func(ctx context.Context, userID int64) ([]domain.Account, error)
	getAccountFn        func(ctx context.Context, id, userID int64) (*domain.Account, error)
	createAccountFn     func(ctx context.Context, a *domain.Account) (*domain.Account, error)
	updateAccountFn     func(ctx context.Context, id, userID int64, name, typ, currency string) (*domain.Account, error)
	deleteAccountFn     func(ctx context.Context, id, userID int64) error
	listTransactionsFn    func(ctx context.Context, userID int64, f repository.TransactionFilter) ([]domain.Transaction, int, error)
	createTransactionFn   func(ctx context.Context, t *domain.Transaction) (*domain.Transaction, error)
	deleteTransactionFn   func(ctx context.Context, id, userID int64) error
	getDashboardSummaryFn func(ctx context.Context, userID int64) (*repository.DashboardSummary, error)
	listCategoriesFn      func(ctx context.Context, userID int64) ([]domain.Category, error)
	createCategoryFn      func(ctx context.Context, c *domain.Category) (*domain.Category, error)
	updateCategoryFn      func(ctx context.Context, id, userID int64, name, icon, color string) (*domain.Category, error)
	deleteCategoryFn      func(ctx context.Context, id, userID int64) error
	getReportFn           func(ctx context.Context, userID int64, from, to time.Time) (*repository.ReportResult, error)
	listBudgetsFn    func(ctx context.Context, userID int64, month time.Time) ([]domain.Budget, error)
	createBudgetFn   func(ctx context.Context, b *domain.Budget) (*domain.Budget, error)
	updateBudgetFn   func(ctx context.Context, id, userID int64, limit float64) (*domain.Budget, error)
	deleteBudgetFn   func(ctx context.Context, id, userID int64) error
	listRecurringFn   func(ctx context.Context, userID int64) ([]domain.RecurringPayment, error)
	createRecurringFn func(ctx context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error)
	updateRecurringFn func(ctx context.Context, id, userID int64, name string, amount float64, frequency string, nextDate time.Time, accountID int64, categoryID *int64) (*domain.RecurringPayment, error)
	deleteRecurringFn func(ctx context.Context, id, userID int64) error
	updateUserFn      func(ctx context.Context, id int64, name, email string) (*domain.User, error)
	updatePasswordFn  func(ctx context.Context, id int64, hash string) error
}

func (m *mockStore) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getUserByEmailFn == nil {
		panic("mockStore: GetUserByEmail not configured")
	}
	return m.getUserByEmailFn(ctx, email)
}
func (m *mockStore) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if m.getUserByIDFn == nil {
		panic("mockStore: GetUserByID not configured")
	}
	return m.getUserByIDFn(ctx, id)
}
func (m *mockStore) CreateUser(ctx context.Context, email, name, hash string) (*domain.User, error) {
	if m.createUserFn == nil {
		panic("mockStore: CreateUser not configured")
	}
	return m.createUserFn(ctx, email, name, hash)
}
func (m *mockStore) ListAccounts(ctx context.Context, userID int64) ([]domain.Account, error) {
	if m.listAccountsFn == nil {
		panic("mockStore: ListAccounts not configured")
	}
	return m.listAccountsFn(ctx, userID)
}
func (m *mockStore) GetAccount(ctx context.Context, id, userID int64) (*domain.Account, error) {
	if m.getAccountFn == nil {
		panic("mockStore: GetAccount not configured")
	}
	return m.getAccountFn(ctx, id, userID)
}
func (m *mockStore) CreateAccount(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	if m.createAccountFn == nil {
		panic("mockStore: CreateAccount not configured")
	}
	return m.createAccountFn(ctx, a)
}
func (m *mockStore) UpdateAccount(ctx context.Context, id, userID int64, name, typ, currency string) (*domain.Account, error) {
	if m.updateAccountFn == nil {
		panic("mockStore: UpdateAccount not configured")
	}
	return m.updateAccountFn(ctx, id, userID, name, typ, currency)
}
func (m *mockStore) DeleteAccount(ctx context.Context, id, userID int64) error {
	if m.deleteAccountFn == nil {
		panic("mockStore: DeleteAccount not configured")
	}
	return m.deleteAccountFn(ctx, id, userID)
}
func (m *mockStore) ListTransactions(ctx context.Context, userID int64, f repository.TransactionFilter) ([]domain.Transaction, int, error) {
	if m.listTransactionsFn == nil {
		panic("mockStore: ListTransactions not configured")
	}
	return m.listTransactionsFn(ctx, userID, f)
}
func (m *mockStore) CreateTransaction(ctx context.Context, t *domain.Transaction) (*domain.Transaction, error) {
	if m.createTransactionFn == nil {
		panic("mockStore: CreateTransaction not configured")
	}
	return m.createTransactionFn(ctx, t)
}
func (m *mockStore) DeleteTransaction(ctx context.Context, id, userID int64) error {
	if m.deleteTransactionFn == nil { panic("mockStore: DeleteTransaction not configured") }
	return m.deleteTransactionFn(ctx, id, userID)
}

func (m *mockStore) GetDashboardSummary(ctx context.Context, userID int64) (*repository.DashboardSummary, error) {
	if m.getDashboardSummaryFn == nil {
		panic("mockStore: GetDashboardSummary not configured")
	}
	return m.getDashboardSummaryFn(ctx, userID)
}

func (m *mockStore) ListCategories(ctx context.Context, userID int64) ([]domain.Category, error) {
	if m.listCategoriesFn == nil {
		panic("mockStore: ListCategories not configured")
	}
	return m.listCategoriesFn(ctx, userID)
}

func (m *mockStore) CreateCategory(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	if m.createCategoryFn == nil {
		panic("mockStore: CreateCategory not configured")
	}
	return m.createCategoryFn(ctx, c)
}

func (m *mockStore) UpdateCategory(ctx context.Context, id, userID int64, name, icon, color string) (*domain.Category, error) {
	if m.updateCategoryFn == nil {
		panic("mockStore: UpdateCategory not configured")
	}
	return m.updateCategoryFn(ctx, id, userID, name, icon, color)
}

func (m *mockStore) DeleteCategory(ctx context.Context, id, userID int64) error {
	if m.deleteCategoryFn == nil {
		panic("mockStore: DeleteCategory not configured")
	}
	return m.deleteCategoryFn(ctx, id, userID)
}

func (m *mockStore) GetReport(ctx context.Context, userID int64, from, to time.Time) (*repository.ReportResult, error) {
	if m.getReportFn == nil {
		panic("mockStore: GetReport not configured")
	}
	return m.getReportFn(ctx, userID, from, to)
}

func (m *mockStore) ListBudgets(ctx context.Context, userID int64, month time.Time) ([]domain.Budget, error) {
	if m.listBudgetsFn == nil { panic("mockStore: ListBudgets not configured") }
	return m.listBudgetsFn(ctx, userID, month)
}
func (m *mockStore) CreateBudget(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	if m.createBudgetFn == nil { panic("mockStore: CreateBudget not configured") }
	return m.createBudgetFn(ctx, b)
}
func (m *mockStore) UpdateBudget(ctx context.Context, id, userID int64, limit float64) (*domain.Budget, error) {
	if m.updateBudgetFn == nil { panic("mockStore: UpdateBudget not configured") }
	return m.updateBudgetFn(ctx, id, userID, limit)
}
func (m *mockStore) DeleteBudget(ctx context.Context, id, userID int64) error {
	if m.deleteBudgetFn == nil { panic("mockStore: DeleteBudget not configured") }
	return m.deleteBudgetFn(ctx, id, userID)
}
func (m *mockStore) ListRecurring(ctx context.Context, userID int64) ([]domain.RecurringPayment, error) {
	if m.listRecurringFn == nil { panic("mockStore: ListRecurring not configured") }
	return m.listRecurringFn(ctx, userID)
}
func (m *mockStore) CreateRecurring(ctx context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error) {
	if m.createRecurringFn == nil { panic("mockStore: CreateRecurring not configured") }
	return m.createRecurringFn(ctx, r)
}
func (m *mockStore) UpdateRecurring(ctx context.Context, id, userID int64, name string, amount float64, frequency string, nextDate time.Time, accountID int64, categoryID *int64) (*domain.RecurringPayment, error) {
	if m.updateRecurringFn == nil { panic("mockStore: UpdateRecurring not configured") }
	return m.updateRecurringFn(ctx, id, userID, name, amount, frequency, nextDate, accountID, categoryID)
}
func (m *mockStore) DeleteRecurring(ctx context.Context, id, userID int64) error {
	if m.deleteRecurringFn == nil { panic("mockStore: DeleteRecurring not configured") }
	return m.deleteRecurringFn(ctx, id, userID)
}
func (m *mockStore) UpdateUser(ctx context.Context, id int64, name, email string) (*domain.User, error) {
	if m.updateUserFn == nil { panic("mockStore: UpdateUser not configured") }
	return m.updateUserFn(ctx, id, name, email)
}
func (m *mockStore) UpdatePassword(ctx context.Context, id int64, hash string) error {
	if m.updatePasswordFn == nil { panic("mockStore: UpdatePassword not configured") }
	return m.updatePasswordFn(ctx, id, hash)
}

// testEnv wires all handlers with a mock store and a real token maker.
type testEnv struct {
	router *gin.Engine
	maker  *token.Maker
	store  *mockStore
}

func newTestEnv() *testEnv {
	gin.SetMode(gin.TestMode)
	store := &mockStore{}
	maker, _ := token.New("", time.Hour)
	log := slog.Default()

	authSvc := service.NewAuth(store, maker, log)
	accSvc := service.NewAccount(store, log)
	txSvc := service.NewTransaction(store, log)
	dashboardSvc := service.NewDashboard(store, log)
	categorySvc := service.NewCategory(store, log)
	budgetSvc := service.NewBudget(store, log)
	recurringSvc := service.NewRecurring(store, log)

	r := gin.New()
	authH := &authHandler{svc: authSvc, log: log}
	accH := &accountHandler{svc: accSvc, log: log}
	txH := &transactionHandler{svc: txSvc, log: log}
	dashboardH := &dashboardHandler{svc: dashboardSvc, log: log}
	catH := &categoryHandler{svc: categorySvc, log: log}
	budgetH := &budgetHandler{svc: budgetSvc, log: log}
	recurringH := &recurringHandler{svc: recurringSvc, log: log}

	api := r.Group("/api")
	api.POST("/auth/login", authH.login)
	api.POST("/auth/register", authH.register)
	api.POST("/auth/logout", authH.logout)

	protected := api.Group("/")
	protected.Use(authMiddleware(authSvc))
	protected.GET("/auth/me", authH.me)
	protected.GET("/accounts", accH.list)
	protected.POST("/accounts", accH.create)
	protected.PUT("/accounts/:id", accH.update)
	protected.DELETE("/accounts/:id", accH.delete)
	protected.GET("/transactions", txH.list)
	protected.POST("/transactions", txH.create)
	protected.DELETE("/transactions/:id", txH.delete)
	protected.GET("/dashboard/summary", dashboardH.summary)
	protected.GET("/dashboard/report", dashboardH.report)
	protected.GET("/budgets", budgetH.list)
	protected.POST("/budgets", budgetH.create)
	protected.PUT("/budgets/:id", budgetH.update)
	protected.DELETE("/budgets/:id", budgetH.delete)
	protected.GET("/recurring", recurringH.list)
	protected.POST("/recurring", recurringH.create)
	protected.PUT("/recurring/:id", recurringH.update)
	protected.DELETE("/recurring/:id", recurringH.delete)
	protected.GET("/categories", catH.list)
	protected.POST("/categories", catH.create)
	protected.PUT("/categories/:id", catH.update)
	protected.DELETE("/categories/:id", catH.delete)

	return &testEnv{router: r, maker: maker, store: store}
}

// bearerToken returns a valid "Bearer <token>" header value for the given user.
func (e *testEnv) bearerToken(userID int64) string {
	tok, _, _ := e.maker.CreateToken(userID, "test@example.com")
	return "Bearer " + tok
}

// makeHash creates a bcrypt hash with MinCost for use in test fixtures.
func makeHash(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h)
}
