package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/waliqueiroz/devbook-api/database"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/response"
	"github.com/waliqueiroz/devbook-api/security"
)

type authController struct{}

// NewAuthController creates a new authController
func NewAuthController() *authController {
	return &authController{}
}

// Login authenticates an user
func (controller authController) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repository.NewUserRepository(db)

	userStored, err := repository.FindByEmail(user.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	err = security.Verify(userStored.Password, user.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	w.Write([]byte("Logou"))

}
