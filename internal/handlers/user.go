package handlers

import (
	"encoding/json"
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

    // Validate input
    if user.Name == "" || user.Email == "" {
        http.Error(w, "Name and Email are required", http.StatusBadRequest)
        return
    }

    // Ensure the table exists
    if err := ensureUsersTableExists(); err != nil {
        logger.Log.Error().Err(err).Msg("Failed to ensure users table exists") // Corrected logger usage
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Insert user into the database
    _, err := db.DB.NewInsert().Model(&user).Exec(r.Context())
    if err != nil {
        logger.Log.Error().Err(err).Msg("Failed to insert user") // Corrected logger usage
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

// func updateUser(w http.ResponseWriter,r *http.Request){
// 	if r.Method != http.MethodPut {
// 		http.Error(w,"invalid request path", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	path := strings.TrimPrefix(r.URL.Path, "/users")
// 	if path == "" {
// 		http.Error(w, "UserId is required",http.StatusBadRequest)
// 	}
// 	userId := path

// 	// Decode the request body to get the new user data

// 	var updates map[string]interface{}
// 	err := json.NewDecoder(r.Body).Decode(&updates)

// 	if err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
//         return
// 	}

// 	if len(updates) == 0 {
// 		http.Error(w, "No Fileds to update",http.StatusBadRequest)
// 		return
// 	}

// 	updateQuery := db.DB.NewUpdate().Model(&User{}).Where("id = ?", userId)

// 	if name, ok := updates["name"]; ok && name != "" {
// 		updateQuery.Set("name = ?", name)
// 	}
// 	if email, ok := updates["email"]; ok && email != "" {
// 		updateQuery.Set("email = ?", email)
// 	}
// 	 result, err := updateQuery.Exec(r.Context())
//   	if err != nil {
//         logger.Log.Error().Err(err).Msg("Failed to update user")
//         http.Error(w, "Internal server error", http.StatusInternalServerError)
//         return
//     }

//     rowsAffected, _ := result.RowsAffected()
//     if rowsAffected == 0 {
//         http.Error(w, "User not found", http.StatusNotFound)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
// }

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
