package controller_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/devbook-api/controller"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/test/mock"
)

// TestCreateUser run test for user creation
func TestCreateUser(t *testing.T) {
	inputUserJson, _ := ioutil.ReadFile("../test/resource/json/input_user.json")

	// subTests := []struct {
	// 	input              string
	// 	expectedStatusCode int
	// 	expectedResponse string
	// }{

	// }

	userRepository := mock.NewUserRepositoryMock()
	userController := controller.NewUserController(userRepository)

	request := httptest.NewRequest("POST", "/users", strings.NewReader(string(inputUserJson)))
	request.Header.Add("Content-Type", "application/json")

	response := httptest.NewRecorder()

	userController.Create(response, request)

	responseStatus := response.Result().StatusCode

	var inputUser model.User
	var createdUser model.User

	json.Unmarshal(inputUserJson, &inputUser)
	json.Unmarshal(response.Body.Bytes(), &createdUser)

	assert.Equal(t, http.StatusCreated, responseStatus)
	assert.Equal(t, inputUser.Name, createdUser.Name)
}
