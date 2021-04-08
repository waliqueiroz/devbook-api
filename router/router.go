package router

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

// Generate will return a router with the configured routes
func Generate(applicationRoutes []Route) *mux.Router {
	r := mux.NewRouter()
	return config(r, applicationRoutes)
}

// Config put all the routes inside router
func config(r *mux.Router, applicationRoutes []Route) *mux.Router {

	for _, route := range applicationRoutes {

		if route.RequiresAuth {
			r.HandleFunc(route.URI, middleware.Logger(middleware.Authenticate(route.Function))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middleware.Logger(route.Function)).Methods(route.Method)
		}
	}

	return r
}
