package util

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// ConnectDatabase initializes the MySQL database connection
func ConnectDatabase() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/test" // Update with your username, password, and database name
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connected to the database successfully!")
}
