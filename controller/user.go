package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/waliqueiroz/devbook-api/database"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/response"
)

type userController struct{}

// NewUserController create a new UserController
func NewUserController() *userController {
	return &userController{}
}

// Index show all users
func (controller userController) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listando usuários"))
}

// Create an user
func (controller userController) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repository.NewUserRepository(db)

	user.ID, err = repository.Create(user)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}

	response.JSON(w, http.StatusCreated, user)
}

// Show returns a specific user
func (controller userController) Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listando usuário"))
}

// Update an user
func (controller userController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Atualizando usuário"))
}

// Delete an user
func (controller userController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deletando usuário"))
}
