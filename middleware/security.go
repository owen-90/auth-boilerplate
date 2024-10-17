package middleware

import (
	"auth-boilerplate/models"
	"auth-boilerplate/utils"
	"strings"

	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

// JWT secret key
var jwtSecret []byte

func init() {
	utils.ReadENV()

	jwtSecret = []byte(os.Getenv("JWT_TOKEN")) // defined as a byte
}

// GenerateJWT Generate a token for the user
func GenerateJWT(user models.User) (string, error) {
	// Define the claims (user data and expiration time)
	claims := jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create a new JWT token with the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}

// AuthMiddleware validates the  JWT token in the Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.HandleError(w, http.StatusUnauthorized, "Unauthorized. Authorization header missing.")
			return
		}

		// Split "Bearer {token}" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid{
			utils.HandleError(w, http.StatusUnauthorized, "Invalid or expired token.")
		}

		next.ServeHTTP(w, r)
	})
}
