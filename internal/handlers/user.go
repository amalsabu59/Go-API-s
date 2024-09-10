package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/amalsabu59/onboard/internal/db"
	"github.com/amalsabu59/onboard/internal/logger"
)

type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func AddUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }


    if user.Name == "" || user.Email == "" {
        http.Error(w, "Name and Email are required", http.StatusBadRequest)
        return
    }

    if err := ensureUsersTableExists(); err != nil {
        logger.Log.Error().Err(err).Msg("Failed to ensure users table exists") 
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

  
    _, err := db.DB.NewInsert().Model(&user).Exec(r.Context())
    if err != nil {
        logger.Log.Error().Err(err).Msg("Failed to insert user")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
func GetAllUsers(w http.ResponseWriter,r *http.Request){
	  if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

	var users []User
	err := db.DB.NewSelect().Model(&users).Scan(r.Context())

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed t fetch user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(users)

}

func GetUserById(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path,"/users/")
	if path == "" {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	userId := path
	var user []User

	err := db.DB.NewSelect().Model(&user).Where("id = ?",userId).Scan(r.Context())

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch user with this Id")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter,r *http.Request){
	if r.Method != http.MethodPut {
		http.Error(w,"invalid request path", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	if path == "" {
		http.Error(w, "UserId is required",http.StatusBadRequest)
	}
	userId := path

	// Decode the request body to get the new user data

	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
	}

	if len(updates) == 0 {
		http.Error(w, "No Fileds to update",http.StatusBadRequest)
		return
	}

	updateQuery := db.DB.NewUpdate().Model(&User{}).Where("id = ?", userId)

	if name, ok := updates["name"]; ok && name != "" {
		updateQuery.Set("name = ?", name)
	}
	if email, ok := updates["email"]; ok && email != "" {
		updateQuery.Set("email = ?", email)
	}
      fmt.Println("Update Query:", updateQuery)
	 result, err := updateQuery.Exec(r.Context())
  	if err != nil {
        logger.Log.Error().Err(err).Msg("Failed to update user")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodDelete{
        http.Error(w,"invalid request path", http.StatusMethodNotAllowed)
        return
    }
    path := strings.TrimPrefix(r.URL.Path,"/users/")
    if path == "" {
        http.Error(w,"User Id is required", http.StatusBadRequest)
    }
    userId := path

    _, err := db.DB.NewDelete().Model((*User)(nil)).Where("id = ?", userId).Exec(r.Context())

    if err != nil {
        logger.Log.Error().Err(err).Msg("Failed to delete")
        http.Error(w,"Failed to delete",http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func ensureUsersTableExists() error {
    _, err := db.DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) NOT NULL UNIQUE
        );
    `)
    return err
}
