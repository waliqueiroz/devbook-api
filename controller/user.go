package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/waliqueiroz/devbook-api/authentication"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/response"
	"github.com/waliqueiroz/devbook-api/security"
)

type UserController struct {
	userRepository *repository.UserRepository
}

// NewUserController creates a new UserController
func NewUserController(userRepository *repository.UserRepository) *UserController {
	return &UserController{
		userRepository,
	}
}

// Index shows all users
func (controller UserController) Index(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	users, err := controller.userRepository.FindByNameOrNick(nameOrNick)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// Create an user
func (controller UserController) Create(w http.ResponseWriter, r *http.Request) {
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

	if err = user.Prepare("register"); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	newUser, err := controller.userRepository.Create(user)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, newUser)
}

// Show returns a specific user
func (controller UserController) Show(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userRepository.FindByID(userID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// Updates a user
func (controller UserController) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	tokenUserID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != tokenUserID {
		response.Error(w, http.StatusForbidden, errors.New("is not possible to update an user other than your own"))
		return
	}

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

	if err = user.Prepare("update"); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = controller.userRepository.Update(userID, user)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// Deletes a user
func (controller UserController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	tokenUserID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != tokenUserID {
		response.Error(w, http.StatusForbidden, errors.New("is not possible to delete an user other than your own"))
		return
	}

	err = controller.userRepository.Delete(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// FollowUser allows an user to unfollow another
func (controller UserController) FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if userID == followerID {
		response.Error(w, http.StatusForbidden, errors.New("is not possible to follow yourself"))
		return
	}

	if err := controller.userRepository.Follow(userID, followerID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// UnfollowUser allows an user to unfollow another
func (controller UserController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if userID == followerID {
		response.Error(w, http.StatusForbidden, errors.New("is not possible to unfollow yourself"))
		return
	}

	if err := controller.userRepository.Unfollow(userID, followerID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// SearchFollowers returns a list of followers for a given user
func (controller UserController) SearchFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	followers, err := controller.userRepository.SearchFollowers(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}

// SearchFollowing returns a list of users that a given user is following
func (controller UserController) SearchFollowing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	followers, err := controller.userRepository.SearchFollowing(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}

// UpdatePassword updates the user password
func (controller UserController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	tokenUserID, err := authentication.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if userID != tokenUserID {
		response.Error(w, http.StatusForbidden, errors.New("you cannot update a user other than your own"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var password model.Password
	err = json.Unmarshal(body, &password)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := controller.userRepository.FindPassword(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err := security.Verify(hashedPassword, password.Current); err != nil {
		response.Error(w, http.StatusUnauthorized, errors.New("the current password does not match the one saved in the database"))
		return
	}

	newHasedPassword, err := security.Hash(password.New)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = controller.userRepository.UpdatePassword(userID, string(newHasedPassword))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
