package controller_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/devbook-api/controller"
	"github.com/waliqueiroz/devbook-api/test/mock"
)

func TestLogin(t *testing.T) {
	loginInput, _ := ioutil.ReadFile("../test/resource/json/login_input.json")
	invalidCredentials, _ := ioutil.ReadFile("../test/resource/json/login_input_with_invalid_credentials.json")
	invalidLoginInput, _ := ioutil.ReadFile("../test/resource/json/invalid_login_input.json")

	subTests := []struct {
		name               string
		input              io.Reader
		expectedStatusCode int
	}{
		{
			name:               "Login with correct credentials",
			input:              bytes.NewReader(loginInput),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Login with invalid credentials",
			input:              bytes.NewReader(invalidCredentials),
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Login with invalid body payload",
			input:              mock.NewReader(),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Login with invalid data",
			input:              bytes.NewReader(invalidLoginInput),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	userRepository := mock.NewUserRepository()
	authController := controller.NewAuthController(userRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/login", subTest.input)
			request.Header.Add("Content-Type", "application/json")

			response := httptest.NewRecorder()

			authController.Login(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")
			assert.NotEmpty(t, response.Body.String(), "Response body is empty")
		})
	}
}
