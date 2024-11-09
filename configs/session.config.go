package configs

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v2"
)

func InitSession(storage *redis.Storage) *session.Store {
	return session.New(session.Config{
		KeyLookup:      "cookie:session_id",
		Expiration:     24 * time.Hour, // 24 hours
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookiePath:     "/",
		Storage:        storage,
	})
}
