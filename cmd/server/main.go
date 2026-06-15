package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/middleware"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"
)

func main() {
	zapLogger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = zapLogger.Sync() }()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.LogErrors(zapLogger),
	})
	app.Use(middleware.Recover(zapLogger))

	// TODO: replace with config-driven DB initialization in the next phase.
	var repo *repository.UserRepository
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	routes.Register(app, userHandler)

	middleware.LogStartup(zapLogger, ":3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
