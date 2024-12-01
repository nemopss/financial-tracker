package repository

import (
	"context"
	"log"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

// Методы для категорий
func (m *MockRepo) CreateCategory(ctx context.Context, userID int, name string) (int, error) {
	args := m.Called(ctx, userID, name)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) GetCategories(ctx context.Context, userID int) ([]Category, error) {
	args := m.Called(ctx, userID)
	categories, ok := args.Get(0).([]Category)
	if !ok {
		log.Printf("MockRepo.GetCategories: failed to cast to []Category. Got: %+v", args.Get(0))
	}
	return categories, args.Error(1)
}

func (m *MockRepo) UpdateCategory(ctx context.Context, userID, categoryID int, name string) error {
	args := m.Called(ctx, userID, categoryID, name)
	return args.Error(0)
}

func (m *MockRepo) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	args := m.Called(ctx, userID, categoryID)
	return args.Error(0)
}

// Методы для транзакций
func (m *MockRepo) CreateTransaction(ctx context.Context, txn Transaction) (int, error) {
	args := m.Called(ctx, txn)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) GetTransactions(ctx context.Context, userID int) ([]Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Transaction), args.Error(1)
}

func (m *MockRepo) UpdateTransaction(ctx context.Context, txn Transaction) error {
	args := m.Called(ctx, txn)
	return args.Error(0)
}

func (m *MockRepo) DeleteTransaction(ctx context.Context, userID, txnID int) error {
	args := m.Called(ctx, userID, txnID)
	return args.Error(0)
}

// Методы для аналитики
func (m *MockRepo) GetIncomeAndExpenses(ctx context.Context, userID int) (*Analytics, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*Analytics), args.Error(1)
}

func (m *MockRepo) GetCategoryAnalytics(ctx context.Context, userID int) ([]CategoryAnalytics, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]CategoryAnalytics), args.Error(1)
}

// Методы для аналитики
func (m *MockRepo) GetIncomeAndExpensesFiltered(ctx context.Context, userID int, startDate, endDate string) (*Analytics, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	return args.Get(0).(*Analytics), args.Error(1)
}

func (m *MockRepo) GetCategoryAnalyticsFiltered(ctx context.Context, userID int, startDate, endDate string) ([]CategoryAnalytics, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	return args.Get(0).([]CategoryAnalytics), args.Error(1)
}

// Реализация метода CreateUser
func (m *MockRepo) CreateUser(ctx context.Context, username, hashedPassword string) (int, error) {
	args := m.Called(ctx, username, hashedPassword)
	return args.Int(0), args.Error(1)
}

// Реализация метода GetUserByUsername
func (m *MockRepo) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	user, ok := args.Get(0).(*User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}
