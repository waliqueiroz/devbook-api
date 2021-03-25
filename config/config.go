package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var DBHost = ""
var DBPort = ""
var DBDatabase = ""
var DBUsername = ""
var DBPassword = ""
var APIPort = 0

func Load() {
	var err error

	err = godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	APIPort, err = strconv.Atoi(os.Getenv("API_PORT"))

	if err != nil {
		APIPort = 8000
	}

	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBDatabase = os.Getenv("DB_DATABASE")
	DBUsername = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")

}
