package repository

import (
	"database/sql"
	"fmt"

	"github.com/waliqueiroz/devbook-api/model"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

// Create inserts an user into database
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

// FindByNameOrNick returns all users that name or nick match with the argument
func (repository userRepository) FindByNameOrNick(nameOrNick string) ([]model.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick)

	rows, err := repository.db.Query("select id, name, nick, email, created_at from users where name like ? or nick like ?", nameOrNick, nameOrNick)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		err = rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// FindByID returns all users that name or nick match with the argument
func (repository userRepository) FindByID(userID uint64) (model.User, error) {

	rows, err := repository.db.Query("select id, name, nick, email, created_at from users where id = ?", userID)

	if err != nil {
		return model.User{}, err
	}

	defer rows.Close()

	var user model.User

	if rows.Next() {

		err = rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt)

		if err != nil {
			return model.User{}, err
		}

	}

	return user, nil
}

// Update updates an user in database
func (repository userRepository) Update(userID uint64, user model.User) error {

	statement, err := repository.db.Prepare("update users set name = ?, nick = ?, email = ? where id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Nick, user.Email, userID)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an user in database
func (repository userRepository) Delete(userID uint64) error {

	statement, err := repository.db.Prepare("delete from users where id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}
