package configs

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDBCon() *pgxpool.Pool {
	DSN := os.Getenv("DSN")

	db, err := pgxpool.New(context.Background(), DSN)
	if err != nil {
		log.Fatalf("Error initializing db connection: %s", err.Error())
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Error ping-ing db: %s", err.Error())
	}

	return db
}
