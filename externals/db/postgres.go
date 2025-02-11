package db

import (
	"boutonsjob/internals/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) AddExecutorCategory(name string, data string) error {
	_, err := db.db.Exec("INSERT INTO executor_categories (name, data) VALUES ($1, $2)", name, data)
	return err
}

func (db *DB) GetExecutorCategories() ([]models.ExecutorCategory, error) {
	var categories []models.ExecutorCategory
	err := db.db.Select(&categories, "SELECT id, name, data FROM executor_categories")
	return categories, err
}
