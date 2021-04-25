package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/mattn/go-sqlite3"
	"github.com/waliqueiroz/devbook-api/config"
)

// Connect opens a connection with MySQL database
func Connect() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBDatabase)

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

// ConnectSQLite opens a connection with SQLite database
// func ConnectSQLite() (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = db.Ping(); err != nil {
// 		db.Close()
// 		return nil, err
// 	}

// 	db.SetMaxOpenConns(1)

// 	return db, nil
// }
