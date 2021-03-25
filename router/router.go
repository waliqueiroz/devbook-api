package router

import (
	"github.com/gorilla/mux"
	"github.com/waliqueiroz/devbook-api/router/routes"
)

// Generate will return a router with the configured routes
func Generate() *mux.Router {
	r := mux.NewRouter()
	return routes.Config(r)
}
