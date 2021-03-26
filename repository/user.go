package repository

import (
	"database/sql"

	"github.com/waliqueiroz/devbook-api/model"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository create a new user repository
func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

// Create insert an user into database
func (repository userRepository) Create(user model.User) (uint64, error) {
	statement, err := repository.db.Prepare("insert into users (name, nick, email, password) values (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertID), nil
}
