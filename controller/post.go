package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/waliqueiroz/devbook-api/authentication"
	"github.com/waliqueiroz/devbook-api/database"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/response"
)

type postController struct{}

func NewPostController() *postController {
	return &postController{}
}

// Index shows all posts by a user and from who they are following
func (controller postController) Index(w http.ResponseWriter, r *http.Request) {}

// Create a post
func (controller postController) Create(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repository.NewPostRepository(db)

	post.ID, err = repository.Create(post)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

// Show a post
func (controller postController) Show(w http.ResponseWriter, r *http.Request) {}

// Update a post
func (controller postController) Update(w http.ResponseWriter, r *http.Request) {}

// Delete a post
func (controller postController) Delete(w http.ResponseWriter, r *http.Request) {}
