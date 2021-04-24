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

func TestIndex(t *testing.T) {
	expectedPostListJson, _ := ioutil.ReadFile("../test/resource/json/stored_post_list.json")

	var expectedPostList []model.Post
	json.Unmarshal(expectedPostListJson, &expectedPostList)

	userID := uint64(1)
	token, _ := authentication.CreateToken(userID)

	subTests := []struct {
		name               string
		expectedStatusCode int
		expectedResponse   []model.Post
		token              string
	}{
		{
			name:               "Get posts",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   expectedPostList,
			token:              token,
		},
		{
			name:               "Get posts with invalid token",
			expectedStatusCode: http.StatusUnauthorized,
			token:              "teste=",
		},
	}

	postRepository := mock.NewPostRepository()
	postController := controller.NewPostController(postRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/posts", nil)
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", subTest.token))

			response := httptest.NewRecorder()

			postController.Index(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")

			if subTest.expectedStatusCode == http.StatusOK {
				var postList []model.Post
				json.Unmarshal(response.Body.Bytes(), &postList)
				assert.Equal(t, subTest.expectedResponse, postList, "Post list does not match with expected")
			} else {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}
}

func TestShowPost(t *testing.T) {
	expectedPostJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var expectedPost model.Post
	json.Unmarshal(expectedPostJson, &expectedPost)

	userID := uint64(1)

	subTests := []struct {
		name               string
		routeVariable      string
		expectedStatusCode int
		expectedResponse   model.Post
	}{
		{
			name:               "Get post with a valid user ID",
			routeVariable:      fmt.Sprintf("%d", userID),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   expectedPost,
		},
		{
			name:               "Get post with an invalid user ID",
			routeVariable:      "teste",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	postRepository := mock.NewPostRepository()
	postController := controller.NewPostController(postRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/posts/"+subTest.routeVariable, nil)
			request = mux.SetURLVars(request, map[string]string{
				"postID": subTest.routeVariable,
			})
			request.Header.Add("Content-Type", "application/json")

			response := httptest.NewRecorder()

			postController.Show(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")

			if subTest.expectedStatusCode == http.StatusOK {
				var createdPost model.Post
				json.Unmarshal(response.Body.Bytes(), &createdPost)
				assert.Equal(t, subTest.expectedResponse, createdPost, "Post does not match with expected")
			} else {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	postInputJson, _ := ioutil.ReadFile("../test/resource/json/post_input.json")
	invalidPostInputJson, _ := ioutil.ReadFile("../test/resource/json/invalid_post_input.json")
	incompletePostInputJson, _ := ioutil.ReadFile("../test/resource/json/incomplete_post_input.json")

	expectedPostJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var expectedPost model.Post
	json.Unmarshal(expectedPostJson, &expectedPost)

	userID := uint64(1)
	token, _ := authentication.CreateToken(userID)

	subTests := []struct {
		name               string
		input              io.Reader
		expectedStatusCode int
		token              string
		expectedResponse   model.Post
	}{
		{
			name:               "Create post with valid data",
			input:              bytes.NewReader(postInputJson),
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   expectedPost,
			token:              token,
		},
		{
			name:               "Create post with invalid body payload",
			input:              mock.NewReader(),
			expectedStatusCode: http.StatusUnprocessableEntity,
			token:              token,
		},
		{
			name:               "Create post with invalid data",
			input:              bytes.NewReader(invalidPostInputJson),
			expectedStatusCode: http.StatusBadRequest,
			token:              token,
		},
		{
			name:               "Create post with incomplete data",
			input:              bytes.NewReader(incompletePostInputJson),
			expectedStatusCode: http.StatusBadRequest,
			token:              token,
		},
		{
			name:               "Create post with invalid token",
			input:              bytes.NewReader(postInputJson),
			expectedStatusCode: http.StatusUnauthorized,
			token:              "teste=",
		},
	}

	postRepository := mock.NewPostRepository()
	postController := controller.NewPostController(postRepository)

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/posts", subTest.input)
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", subTest.token))

			response := httptest.NewRecorder()

			postController.Create(response, request)

			assert.Equal(t, subTest.expectedStatusCode, response.Code, "Status code does not match with expected")

			if subTest.expectedStatusCode == http.StatusCreated {
				var createdPost model.Post
				json.Unmarshal(response.Body.Bytes(), &createdPost)
				assert.Equal(t, subTest.expectedResponse, createdPost, "Post does not match with expected")
			} else {
				assert.NotEmpty(t, response.Body.String(), "Response body is empty")
			}
		})
	}
}
