package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/waliqueiroz/devbook-api/config"
	"github.com/waliqueiroz/devbook-api/controller"
	"github.com/waliqueiroz/devbook-api/database"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/router"
	"github.com/waliqueiroz/devbook-api/router/routes"
)

func main() {
	config.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)

	authController := controller.NewAuthController(userRepository)
	userController := controller.NewUserController(userRepository)
	postController := controller.NewPostController(postRepository)

	var applicationRoutes []router.Route

	applicationRoutes = append(applicationRoutes, routes.Auth(authController)...)
	applicationRoutes = append(applicationRoutes, routes.User(userController)...)
	applicationRoutes = append(applicationRoutes, routes.Post(postController)...)

	r := router.Generate(applicationRoutes)

	fmt.Printf("Listening on port %d...\n", config.APIPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), r)
}
