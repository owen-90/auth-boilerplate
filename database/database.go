package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	initUser = `CREATE TABLE IF NOT EXISTS users (
id INT PRIMARY KEY AUTO_INCREMENT,
name VARCHAR(100) NOT NULL,
email VARCHAR(255) UNIQUE NOT NULL
);`
	initLogin = `CREATE TABLE IF NOT EXISTS login (
email VARCHAR(255) PRIMARY KEY,
username VARCHAR(255),
password VARCHAR(60) NOT NULL,
token VARCHAR(255) UNIQUE NOT NULL,
FOREIGN KEY (email) REFERENCES users (email)
);`
)
// Global DB instance
var db *sql.DB

// Connect to the database and assign the connection to the global `db` variable
func Connect(user, password, dsn string) error {
	var err error
	// Assign to the global `db` variable
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s", user, password, dsn))
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	fmt.Println("Connected to database")
	return nil
}

func Init() error {
	//utils.ReadENV()

	// Read credentials from environment variables
	user := os.Getenv("DB_USER")
	if user == "" {
		return fmt.Errorf("missing environment variable DB_USER")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return fmt.Errorf("missing environment variable DB_PASSWORD")
	}

	dsn := os.Getenv("DB_DSN") // Data source name := database name
	if dsn == "" {
		return fmt.Errorf("missing environment variable DB_DSN")
	}

	// Connect to the database
	err := Connect(user, password, dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	_, err = db.Exec(initUser)
	if err != nil {
		return fmt.Errorf("error creating USER table: %w", err)
	}

	_, err = db.Exec(initLogin)
	if err != nil {
		return fmt.Errorf("error creating LOGIN table: %w", err)
	}

	return nil
}

func DB() *sql.DB {
	if db == nil {
		// If not connected, attempt to connect on-demand
		err := Init()
		if err != nil {
			log.Fatal(err) // Handle connection error more gracefully in your application
		}
	}
	return db
}
