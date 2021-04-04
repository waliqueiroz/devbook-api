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

// FindByID returns a post that match with a given ID
func (repository postRepository) FindByID(postID uint64) (model.Post, error) {

	rows, err := repository.db.Query("select p.*, u.nick from posts p join users u on p.author_id = u.id where p.id = ?", postID)

	if err != nil {
		return model.Post{}, err
	}

	defer rows.Close()

	var post model.Post

	if rows.Next() {

		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.Likes, &post.CreatedAt, &post.AuthorNick)

		if err != nil {
			return model.Post{}, err
		}

	}

	return post, nil
}
