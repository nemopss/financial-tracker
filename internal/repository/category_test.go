package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}
	mock.ExpectQuery("INSERT INTO categories").
		WithArgs("Groceries", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := r.CreateCategory(context.Background(), 1, "Groceries")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}
	mock.ExpectQuery("SELECT id, name FROM categories WHERE user_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Groceries").
			AddRow(2, "Entertainment"))

	categories, err := r.GetCategories(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Groceries", categories[0].Name)
	assert.Equal(t, "Entertainment", categories[1].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}
	mock.ExpectExec("UPDATE categories SET name = \\$1 WHERE id = \\$2 AND user_id = \\$3").
		WithArgs("Updated Category", 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = r.UpdateCategory(context.Background(), 1, 1, "Updated Category")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCategoryNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}
	mock.ExpectExec("UPDATE categories SET name = \\$1 WHERE id = \\$2 AND user_id = \\$3").
		WithArgs("Updated Category", 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = r.UpdateCategory(context.Background(), 1, 1, "Updated Category")
	assert.Error(t, err)
	assert.Equal(t, "no category found or not authorized", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}
	mock.ExpectExec("DELETE FROM categories WHERE id = \\$1 AND user_id = \\$2").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = r.DeleteCategory(context.Background(), 1, 1)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCategoryNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := &DB{Conn: db}
	mock.ExpectExec("DELETE FROM categories WHERE id = \\$1 AND user_id = \\$2").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = r.DeleteCategory(context.Background(), 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "no category found or not authorized", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
