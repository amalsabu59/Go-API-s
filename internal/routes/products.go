package routes

import (
	"net/http"

	"github.com/amalsabu59/onboard/internal/handlers"
)

func ProductRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/product", productRouteHandler)
}


func productRouteHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost {
		handlers.AddProduct(w, r)
	}
	if r.Method == http.MethodGet {
		handlers.GetAllProducts(w, r)
	}
}
