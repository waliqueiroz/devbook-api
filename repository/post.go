package repository

import (
	"database/sql"

	"github.com/waliqueiroz/devbook-api/model"
)

type postRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{db}
}

// Create inserts a post into database
func (repository postRepository) Create(post model.Post) (uint64, error) {
	statement, err := repository.db.Prepare("insert into posts (title, content, author_id) values (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(post.Title, post.Content, post.AuthorID)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertID), nil
}
