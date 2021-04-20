package mock

import (
	"encoding/json"
	"io/ioutil"

	"github.com/waliqueiroz/devbook-api/model"
)

type PostRepositoryMock struct{}

// NewPostRepositoryMock creates a new post repository
func NewPostRepository() *PostRepositoryMock {
	return &PostRepositoryMock{}
}

// Create inserts a post into database
func (repository PostRepositoryMock) Create(post model.Post) (model.Post, error) {
	return model.Post{}, nil
}

// FindByID returns a post that match with a given ID
func (repository PostRepositoryMock) FindByID(postID uint64) (model.Post, error) {
	return repository.getStoredPost()
}

// Index returns all posts by a user and from who they are following
func (repository PostRepositoryMock) Index(userID uint64) ([]model.Post, error) {
	return repository.getStoredPostList()
}

// Update updates a post in database
func (repository PostRepositoryMock) Update(postID uint64, post model.Post) error {
	return nil
}

// Update deletes a post from database
func (repository PostRepositoryMock) Delete(postID uint64) error {
	return nil
}

// FindByUser returns all posts from a given user
func (repository PostRepositoryMock) FindByUser(userID uint64) ([]model.Post, error) {
	return []model.Post{}, nil
}

// LikePost increases the number of likes in a post
func (repository PostRepositoryMock) LikePost(postID uint64) error {
	return nil
}

// DeslikePost decreases the number of likes in a post
func (repository PostRepositoryMock) DeslikePost(postID uint64) error {
	return nil
}

func (repository PostRepositoryMock) getStoredPostList() ([]model.Post, error) {
	storedPostListJson, _ := ioutil.ReadFile("../test/resource/json/stored_post_list.json")

	var storedPostList []model.Post

	json.Unmarshal(storedPostListJson, &storedPostList)

	return storedPostList, nil
}

func (repository PostRepositoryMock) getStoredPost() (model.Post, error) {
	storedPostJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var storedPost model.Post

	json.Unmarshal(storedPostJson, &storedPost)

	return storedPost, nil
}
