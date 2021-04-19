package mock

import (
	"encoding/json"
	"io/ioutil"

	"github.com/waliqueiroz/devbook-api/model"
)

type UserRepositoryMock struct{}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepositoryMock {
	return &UserRepositoryMock{}
}

// Create inserts a user into database
func (repository UserRepositoryMock) Create(user model.User) (model.User, error) {
	return repository.getStoredUser()
}

// FindByNameOrNick returns all users that name or nick match with the argument
func (repository UserRepositoryMock) FindByNameOrNick(nameOrNick string) ([]model.User, error) {
	return repository.getStoredUserList()
}

// FindByID returns a user thar match with a given ID
func (repository UserRepositoryMock) FindByID(userID uint64) (model.User, error) {
	return repository.getStoredUser()
}

// Update updates a user in database
func (repository UserRepositoryMock) Update(userID uint64, user model.User) error {
	return nil
}

// Delete deletes a user in database
func (repository UserRepositoryMock) Delete(userID uint64) error {
	return nil
}

// FindByEmail returns all users that email match with the argument
func (repository UserRepositoryMock) FindByEmail(email string) (model.User, error) {
	storedUserJson, _ := ioutil.ReadFile("../test/resource/json/stored_user.json")

	var storedUser model.User

	json.Unmarshal(storedUserJson, &storedUser)

	return storedUser, nil
}

// Follow allows a user to follow another
func (repository UserRepositoryMock) Follow(userID, followerID uint64) error {
	return nil
}

// Unfollow allows a user to unfollow another
func (repository UserRepositoryMock) Unfollow(userID, followerID uint64) error {
	return nil
}

// SearchFollowers returns a list of followers for a given user
func (repository UserRepositoryMock) SearchFollowers(userID uint64) ([]model.User, error) {
	return repository.getStoredUserList()
}

// SearchFollowing returns a list of users that a given user is following
func (repository UserRepositoryMock) SearchFollowing(userID uint64) ([]model.User, error) {
	return repository.getStoredUserList()
}

// FindPassword returns the hashed password of a given user
func (repository UserRepositoryMock) FindPassword(userID uint64) (string, error) {
	return "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6", nil
}

// UpdatePassword updates the password for a given user
func (repository UserRepositoryMock) UpdatePassword(userID uint64, password string) error {
	return nil
}

func (repository UserRepositoryMock) getStoredUser() (model.User, error) {
	storedUserJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var storedUser model.User

	json.Unmarshal(storedUserJson, &storedUser)

	return storedUser, nil
}

func (repository UserRepositoryMock) getStoredUserList() ([]model.User, error) {
	storedUserlistJson, _ := ioutil.ReadFile("../test/resource/json/stored_user_list.json")

	var storedUserList []model.User

	json.Unmarshal(storedUserlistJson, &storedUserList)

	return storedUserList, nil
}
