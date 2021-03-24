package controllers

import "net/http"

// Index show all users
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listando usuários"))
}

// Create an user
func Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
}

// Show returns a specific user
func Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Listando usuário"))
}

// Update an user
func Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Atualizando usuário"))
}

// Delete an user
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deletando usuário"))
}
