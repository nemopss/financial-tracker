package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetIncomeAndExpenses(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(CASE WHEN amount > 0 THEN amount ELSE 0 END\\), 0\\) AS total_income,").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"total_income", "total_expense"}).AddRow(5000.00, -2000.00))

	analytics, err := r.GetIncomeAndExpenses(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 5000.00, analytics.TotalIncome)
	assert.Equal(t, -2000.00, analytics.TotalExpense)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetIncomeAndExpensesError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(CASE WHEN amount > 0 THEN amount ELSE 0 END\\), 0\\) AS total_income,").
		WithArgs(1).
		WillReturnError(fmt.Errorf("db error"))

	analytics, err := r.GetIncomeAndExpenses(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, analytics)
	assert.Equal(t, "failed to fetch analytics: db error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetIncomeAndExpensesFiltered(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(CASE WHEN amount > 0 THEN amount ELSE 0 END\\), 0\\) AS total_income,").
		WithArgs(1, "2024-01-01", "2024-12-31").
		WillReturnRows(sqlmock.NewRows([]string{"total_income", "total_expense"}).AddRow(3000.00, -1000.00))

	analytics, err := r.GetIncomeAndExpensesFiltered(context.Background(), 1, "2024-01-01", "2024-12-31")
	assert.NoError(t, err)
	assert.Equal(t, 3000.00, analytics.TotalIncome)
	assert.Equal(t, -1000.00, analytics.TotalExpense)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategoryAnalytics(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT c.name AS category_name, COALESCE\\(SUM\\(t.amount\\), 0\\) AS total_amount").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"category_name", "total_amount"}).
			AddRow("Groceries", -150.00).
			AddRow("Salary", 5000.00))

	analytics, err := r.GetCategoryAnalytics(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, analytics, 2)
	assert.Equal(t, "Groceries", analytics[0].CategoryName)
	assert.Equal(t, -150.00, analytics[0].TotalAmount)
	assert.Equal(t, "Salary", analytics[1].CategoryName)
	assert.Equal(t, 5000.00, analytics[1].TotalAmount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategoryAnalyticsFiltered(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT c.name AS category_name, COALESCE\\(SUM\\(t.amount\\), 0\\) AS total_amount").
		WithArgs(1, "2024-01-01", "2024-12-31").
		WillReturnRows(sqlmock.NewRows([]string{"category_name", "total_amount"}).
			AddRow("Groceries", -200.00).
			AddRow("Entertainment", -100.00))

	analytics, err := r.GetCategoryAnalyticsFiltered(context.Background(), 1, "2024-01-01", "2024-12-31")
	assert.NoError(t, err)
	assert.Len(t, analytics, 2)
	assert.Equal(t, "Groceries", analytics[0].CategoryName)
	assert.Equal(t, -200.00, analytics[0].TotalAmount)
	assert.Equal(t, "Entertainment", analytics[1].CategoryName)
	assert.Equal(t, -100.00, analytics[1].TotalAmount)

	assert.NoError(t, mock.ExpectationsWereMet())
}
