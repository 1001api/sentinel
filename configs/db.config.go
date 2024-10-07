package configs

import (
	"context"
	"log"
	"os"

	redis "github.com/gofiber/storage/redis/v2"
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

func InitRedis() *redis.Storage {
	REDIS_URL := os.Getenv("REDIS_URL")

	store := redis.New(redis.Config{
		URL:   REDIS_URL,
		Reset: false,
	})

	return store
}
