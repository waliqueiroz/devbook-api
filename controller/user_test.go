package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/devbook-api/authentication"
	"github.com/waliqueiroz/devbook-api/controller"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/test/mock"
)

// TestCreateUser run test for user creation
func TestCreateUser(t *testing.T) {
	userInputJson, _ := ioutil.ReadFile("../test/resource/json/user_input.json")
	invalidUserInputJson, _ := ioutil.ReadFile("../test/resource/json/invalid_user_input.json")
	incompleteUserInputJson, _ := ioutil.ReadFile("../test/resource/json/incomplete_user_input.json")

	expectedUserJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var expectedUser model.User
	json.Unmarshal(expectedUserJson, &expectedUser)

	subTests := []struct {
		name               string
		input              io.Reader
		expectedStatusCode int
		expectedResponse   model.User
	}{
		{
			name:               "Create user with valid data",
			input:              bytes.NewReader(userInputJson),
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   expectedUser,
		},
		{
			name:               "Create user with invalid body payload",
			input:              mock.NewReader(),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Create user with invalid data",
			input:              bytes.NewReader(invalidUserInputJson),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Create user with incomplete data",
			input:              bytes.NewReader(incompleteUserInputJson),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	userRepository := mock.NewUserRepository()
	userController := controller.NewUserController(userRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/users", subTest.input)
			request.Header.Add("Content-Type", "application/json")

			response := httptest.NewRecorder()

			userController.Create(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")

			if subTest.expectedStatusCode == http.StatusCreated {
				var createdUser model.User
				json.Unmarshal(response.Body.Bytes(), &createdUser)

				assert.Equal(t, subTest.expectedResponse, createdUser, "Created user does not match with expected")
			} else {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}
}

func TestFindUsers(t *testing.T) {
	expectedUserListJson, _ := ioutil.ReadFile("../test/resource/json/stored_user_list.json")

	var expectedUserList []model.User
	json.Unmarshal(expectedUserListJson, &expectedUserList)

	userRepository := mock.NewUserRepository()
	userController := controller.NewUserController(userRepository)

	request := httptest.NewRequest("GET", "/users?user=Juliette", nil)
	request.Header.Add("Content-Type", "application/json")

	response := httptest.NewRecorder()

	userController.Index(response, request)

	var userList []model.User
	json.Unmarshal(response.Body.Bytes(), &userList)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode, "Status code does not match with expected")
	assert.Equal(t, expectedUserList, userList, "User list does not match with expected")
}

func TestShowUser(t *testing.T) {
	expectedUserJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var expectedUser model.User
	json.Unmarshal(expectedUserJson, &expectedUser)

	subTests := []struct {
		name               string
		routeVariable      string
		expectedStatusCode int
		expectedResponse   model.User
	}{
		{
			name:               "Get with a valid user ID",
			routeVariable:      "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   expectedUser,
		},
		{
			name:               "Get with an invalid user ID",
			routeVariable:      "teste",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	userRepository := mock.NewUserRepository()
	userController := controller.NewUserController(userRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/users/"+subTest.routeVariable, nil)
			request = mux.SetURLVars(request, map[string]string{
				"userID": subTest.routeVariable,
			})
			request.Header.Add("Content-Type", "application/json")

			response := httptest.NewRecorder()

			userController.Show(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")

			if subTest.expectedStatusCode == http.StatusOK {
				var createdUser model.User
				json.Unmarshal(response.Body.Bytes(), &createdUser)
				assert.Equal(t, subTest.expectedResponse, createdUser, "Created user does not match with expected")
			} else {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	userInputJson, _ := ioutil.ReadFile("../test/resource/json/user_input_update.json")
	invalidUserInputJson, _ := ioutil.ReadFile("../test/resource/json/invalid_user_input.json")
	incompleteUserInputJson, _ := ioutil.ReadFile("../test/resource/json/incomplete_user_input.json")

	expectedUserJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var expectedUser model.User
	json.Unmarshal(expectedUserJson, &expectedUser)

	userID := uint64(1)
	token, _ := authentication.CreateToken(userID)

	subTests := []struct {
		name               string
		input              io.Reader
		routeVariable      string
		expectedStatusCode int
		token              string
	}{
		{
			name:               "Update user with valid data",
			input:              bytes.NewReader(userInputJson),
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusNoContent,
			token:              token,
		},
		{
			name:               "Update user with an invalid authorization token",
			input:              bytes.NewReader(userInputJson),
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusUnauthorized,
			token:              "teste=",
		},
		{
			name:               "Update user with an invalid user ID",
			input:              bytes.NewReader(userInputJson),
			routeVariable:      "teste",
			expectedStatusCode: http.StatusBadRequest,
			token:              token,
		},
		{
			name:               "Update user with invalid body payload",
			input:              mock.NewReader(),
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusUnprocessableEntity,
			token:              token,
		},
		{
			name:               "Try to update a user other than your own",
			input:              bytes.NewReader(userInputJson),
			routeVariable:      "2",
			expectedStatusCode: http.StatusForbidden,
			token:              token,
		},
		{
			name:               "Update user with invalid data",
			input:              bytes.NewReader(invalidUserInputJson),
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusBadRequest,
			token:              token,
		},
		{
			name:               "Update user with incomplete data",
			input:              bytes.NewReader(incompleteUserInputJson),
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusBadRequest,
			token:              token,
		},
	}

	userRepository := mock.NewUserRepository()
	userController := controller.NewUserController(userRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("PUT", "/users/"+subTest.routeVariable, subTest.input)
			request = mux.SetURLVars(request, map[string]string{
				"userID": subTest.routeVariable,
			})
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", subTest.token))

			response := httptest.NewRecorder()

			userController.Update(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")

			if subTest.expectedStatusCode != http.StatusNoContent {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}

}

func TestDeleteUser(t *testing.T) {

	userID := uint64(1)
	token, _ := authentication.CreateToken(userID)

	subTests := []struct {
		name               string
		routeVariable      string
		expectedStatusCode int
		token              string
	}{
		{
			name:               "Delete user",
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusNoContent,
			token:              token,
		},
		{
			name:               "Delete user with an invalid token",
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusUnauthorized,
			token:              "teste=",
		},
		{
			name:               "Delete user with an invalid user ID",
			routeVariable:      "teste",
			expectedStatusCode: http.StatusBadRequest,
			token:              token,
		},
		{
			name:               "Try to delete a user other than your own",
			routeVariable:      "2",
			expectedStatusCode: http.StatusForbidden,
			token:              token,
		},
	}

	userRepository := mock.NewUserRepository()
	userController := controller.NewUserController(userRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("PUT", "/users/"+subTest.routeVariable, nil)
			request = mux.SetURLVars(request, map[string]string{
				"userID": subTest.routeVariable,
			})
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", subTest.token))

			response := httptest.NewRecorder()

			userController.Delete(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")

			if subTest.expectedStatusCode != http.StatusNoContent {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}
}
