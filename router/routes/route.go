package routes

import (
	"net/http"

	"github.com/gorilla/mux"
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
	routes := usersRoutes

	for _, route := range routes {
		r.HandleFunc(route.URI, route.Function).Methods(route.Method)
	}

	return r
}
