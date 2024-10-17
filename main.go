package main

import (
	"auth-boilerplate/controllers"
	"auth-boilerplate/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"net/http"
)

func main() {

	router := mux.NewRouter()

	// Subrouter for API routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Test route
	apiRouter.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test endpoint reached"))
	}).Methods("GET")

	// Test route
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	// Switched routes
	http.HandleFunc("/register", controllers.RegisterUser)

	//apiRouter.HandleFunc("/", controllers.HomeHandler)
	apiRouter.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	apiRouter.HandleFunc("/login", controllers.Login).Methods("POST")

	// Subrouter for protected API routes (authentication required)
	protectedAPIRouter := apiRouter.PathPrefix("/").Subrouter()
	protectedAPIRouter.Use(middleware.AuthMiddleware)

	// Protected routes (require JWT authentication)

	http.HandleFunc("/test", controllers.HomeHandler)

	// Start the server
	fmt.Println("Starting server on :8025...")
	if err := http.ListenAndServe(":8025", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
