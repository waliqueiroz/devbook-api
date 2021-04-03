package routes

import (
	"net/http"

	"github.com/waliqueiroz/devbook-api/controller"
)

var postController = controller.NewPostController()

var postRoutes = []Route{
	{
		URI:          "/posts",
		Method:       http.MethodPost,
		Function:     postController.Create,
		RequiresAuth: false,
	},
	{
		URI:          "/posts",
		Method:       http.MethodGet,
		Function:     postController.Index,
		RequiresAuth: true,
	},
	{
		URI:          "/posts/{postID}",
		Method:       http.MethodGet,
		Function:     postController.Show,
		RequiresAuth: true,
	},
	{
		URI:          "/posts/{postID}",
		Method:       http.MethodPut,
		Function:     postController.Update,
		RequiresAuth: true,
	},
	{
		URI:          "/posts/{postID}",
		Method:       http.MethodDelete,
		Function:     postController.Delete,
		RequiresAuth: true,
	},
}
