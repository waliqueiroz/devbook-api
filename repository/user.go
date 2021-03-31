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

// FindByEmail returns all users that email match with the argument
func (repository userRepository) FindByEmail(email string) (model.User, error) {

	rows, err := repository.db.Query("select id, password from users where email = ?", email)

	if err != nil {
		return model.User{}, err
	}

	defer rows.Close()

	var user model.User

	if rows.Next() {

		err = rows.Scan(&user.ID, &user.Password)

		if err != nil {
			return model.User{}, err
		}

	}

	return user, nil
}

// Follow allows an user to follow another
func (repository userRepository) Follow(userID, followerID uint64) error {
	statement, err := repository.db.Prepare("insert ignore into followers (user_id, follower_id) values (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userID, followerID)
	if err != nil {
		return err
	}

	return nil
}

// Unfollow allows an user to unfollow another
func (repository userRepository) Unfollow(userID, followerID uint64) error {
	statement, err := repository.db.Prepare("delete from followers where user_id = ? and follower_id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userID, followerID)
	if err != nil {
		return err
	}

	return nil
}

// SearchFollowers returns a list of followers for a given user
func (repository userRepository) SearchFollowers(userID uint64) ([]model.User, error) {
	rows, err := repository.db.Query(`select u.id, u.name, u.nick, u.email, u.created_at 
									from users u join followers f on u.id = f.follower_id where f.user_id = ?`,
		userID)
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

// SearchFollowing returns a list of users that a given user is following
func (repository userRepository) SearchFollowing(userID uint64) ([]model.User, error) {
	rows, err := repository.db.Query(`select u.id, u.name, u.nick, u.email, u.created_at 
									from users u join followers f on u.id = f.user_id where f.follower_id = ?`,
		userID)
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

// FindPassword returns the hashed password of a given user
func (repository userRepository) FindPassword(userID uint64) (string, error) {
	rows, err := repository.db.Query(`select password from users where id = ?`,
		userID)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	var user model.User
	if rows.Next() {

		err = rows.Scan(&user.Password)

		if err != nil {
			return "", err
		}

	}

	return user.Password, nil

}

// UpdatePassword updates the password for a given user
func (repository userRepository) UpdatePassword(userID uint64, password string) error {
	statement, err := repository.db.Prepare("update users set password = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(password, userID)
	if err != nil {
		return err
	}

	return nil

}
