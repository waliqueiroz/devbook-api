package routes

import (
	"net/http"

	"github.com/waliqueiroz/devbook-api/controller"
)

var authController = controller.NewAuthController()

var authRoutes = []Route{
	{
		URI:          "/login",
		Method:       http.MethodPost,
		Function:     authController.Login,
		RequiresAuth: false,
	},
}
