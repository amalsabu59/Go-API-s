package routes

import (
	"net/http"
	"strings"

	"github.com/amalsabu59/onboard/internal/handlers"
)

func UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/", usersHandler)
}

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
