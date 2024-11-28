package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	Conn *sql.DB
}

func NewDB(host, port, password, user, dbname string) (*DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbname,
	)

	log.Printf("Connecting to DB: %s", dsn)
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	log.Printf("database connection estabilished")

	return &DB{Conn: db}, nil
}
