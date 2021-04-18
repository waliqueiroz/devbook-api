package interfaces

import "github.com/waliqueiroz/devbook-api/model"

// PostRepository describes a post repository interface
type PostRepository interface {
	Create(model.Post) (model.Post, error)
	FindByID(uint64) (model.Post, error)
	Index(uint64) ([]model.Post, error)
	Update(uint64, model.Post) error
	Delete(uint64) error
	FindByUser(uint64) ([]model.Post, error)
	LikePost(uint64) error
	DeslikePost(uint64) error
}
