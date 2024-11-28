package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt string
}

func (db *DB) CreateUser(ctx context.Context, username, passwordHash string) (int, error) {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
	var id int
	err := db.Conn.QueryRowContext(ctx, query, username, passwordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (db *DB) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := "SELECT id, username, password_hash, created_at FROM users WHERE username = $1"
	var user User
	err := db.Conn.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
