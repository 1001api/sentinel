package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/hubkudev/sentinel/configs"
	repositories "github.com/hubkudev/sentinel/repos"
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
		Views: engine,
	})

	app.Static("/static", "./views/public", fiber.Static{
		Compress:  true,
		ByteRange: true,
		MaxAge:    3600,
	})

	// initialize database connection
	db := configs.InitDBCon()
	redisCon := configs.InitRedis()
	defer db.Close()
	defer redisCon.Close()

	// init sessions
	sessionStore := configs.InitSession(redisCon)
	stateStore := configs.InitStateSession(redisCon)

	// init repo
	userRepo := repositories.UserRepoImpl{DB: db}

	// init services
	utilService := services.UtilServiceImpl{}
	userService := services.UserServiceImpl{
		UserRepo: &userRepo,
	}
	authService := services.AuthServiceImpl{
		UtilService:  &utilService,
		UserService:  &userService,
		SessionStore: sessionStore,
		StateStore:   stateStore,
	}
	eventService := services.EventServiceImpl{}

	// init routes
	routes.InitAuthRoute(app, &authService)
	routes.InitEventRoute(app, &eventService)
	routes.InitWebRoute(app)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	if err := app.Listen(fmt.Sprintf(":%s", PORT)); err != nil {
		log.Fatalf("Failed to serve the server at port: %s", PORT)
	}
}
