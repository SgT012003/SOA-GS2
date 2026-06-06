package dao

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/agriai?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("Falha ao inicializar banco de dados", slog.String("error", err.Error()))
		return err
	}

	if err = DB.Ping(); err != nil {
		slog.Error("Falha no ping do banco de dados", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Successfully connected to the database")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
