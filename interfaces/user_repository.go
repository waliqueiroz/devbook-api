package interfaces

import "github.com/waliqueiroz/devbook-api/model"

// UserRepository describes a user repository interface
type UserRepository interface {
	Create(model.User) (model.User, error)
	FindByNameOrNick(string) ([]model.User, error)
	FindByID(uint64) (model.User, error)
	Update(uint64, model.User) error
	Delete(uint64) error
	FindByEmail(string) (model.User, error)
	Follow(uint64, uint64) error
	Unfollow(uint64, uint64) error
	SearchFollowers(uint64) ([]model.User, error)
	SearchFollowing(uint64) ([]model.User, error)
	FindPassword(uint64) (string, error)
	UpdatePassword(uint64, string) error
}
