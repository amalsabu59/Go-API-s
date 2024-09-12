package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amalsabu59/onboard/internal/db"
	"github.com/amalsabu59/onboard/internal/logger"
)

type Product struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`       
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`     
	ImageURL     string  `json:"image_url"`   
	StockQuantity int64  `json:"stock_quantity"`
	IsDeleted    bool    `json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"` 
}


func AddProduct(w http.ResponseWriter ,r *http.Request){
	fmt.Println("entering add product")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method",http.StatusMethodNotAllowed)
		return
	}

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product) ; err != nil {
		http.Error(w, "Requied fields are missing",http.StatusBadRequest)
	}
	if product.Name == "" || product.Description  == "" || product.Price == 0 {
 	 http.Error(w, "Name Description and Price are required", http.StatusBadRequest)
	 return
	}

	if err := ensureProductTableExists(); err != nil {
	    logger.Log.Error().Err(err).Msg("Failed to ensure products table exists") 
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
	}

	_, err := db.DB.NewInsert().Model(&product).Exec(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to insert product")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func ensureProductTableExists() error {
	_, err := db.DB.Exec(
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			category VARCHAR(100) DEFAULT NULL,
			image_url TEXT DEFAULT NULL,
			stock_quantity INTEGER NOT NULL DEFAULT 1,
			is_deleted BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	)
	return err
}


