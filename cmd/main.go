package main

import (
	"net/http"

	"github.com/amalsabu59/onboard/internal/db"
	"github.com/amalsabu59/onboard/internal/logger"
	"github.com/amalsabu59/onboard/internal/routes"
)

func main() {
	logger.SetupLogger()
	db.SetupDB()
	mux := http.NewServeMux()
	routes.UserRoutes(mux)

	logger.Log.Info().Msg("Server starting at :8080")
	http.ListenAndServe(":8080", mux)
}

