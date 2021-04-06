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

// Index returns all posts by a user and from who they are following
func (repository postRepository) Index(userID uint64) ([]model.Post, error) {

	rows, err := repository.db.Query(`select distinct
										p.*,
										u.nick
									from
										posts p
									join users u on
										p.author_id = u.id
									join followers f on 
										p.author_id = f.user_id 
									where
										u.id = ? or f.follower_id = ? order by p.id desc`, userID, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.Likes, &post.CreatedAt, &post.AuthorNick)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// Update updates a post in database
func (repository postRepository) Update(postID uint64, post model.Post) error {

	statement, err := repository.db.Prepare("update posts set title = ?, content = ? where id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(post.Title, post.Content, postID)
	if err != nil {
		return err
	}

	return nil
}

// Update deletes a post from database
func (repository postRepository) Delete(postID uint64) error {
	statement, err := repository.db.Prepare("delete from posts where id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(postID)
	if err != nil {
		return err
	}

	return nil
}

// FindByUser returns all posts from a given user
func (repository postRepository) FindByUser(userID uint64) ([]model.Post, error) {

	rows, err := repository.db.Query(`select distinct
										p.*,
										u.nick
									from
										posts p
									join users u on
										p.author_id = u.id
									where
										u.id = ?`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.Likes, &post.CreatedAt, &post.AuthorNick)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
