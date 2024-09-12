package routes

import (
	"net/http"

	"github.com/amalsabu59/onboard/internal/handlers"
)

func ProductRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/product", handlers.AddProduct)
}
