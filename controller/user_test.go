package controller_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/devbook-api/controller"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/test/mock"
)

// TestCreateUser run test for user creation
func TestCreateUser(t *testing.T) {
	inputUserJson, _ := ioutil.ReadFile("../test/resource/json/input_user.json")
	invalidInputUserJson, _ := ioutil.ReadFile("../test/resource/json/invalid_input_user.json")
	incompleteInputUserJson, _ := ioutil.ReadFile("../test/resource/json/incomplete_input_user.json")

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
			name:               "Create valid user",
			input:              bytes.NewReader(inputUserJson),
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   expectedUser,
		},
		{
			name:               "Create invalid body payload",
			input:              mock.NewReader(),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Create invalid user",
			input:              bytes.NewReader(invalidInputUserJson),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Create incomplete user",
			input:              bytes.NewReader(incompleteInputUserJson),
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
			name:               "Get with a invalid user ID",
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
