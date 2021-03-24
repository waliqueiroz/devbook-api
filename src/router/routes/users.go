package routes

import (
	"net/http"

	"github.com/waliqueiroz/devbook-api/src/controllers"
)

var usersRoutes = []Route{
	{
		URI:          "/users",
		Method:       http.MethodPost,
		Function:     controllers.Create,
		RequiresAuth: false,
	},
	{
		URI:          "/users",
		Method:       http.MethodGet,
		Function:     controllers.Index,
		RequiresAuth: false,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodGet,
		Function:     controllers.Show,
		RequiresAuth: false,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodPut,
		Function:     controllers.Update,
		RequiresAuth: false,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodDelete,
		Function:     controllers.Delete,
		RequiresAuth: false,
	},
}
