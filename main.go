package main

import (
	"fmt"
	"net/http"

	"github.com/waliqueiroz/devbook-api/config"
	"github.com/waliqueiroz/devbook-api/router"
)

func main() {
	config.Load()

	fmt.Printf("Running API on port %d...\n", config.APIPort)
	r := router.Generate()

	http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), r)
}
