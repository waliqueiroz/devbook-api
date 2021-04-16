package controller_test

import (
	"encoding/json"
	"errors"
	"io"
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

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

// TestCreateUser run test for user creation
func TestCreateUser(t *testing.T) {
	inputUserJson, _ := ioutil.ReadFile("../test/resource/json/input_user.json")
	invalidInputUserJson, _ := ioutil.ReadFile("../test/resource/json/invalid_input_user.json")
	incompleteInputUserJson, _ := ioutil.ReadFile("../test/resource/json/incomplete_input_user.json")

	var expectedUser model.User
	json.Unmarshal(inputUserJson, &expectedUser)

	subTests := []struct {
		name               string
		input              io.Reader
		expectedStatusCode int
		expectedResponse   model.User
	}{
		{
			name:               "Create valid user",
			input:              strings.NewReader(string(inputUserJson)),
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   expectedUser,
		},
		{
			name:               "Create invalid body payload",
			input:              errReader(0),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Create invalid user",
			input:              strings.NewReader(string(invalidInputUserJson)),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Create incomplete user",
			input:              strings.NewReader(string(incompleteInputUserJson)),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	userRepository := mock.NewUserRepositoryMock()
	userController := controller.NewUserController(userRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/users", subTest.input)
			request.Header.Add("Content-Type", "application/json")

			response := httptest.NewRecorder()

			userController.Create(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Result().StatusCode, "Status code does not match with expected")

			if subTest.expectedStatusCode == http.StatusCreated {
				var createdUser model.User

				json.Unmarshal(response.Body.Bytes(), &createdUser)

				assert.Equal(t, subTest.expectedResponse.Name, createdUser.Name, "User name does not match with expected")
				assert.Equal(t, subTest.expectedResponse.Email, createdUser.Email, "User email does not match with expected")
				assert.Equal(t, subTest.expectedResponse.Nick, createdUser.Nick, "User nick does not match with expected")
			} else {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}
}
