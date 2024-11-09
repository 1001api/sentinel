package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/hubkudev/sentinel/configs"
	gen "github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/middlewares"
	"github.com/hubkudev/sentinel/routes"
	"github.com/hubkudev/sentinel/services"
	"github.com/joho/godotenv"
)

var (
	engine = html.New("./views", ".html")
)

func main() {
	godotenv.Load()

	app := fiber.New(fiber.Config{
		Views:       engine,
		ProxyHeader: "CF-Connecting-IP",
	})

	// Encrypt Cookie Config
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SALT"),
	}))

	// CORS Policy
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Logging
	// comment this code to disable log every endpoint hit.
	// See more: https://docs.gofiber.io/api/middleware/logger
	app.Use(logger.New(logger.Config{
		Format: "${yellow} [${time}] ${status} - ${method} ${path} ${latency}\n",
	}))

	app.Static("/static", "./views/public", fiber.Static{
		Compress:  true,
		ByteRange: true,
		MaxAge:    3600,
	})

	// initialize database connection
	db := configs.InitDBCon()
	redisCon := configs.InitRedis()
	ipdbCon := configs.InitIPDBCon()
	defer db.Close()
	defer redisCon.Close()
	defer ipdbCon.Close()

	app.Use("api/event/download", limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			path := c.OriginalURL()
			return !strings.Contains(path, "/api/event/download")
		},
		Max:        4,
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).SendString(`
				<div x-data="{ show: true }" x-show="show" x-init="setTimeout(() => show = false, 5000)">
					Too many request, please wait for 30 seconds
				</div>
			`)
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("CF-Connecting-IP")
		},
	}))

	// init class validator
	var validate = validator.New()
	_ = validate.RegisterValidation("timestamp", services.IsISO8601Date)
	_ = validate.RegisterValidation("password", services.IsStrongPassword)

	// init sessions
	sessionStore := configs.InitSession(redisCon)

	// init repo
	repository := gen.New(db)

	// init services
	utilService := services.InitUtilService(validate, ipdbCon)
	cacheService := services.InitCacheService(redisCon)
	downloadService := services.InitDownloadService(repository)
	userService := services.InitUserService(&utilService, repository)
	authService := services.InitAuthService(&utilService, &userService, sessionStore)
	projectService := services.InitProjectService(repository, db)
	eventService := services.InitEventService(&utilService, &cacheService, repository)
	apiService := services.InitAPIService(&projectService, &eventService, &downloadService, &cacheService)
	webService := services.InitWebService(&userService, &projectService, &eventService)

	// init middleware
	m := middlewares.InitMiddleware(&userService, sessionStore, &cacheService)

	// init routes
	routes.InitAuthRoute(app, &authService)
	routes.InitEventRoute(app, &eventService)
	routes.InitAPIRoute(app, &m, &apiService, &eventService)
	routes.InitWebRoute(app, &m, &webService)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	if err := app.Listen(fmt.Sprintf(":%s", PORT)); err != nil {
		log.Fatalf("Failed to serve the server at port: %s", PORT)
	}
}
