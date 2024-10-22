package configs

import (
	"context"
	"log"
	"os"

	redis "github.com/gofiber/storage/redis/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// run migrations
	m, err := migrate.New("file://migrations/", DSN)
	if err != nil {
		log.Fatalf("Error running database migrations: %f", err)
	}
	if err := m.Up(); err != nil {
		log.Printf("migrations: %s", err)
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
