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

func (db *DB) GetExecutorCategoryName(data string) (string, error) {
	var name string
	err := db.db.Get(&name, "SELECT name FROM executor_categories WHERE data = $1", data)
	return name, err
}

func (db *DB) GetExecutorSpheres(category string) ([]models.ExecutorSphere, error) {
	var spheres []models.ExecutorSphere
	err := db.db.Select(&spheres, "SELECT es.id, es.name, es.uniq, hash FROM executor_spheres es INNER JOIN executor_categories ec ON es.parent_category = ec.id WHERE ec.data = $1", category)
	return spheres, err
}

func (db *DB) GetExecutorSphereName(uniq string) (string, error) {
	var name string
	err := db.db.Get(&name, "SELECT name FROM executor_spheres WHERE uniq = $1", uniq)
	return name, err
}

func (db *DB) GetOfferCategories() ([]models.OfferCategory, error) {
	var categories []models.OfferCategory
	err := db.db.Select(&categories, "SELECT id, name, data, hash FROM offer_categories")
	return categories, err
}

func (db *DB) GetOfferCategoryName(data string) (string, error) {
	var name string
	err := db.db.Get(&name, "SELECT name FROM offer_categories WHERE data = $1", data)
	return name, err
}

func (db *DB) GetOfferCategoryHash(data string) (string, error) {
	var hash string
	err := db.db.Get(&hash, "SELECT hash FROM offer_categories WHERE data = $1", data)
	return hash, err
}

func (db *DB) AddPost(post models.Post) error {
	_, err := db.db.Exec("INSERT INTO posts (title, post_type, post_link, post_id, user_id, created_at) VALUES ($1, $2, $3, $4, $5, now())", post.Title, post.Type, post.Link, post.MessageId, post.UserId)
	return err
}

func (db *DB) DeletePost(id int) error {
	_, err := db.db.Exec("DELETE FROM posts WHERE id = $1", id)
	return err
}

func (db *DB) GetUserPosts(id int64) ([]models.Post, error) {
	var posts []models.Post
	err := db.db.Select(&posts, "SELECT * FROM posts WHERE user_id = $1", id)
	return posts, err
}
