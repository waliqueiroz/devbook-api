package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/waliqueiroz/devbook-api/authentication"
	"github.com/waliqueiroz/devbook-api/interfaces"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/response"
	"github.com/waliqueiroz/devbook-api/security"
)

type AuthController struct {
	userRepository interfaces.UserRepository
}

// NewAuthController creates a new AuthController
func NewAuthController(userRepository interfaces.UserRepository) *AuthController {
	return &AuthController{
		userRepository,
	}
}

// Login authenticates an user
func (controller AuthController) Login(w http.ResponseWriter, r *http.Request) {
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

	storedUser, err := controller.userRepository.FindByEmail(user.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	err = security.Verify(storedUser.Password, user.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.CreateToken(storedUser.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))

}
