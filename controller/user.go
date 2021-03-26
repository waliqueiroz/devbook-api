package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/waliqueiroz/devbook-api/database"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
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
		log.Fatal(err)
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	repository := repository.NewUserRepository(db)

	userId, err := repository.Create(user)

	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf("Inserted ID: %d", userId)))
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
