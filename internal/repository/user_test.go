package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "hashedpassword").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := r.CreateUser(context.Background(), "testuser", "hashedpassword")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "hashedpassword").
		WillReturnError(fmt.Errorf("db error"))

	id, err := r.CreateUser(context.Background(), "testuser", "hashedpassword")
	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "failed to create user: db error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT id, username, password_hash, created_at FROM users WHERE username = \\$1").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at"}).
			AddRow(1, "testuser", "hashedpassword", "2024-12-01"))

	user, err := r.GetUserByUsername(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "2024-12-01", user.CreatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsernameNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT id, username, password_hash, created_at FROM users WHERE username = \\$1").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	user, err := r.GetUserByUsername(context.Background(), "nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsernameError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}

	mock.ExpectQuery("SELECT id, username, password_hash, created_at FROM users WHERE username = \\$1").
		WithArgs("testuser").
		WillReturnError(fmt.Errorf("db error"))

	user, err := r.GetUserByUsername(context.Background(), "testuser")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "failed to get user: db error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
