package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"time"
)

func GetPool() *pgxpool.Pool {
	dsn := os.Getenv("DB_POSTGRES_URL")
	if dsn == "" {
		log.Printf("DB_POSTGRES_URL is not set")
		dsn = "postgres://postgres:1@localhost:5432/currencies"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	return pool
}
