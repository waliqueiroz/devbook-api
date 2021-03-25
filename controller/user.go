package controller

import "net/http"

// UserController
type UserController struct{}

// NewUserController create a nre UserController
func NewUserController() *UserController {
	return &UserController{}
}

// Index show all users
func (u *UserController) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listando usuários"))
}

// Create an user
func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
}

// Show returns a specific user
func (u *UserController) Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listando usuário"))
}

// Update an user
func (u *UserController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Atualizando usuário"))
}

// Delete an user
func (u *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deletando usuário"))
}
