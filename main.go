package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/hubkudev/sentinel/configs"
	"github.com/hubkudev/sentinel/middlewares"
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

	// init class validator
	var validate = validator.New()
	_ = validate.RegisterValidation("timestamp", services.IsISO8601Date)

	// init sessions
	sessionStore := configs.InitSession(redisCon)
	stateStore := configs.InitStateSession(redisCon)

	// init repo
	userRepo := repositories.UserRepoImpl{DB: db}
	projectRepo := repositories.ProjectRepositoryImpl{DB: db}
	eventRepo := repositories.EventRepoImpl{DB: db}

	// init services
	utilService := services.UtilServiceImpl{
		Validate: validate,
	}
	userService := services.UserServiceImpl{
		UtilService: &utilService,
		UserRepo:    &userRepo,
	}
	authService := services.AuthServiceImpl{
		UtilService:  &utilService,
		UserService:  &userService,
		SessionStore: sessionStore,
		StateStore:   stateStore,
	}
	projectService := services.ProjectServiceImpl{
		ProjectRepo: &projectRepo,
	}
	eventService := services.EventServiceImpl{
		UtilService: &utilService,
		ProjectRepo: &projectRepo,
		EventRepo:   &eventRepo,
	}
	apiService := services.APIServiceImpl{
		ProjectService: &projectService,
		EventService:   &eventService,
	}
	webService := services.WebServiceImpl{
		UserService:    &userService,
		ProjectService: &projectService,
		EventService:   &eventService,
	}

	// init middlewares
	m := middlewares.MiddlewareImpl{
		UserRepo:       &userRepo,
		SessionStorage: sessionStore,
	}

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
