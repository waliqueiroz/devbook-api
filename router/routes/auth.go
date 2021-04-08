package routes

import (
	"net/http"

	"github.com/waliqueiroz/devbook-api/controller"
	"github.com/waliqueiroz/devbook-api/router"
)

func Auth(authController *controller.AuthController) []router.Route {
	return []router.Route{
		{
			URI:          "/login",
			Method:       http.MethodPost,
			Function:     authController.Login,
			RequiresAuth: false,
		},
	}
}
