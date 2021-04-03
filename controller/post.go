package controller

import "net/http"

type postController struct{}

func NewPostController() *postController {
	return &postController{}
}

// Index shows all posts by a user and from who they are following
func (controller postController) Index(w http.ResponseWriter, r *http.Request) {}

// Create a post
func (controller postController) Create(w http.ResponseWriter, r *http.Request) {}

// Show a post
func (controller postController) Show(w http.ResponseWriter, r *http.Request) {}

// Update a post
func (controller postController) Update(w http.ResponseWriter, r *http.Request) {}

// Delete a post
func (controller postController) Delete(w http.ResponseWriter, r *http.Request) {}
