package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/gofiber/storage/redis/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDBCon() *pgxpool.Pool {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	DSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

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
	dbHost := os.Getenv("REDIS_HOST")
	dbPort := os.Getenv("REDIS_PORT")

	DSN := fmt.Sprintf("redis://%s:%s", dbHost, dbPort)

	store := redis.New(redis.Config{
		URL:   DSN,
		Reset: false,
	})

	return store
}
