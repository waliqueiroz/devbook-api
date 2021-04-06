package routes

import (
	"net/http"

	"github.com/waliqueiroz/devbook-api/controller"
)

var userController = controller.NewUserController()

var userRoutes = []Route{
	{
		URI:          "/users",
		Method:       http.MethodPost,
		Function:     userController.Create,
		RequiresAuth: false,
	},
	{
		URI:          "/users",
		Method:       http.MethodGet,
		Function:     userController.Index,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}",
		Method:       http.MethodGet,
		Function:     userController.Show,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}",
		Method:       http.MethodPut,
		Function:     userController.Update,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}",
		Method:       http.MethodDelete,
		Function:     userController.Delete,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/follow",
		Method:       http.MethodPost,
		Function:     userController.FollowUser,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/unfollow",
		Method:       http.MethodPost,
		Function:     userController.UnfollowUser,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/followers",
		Method:       http.MethodGet,
		Function:     userController.SearchFollowers,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/following",
		Method:       http.MethodGet,
		Function:     userController.SearchFollowing,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/update-password",
		Method:       http.MethodPost,
		Function:     userController.UpdatePassword,
		RequiresAuth: true,
	},
}
