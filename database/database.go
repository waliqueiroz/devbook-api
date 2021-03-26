package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/waliqueiroz/devbook-api/config"
)

// Connect opens a connection with the database
func Connect() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBDatabase)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
