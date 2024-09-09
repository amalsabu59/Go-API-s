package main

import (
	"net/http"
	"strings"

	"github.com/amalsabu59/onboard/internal/db"
	"github.com/amalsabu59/onboard/internal/handlers"
	"github.com/amalsabu59/onboard/internal/logger"
)

func main() {
	logger.SetupLogger()
	db.SetupDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handlers.HelloHandler)
	mux.HandleFunc("/users/", usersHandler) // Register handler for /users/ (note the trailing slash)

	logger.Log.Info().Msg("Server starting at :8080")
	http.ListenAndServe(":8080", mux)
}

// usersHandler differentiates between HTTP methods for /users and /users/{userId}
func usersHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the path and determine if this is a request for /users or /users/{userId}
	path := strings.TrimPrefix(r.URL.Path, "/users/")

	// If the path is empty, it means it's a request to /users
	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodPost:
			handlers.AddUser(w, r) // Call AddUser if method is POST
		case http.MethodGet:
			handlers.GetAllUsers(w, r) // Call GetAllUsers if method is GET
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		// If the path is not empty, it's a request for /users/{userId}
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserById(w, r) // Call GetUserById to get a user by ID
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
