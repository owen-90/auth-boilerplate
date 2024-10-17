package utils

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// HashPassword Hash password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash Check password against bcrypt hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsStrongPassword Check password strength
func IsStrongPassword(password string) bool {
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
}
/*func IsStrongPassword(password string) bool {
	hasLowercase := false
	hasUppercase := false
	hasDigit := false
	hasSpecialChar := false

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowercase = true
		} else if unicode.IsUpper(char) {
			hasUppercase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			hasSpecialChar = true
		}
	}

	return hasLowercase && hasUppercase && hasDigit && hasSpecialChar
}*/

// IsValidEmail Validate email
func IsValidEmail(email string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return reg.MatchString(email)
}

// StringJoin Join strings
func StringJoin(arr []string, sep string) string {
	return strings.Join(arr, sep)
}

// HandleError Handle errors
func HandleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	// Create a response map
	response := map[string]string{"error": message}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as a JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// ReadENV Read .env file
func ReadENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	/*var myEnv map[string]string
	myEnv, err = godotenv.Read()
	user := myEnv["DB_USER"]*/
}
