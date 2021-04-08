package routes

import (
	"net/http"

	"github.com/waliqueiroz/devbook-api/controller"
	"github.com/waliqueiroz/devbook-api/router"
)

func Post(postController *controller.PostController) []router.Route {
	return []router.Route{
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
		{
			URI:          "/users/{userID}/posts",
			Method:       http.MethodGet,
			Function:     postController.FindByUser,
			RequiresAuth: true,
		},
		{
			URI:          "/posts/{postID}/like",
			Method:       http.MethodPost,
			Function:     postController.LikePost,
			RequiresAuth: true,
		},
		{
			URI:          "/posts/{postID}/deslike",
			Method:       http.MethodPost,
			Function:     postController.DeslikePost,
			RequiresAuth: true,
		},
	}
}
