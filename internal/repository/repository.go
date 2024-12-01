package repository

import "context"

// Интерфейс репозитория с полным CRUD для категорий и транзакций
type Repository interface {
	// Категории
	CreateCategory(ctx context.Context, userID int, name string) (int, error)
	GetCategories(ctx context.Context, userID int) ([]Category, error)
	UpdateCategory(ctx context.Context, userID, categoryID int, name string) error
	DeleteCategory(ctx context.Context, userID, categoryID int) error

	// Транзакции
	CreateTransaction(ctx context.Context, txn Transaction) (int, error)
	GetTransactions(ctx context.Context, userID int) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, txn Transaction) error
	DeleteTransaction(ctx context.Context, userID, txnID int) error

	// Аналитика
	GetIncomeAndExpenses(ctx context.Context, userID int) (*Analytics, error)
	GetCategoryAnalytics(ctx context.Context, userID int) ([]CategoryAnalytics, error)
	GetIncomeAndExpensesFiltered(ctx context.Context, userID int, startDate, endDate string) (*Analytics, error)
	GetCategoryAnalyticsFiltered(ctx context.Context, userID int, startDate, endDate string) ([]CategoryAnalytics, error)

	// Методы для пользователей
	CreateUser(ctx context.Context, username, hashedPassword string) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}
