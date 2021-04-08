package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/waliqueiroz/devbook-api/authentication"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/response"
)

type PostController struct {
	postRepository *repository.PostRepository
}

func NewPostController(postRepository *repository.PostRepository) *PostController {
	return &PostController{
		postRepository,
	}
}

// Index shows all posts by a user and from who they are following
func (controller PostController) Index(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	posts, err := controller.postRepository.Index(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, posts)

}

// Create creates a post
func (controller PostController) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post model.Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorID = userID

	if err := post.Prepare(); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post.ID, err = controller.postRepository.Create(post)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

// Show shows a post
func (controller PostController) Show(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	postID, err := strconv.ParseUint(params["postID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post, err := controller.postRepository.FindByID(postID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

// Update updates a post
func (controller PostController) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	postID, err := strconv.ParseUint(params["postID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	storedPost, err := controller.postRepository.FindByID(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if userID != storedPost.AuthorID {
		response.Error(w, http.StatusForbidden, errors.New("you cannot update a post that is not yours"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post model.Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = controller.postRepository.Update(postID, post)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)

}

// Delete deletes a post
func (controller PostController) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	postID, err := strconv.ParseUint(params["postID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	storedPost, err := controller.postRepository.FindByID(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if userID != storedPost.AuthorID {
		response.Error(w, http.StatusForbidden, errors.New("you cannot delete a post that is not yours"))
		return
	}

	err = controller.postRepository.Delete(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// FindByUser returns all posts from a given user
func (controller PostController) FindByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post, err := controller.postRepository.FindByUser(userID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

// LikePost increases the number of likes in a post
func (controller PostController) LikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	postID, err := strconv.ParseUint(params["postID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = controller.postRepository.LikePost(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)

}

// DeslikePost decreases the number of likes in a post
func (controller PostController) DeslikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	postID, err := strconv.ParseUint(params["postID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = controller.postRepository.DeslikePost(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)

}
