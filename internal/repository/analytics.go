package repository

import (
	"context"
	"fmt"
)

type Analytics struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}

func (db *DB) GetIncomeAndExpenses(ctx context.Context, userID int) (*Analytics, error) {
	query := `
        SELECT
            COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS total_income,
            COALESCE(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END), 0) AS total_expense
        FROM transactions
        WHERE user_id = $1
    `
	var analytics Analytics

	err := db.Conn.QueryRowContext(ctx, query, userID).Scan(&analytics.TotalIncome, &analytics.TotalExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analytics: %w", err)
	}

	return &analytics, nil
}

type CategoryAnalytics struct {
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
}

func (db *DB) GetCategoryAnalytics(ctx context.Context, userID int) ([]CategoryAnalytics, error) {
	query := `
        SELECT c.name AS category_name, COALESCE(SUM(t.amount), 0) AS total_amount
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.user_id = $1
        GROUP BY c.name
        ORDER BY total_amount DESC
    `
	rows, err := db.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category analytics: %w", err)
	}
	defer rows.Close()

	var analytics []CategoryAnalytics
	for rows.Next() {
		var ca CategoryAnalytics
		if err := rows.Scan(&ca.CategoryName, &ca.TotalAmount); err != nil {
			return nil, err
		}
		analytics = append(analytics, ca)
	}
	return analytics, nil
}
