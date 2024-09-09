package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func SetupDB() {
    // Configure zerolog
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    dsn := ""
    sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
    DB = bun.NewDB(sqldb, pgdialect.New())

    // Test the connection
    if err := DB.Ping(); err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to the database")
    }

    // Log a success message
    log.Info().Msg("Successfully connected to the database")
}
