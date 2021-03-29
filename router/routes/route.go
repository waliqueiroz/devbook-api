package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/waliqueiroz/devbook-api/middleware"
)

// Route represents all the routes in the API
type Route struct {
	URI          string
	Method       string
	Function     func(http.ResponseWriter, *http.Request)
	RequiresAuth bool
}

// Config put all the routes inside router
func Config(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, authRoutes...)

	for _, route := range routes {

		if route.RequiresAuth {
			r.HandleFunc(route.URI, middleware.Logger(middleware.Authenticate(route.Function))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middleware.Logger(route.Function)).Methods(route.Method)
		}
	}

	return r
}
