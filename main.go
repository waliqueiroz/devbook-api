package main

import (
	"fmt"
	"net/http"

	"github.com/waliqueiroz/devbook-api/src/router"
)

func main() {
	fmt.Println("Running API on port 5000...")
	r := router.Generate()

	http.ListenAndServe(":5000", r)
}
