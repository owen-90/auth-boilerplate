package database

import (
	"auth-boilerplate/utils"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
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

func decrypt(encryptedText, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

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
	utils.ReadENV()

	// Read credentials from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := os.Getenv("DB_DSN") // Data source name := database name
	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	if user == "" || password == "" || dsn == "" || encryptionKey == "" {
		return fmt.Errorf("missing required environment variables")
	}

	// Decrypt the password
	password, err := decrypt(password, encryptionKey)
	if err != nil {
		return fmt.Errorf("error decrypting password: %w", err)
	}

	// Connect to the database
	err = Connect(user, password, dsn)
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
