package db

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func Connect() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, reading from environment")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASS")
	name := os.Getenv("POSTGRES_DB_NAME")

	if host == "" || port == "" || user == "" || name == "" {
		return nil, fmt.Errorf("DB environment variables not fully set")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)

	var db *sqlx.DB
	db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open DB connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	return db, nil
}
