package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func loadDatabase(url string, schema string) *sql.DB {
	database, err := sql.Open(
		"pgx", fmt.Sprintf("%s&search_path=%s", url, schema),
	)

	if err != nil {
		slog.Error("could not connect to database", "url", url)
		os.Exit(1)
	}

	if err := database.Ping(); err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}

	return database
}
