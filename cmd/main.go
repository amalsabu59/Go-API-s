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
	mux.HandleFunc("/users/", usersHandler) 

	logger.Log.Info().Msg("Server starting at :8080")
	http.ListenAndServe(":8080", mux)
}

// usersHandler differentiates between HTTP methods for /users and /users/{userId}
func usersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodPost:
			handlers.AddUser(w, r) 
		case http.MethodGet:
			handlers.GetAllUsers(w, r) 
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserById(w, r) 
		case http.MethodPut:
			handlers.UpdateUser(w, r) 
		case http.MethodDelete:
			handlers.DeleteUser(w, r) 
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
