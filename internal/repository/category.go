package repository

import (
	"context"
	"fmt"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (db *DB) CreateCategory(ctx context.Context, userID int, name string) (int, error) {
	query := "INSERT INTO categories (name, user_id) VALUES ($1, $2) RETURNING id"
	var id int
	err := db.Conn.QueryRowContext(ctx, query, name, userID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}
	return id, nil
}

func (db *DB) GetCategories(ctx context.Context, userID int) ([]Category, error) {
	query := "SELECT id, name FROM categories WHERE user_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (db *DB) UpdateCategory(ctx context.Context, userID, categoryID int, name string) error {
	query := "UPDATE categories SET name = $1 WHERE id = $2 AND user_id = $3"
	res, err := db.Conn.ExecContext(ctx, query, name, categoryID, userID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no category found or not authorized")
	}

	return nil
}

func (db *DB) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	query := "DELETE FROM categories WHERE id = $1 AND user_id = $2"
	res, err := db.Conn.ExecContext(ctx, query, categoryID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no category found or not authorized")
	}

	return nil
}
