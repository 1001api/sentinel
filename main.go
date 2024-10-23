package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
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
		ProxyHeader: "X-Forwarded-For",
	})

	// Encrypt Cookie Config
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SALT"),
	}))

	// easen up cors
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
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

	// init class validator
	var validate = validator.New()
	_ = validate.RegisterValidation("timestamp", services.IsISO8601Date)

	// init sessions
	sessionStore := configs.InitSession(redisCon)
	stateStore := configs.InitStateSession(redisCon)

	// init repo
	repository := gen.New(db)

	// init services
	utilService := services.UtilServiceImpl{
		Validate: validate,
		IPReader: ipdbCon,
	}
	userService := services.UserServiceImpl{
		UtilService: &utilService,
		Repo:        repository,
	}
	subService := services.SubServiceImpl{
		Repo: repository,
	}
	authService := services.AuthServiceImpl{
		UtilService:  &utilService,
		UserService:  &userService,
		SessionStore: sessionStore,
		StateStore:   stateStore,
	}
	projectService := services.ProjectServiceImpl{
		SubService: &subService,
		Repo:       repository,
		DB:         db,
	}
	eventService := services.EventServiceImpl{
		UtilService: &utilService,
		Repo:        repository,
		SubService:  &subService,
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
		UserService:    &userService,
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
