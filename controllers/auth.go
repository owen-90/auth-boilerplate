package controllers

import (
	"auth-boilerplate/database"
	"auth-boilerplate/middleware"
	"auth-boilerplate/models"
	"auth-boilerplate/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// HomeHandler Home handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Home Page!")
}

// RegisterUser Register users
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body in a User struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate user input
	err = validateUser(user)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	user.Password = hashedPassword

	// Insert user into the DB
	stmt, err := database.DB().Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "error preparing user insertion statement")
		return
	}
	// Close the prepared statement
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user: %w", err))
		return
	}

	// Insert user login details
	_, err = database.DB().Exec("INSERT INTO login (email, username, password) VALUES (?, ?, ?)", user.Email,
		user.Username, user.Password)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "error creating user login details")
		return
	}

	// Create a success response
	response := map[string]interface{}{
		"message":    "User registered successfully",
		"statusCode": http.StatusOK,
		"user":       user,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	// Validate user input
	if (user.Username == "" && user.Email == "") || user.Password == "" {
		utils.HandleError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	// Retrieve the user from the DB
	var storedUser models.User

	query := "SELECT id, email, username, password FROM users WHERE email = ? AND username = ?"
	row := database.DB().QueryRow(query, user.Email, user.Username)

	err = row.Scan(&storedUser.ID, &storedUser.Username, &storedUser.Email, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleError(w, http.StatusNotFound, "User not found")
		} else {
			utils.HandleError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(user.Password, storedUser.Password) {
		utils.HandleError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	// Generate a JWT token
	token, err := middleware.GenerateJWT(storedUser)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Success response
	response := map[string]interface{}{
		"message":    "successfully logged in",
		"statusCode": http.StatusOK,
		"token":      token,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// User validation
func validateUser(user models.User) error {
	var validationErrors []string

	if user.Username == "" {
		validationErrors = append(validationErrors, "Username is required")
	} else if len(user.Username) < 5 {
		validationErrors = append(validationErrors, "Username must be at least 5 characters long")
	}

	if user.Password == "" {
		validationErrors = append(validationErrors, "Password is required")
	} else if len(user.Password) < 8 {
		validationErrors = append(validationErrors, "Password must be at least 8 characters long")
	} else if !utils.IsStrongPassword(user.Password) {
		validationErrors = append(validationErrors, "Password must contain at least one uppercase letter, "+
			"one lowercase letter, one digit, and one special character")
	}

	if user.Email == "" {
		validationErrors = append(validationErrors, "Email is required")
	} else if !utils.IsValidEmail(user.Email) {
		validationErrors = append(validationErrors, "Invalid email format")
	}

	if len(validationErrors) > 0 {
		return errors.New("Validation errors: " + utils.StringJoin(validationErrors, ", "))
	}

	return nil
}
