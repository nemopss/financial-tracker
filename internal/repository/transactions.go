package repository

import (
	"context"
	"fmt"
	"time"
)

type Transaction struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"time"`
	Description string    `json:"description"`
	CategoryID  int       `json:"category_id"`
	UserID      int       `json:"user_id"`
}

func (db *DB) CreateTransaction(ctx context.Context, txn Transaction) (int, error) {
	query := "INSERT INTO transactions (amount, date, description, category_id, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id int
	err := db.Conn.QueryRowContext(ctx, query, txn.Amount, txn.Date, txn.Description, txn.CategoryID, txn.UserID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}
	return id, nil
}

func (db *DB) GetTransactions(ctx context.Context, userID int) ([]Transaction, error) {
	query := "SELECT id, amount, date, description, category_id, user_id FROM transactions WHERE user_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(&txn.ID, &txn.Amount, &txn.Date, &txn.Description, &txn.CategoryID, &txn.UserID); err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}

func (db *DB) UpdateTransaction(ctx context.Context, txn Transaction) error {
	query := "UPDATE transactions SET amount = $1, date = $2, description = $3, category_id = $4 WHERE id = $5 AND user_id = $6"
	res, err := db.Conn.ExecContext(ctx, query, txn.Amount, txn.Date, txn.Description, txn.CategoryID, txn.ID, txn.UserID)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no transaction found or not authorized")
	}

	return nil
}

func (db *DB) DeleteTransaction(ctx context.Context, userID, txnID int) error {
	query := "DELETE FROM transactions WHERE id = $1 AND user_id = $2"
	res, err := db.Conn.ExecContext(ctx, query, txnID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no transaction found or not authorized")
	}

	return nil
}
