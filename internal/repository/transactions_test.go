package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mockTransaction := Transaction{
		Amount:      100.50,
		Date:        time.Now(),
		Description: "Groceries",
		CategoryID:  1,
		UserID:      1,
	}

	mock.ExpectQuery("INSERT INTO transactions").
		WithArgs(mockTransaction.Amount, mockTransaction.Date, mockTransaction.Description, mockTransaction.CategoryID, mockTransaction.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := r.CreateTransaction(context.Background(), mockTransaction)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT id, amount, date, description, category_id, user_id FROM transactions WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "amount", "date", "description", "category_id", "user_id"}).
			AddRow(1, 100.50, time.Now(), "Groceries", 1, 1).
			AddRow(2, -50.00, time.Now(), "Entertainment", 2, 1))

	transactions, err := r.GetTransactions(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, 100.50, transactions[0].Amount)
	assert.Equal(t, "Groceries", transactions[0].Description)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mockTransaction := Transaction{
		ID:          1,
		Amount:      150.75,
		Date:        time.Now(),
		Description: "Updated Groceries",
		CategoryID:  1,
		UserID:      1,
	}

	mock.ExpectExec("UPDATE transactions SET amount = \\$1, date = \\$2, description = \\$3, category_id = \\$4 WHERE id = \\$5 AND user_id = \\$6").
		WithArgs(mockTransaction.Amount, mockTransaction.Date, mockTransaction.Description, mockTransaction.CategoryID, mockTransaction.ID, mockTransaction.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = r.UpdateTransaction(context.Background(), mockTransaction)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectExec("DELETE FROM transactions WHERE id = \\$1 AND user_id = \\$2").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = r.DeleteTransaction(context.Background(), 1, 1)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
