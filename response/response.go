package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON write a json response to a request
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal(err)
	}
}

// JSON write an error in json format to a request
func Error(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, struct {
		Erro string `json:"error"`
	}{
		Erro: err.Error(),
	})
}
