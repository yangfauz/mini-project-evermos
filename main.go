package main

import (
	"log"
	"mini-project-evermos/configs"
	"mini-project-evermos/handlers"
	"mini-project-evermos/models/entities/migration"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/repositories"
	"mini-project-evermos/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Setup Configuration
	configuration := configs.New()

	// Setup Database
	database := configs.NewMysqlDatabase(configuration)

	// Setup Migration
	migration.Migration(database)

	// Setup Repository
	authRepository := repositories.NewAuthRepository(database)
	userRepository := repositories.NewUserRepository(database)
	categoryRepository := repositories.NewCategoryRepository(database)

	// Setup Service
	authService := services.NewAuthService(&authRepository, &userRepository)
	regionService := services.NewRegionService()
	categoryService := services.NewCategoryService(&categoryRepository)

	// Setup Handler
	authHandler := handlers.NewAuthHandler(&authService)
	regionHandler := handlers.NewRegionHandler(&regionService)
	categoryHandler := handlers.NewCategoryHandler(&categoryService)

	// Setup Fiber
	app := fiber.New(configs.NewFiberConfig())

	fiber.New(configs.NewFiberConfig())

	app.Use(recover.New())
	app.Use(cors.New())

	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${latency} ${status} - ${method} ${path}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(responder.ApiResponse{
			Status:  true,
			Message: "Mini Project Evermos",
			Error:   nil,
			Data:    nil,
		})
	})

	// Setup Routing
	authHandler.Route(app)
	regionHandler.Route(app)
	categoryHandler.Route(app)

	//Not Found in Last
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).JSON(responder.ApiResponse{
			Status:  false,
			Message: "NOT FOUND",
			Error:   &fiber.ErrNotFound.Message,
			Data:    nil,
		})
	})

	chanServer := make(chan os.Signal, 1)
	signal.Notify(chanServer, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	host := ":3000"
	go func() {
		<-chanServer

		log.Printf("Server is shutting down in the %s.", host)
		err := app.Shutdown()
		if err != nil {
			log.Printf("Error in shutting down the server: %v.", err)
		}
	}()

	log.Printf("Server is running in the %s.", host)
	log.Println("Press Ctrl + C to exit the server!")
	err := app.Listen(host)
	if err != nil {
		log.Printf("Error in running the server: %v.", err)
	}
}
