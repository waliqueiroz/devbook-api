package routes

import (
	"net/http"

	"github.com/waliqueiroz/devbook-api/controller"
)

var userController = controller.NewUserController()

var usersRoutes = []Route{
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
		RequiresAuth: false,
	},
	{
		URI:          "/users/{userId}",
		Method:       http.MethodGet,
		Function:     userController.Show,
		RequiresAuth: false,
	},
	{
		URI:          "/users/{userId}",
		Method:       http.MethodPut,
		Function:     userController.Update,
		RequiresAuth: false,
	},
	{
		URI:          "/users/{userId}",
		Method:       http.MethodDelete,
		Function:     userController.Delete,
		RequiresAuth: false,
	},
}
