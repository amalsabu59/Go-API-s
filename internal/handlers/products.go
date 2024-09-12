package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amalsabu59/onboard/internal/db"
	"github.com/amalsabu59/onboard/internal/logger"
)

type Product struct {
    ID           int64     `bun:",pk,autoincrement" json:"id,omitempty"` // Bun recognizes ID as primary key and auto-increment
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Price        float64   `json:"price"`
    Category     string    `json:"category"`
    ImageURL     string    `json:"image_url"`
    StockQuantity int64    `json:"stock_quantity"`
    IsDeleted    bool      `json:"is_deleted"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}


func AddProduct(w http.ResponseWriter, r *http.Request) {
    fmt.Println("entering add product")
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var product Product

    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    product.ID = 0

    if product.Name == "" || product.Description == "" || product.Price <= 0 {
        http.Error(w, "Name, Description, and Price are required", http.StatusBadRequest)
        return
    }
    if err := ensureProductTableExists(); err != nil {
        logger.Log.Error().Err(err).Msg("Failed to ensure products table exists")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    _, err := db.DB.NewInsert().
        Model(&product).
        Exec(r.Context())

    if err != nil {
        logger.Log.Error().Err(err).Msg("Failed to insert product")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}



func GetAllProducts(w http.ResponseWriter, r *http.Request) {
dropProductsTable()	
	dropProductsTable()
	if r.Method != http.MethodGet {
		http.Error(w, "Internal server error", http.StatusMethodNotAllowed)
		return
	}

	var products []Product

	err := db.DB.NewSelect().Model(&products).Scan(r.Context())

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to select products")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)

}

func dropProductsTable() error {
	_, err := db.DB.Exec("DROP TABLE IF EXISTS products;")
	if err != nil {
		return err
	}
	log.Println("Products table dropped successfully")
	return nil
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
);
`,
	)
	return err
}



