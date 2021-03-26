package main

import (
	"fmt"
	"net/http"

	"github.com/waliqueiroz/devbook-api/config"
	"github.com/waliqueiroz/devbook-api/router"
)

func main() {
	config.Load()

	r := router.Generate()

	fmt.Printf("Listening on port %d...\n", config.APIPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), r)
}
