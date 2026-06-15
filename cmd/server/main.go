package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"user-api/internal/handler"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"
)

func main() {
	app := fiber.New()

	// TODO: replace with config-driven DB initialization in the next phase.
	var repo *repository.UserRepository
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	routes.Register(app, userHandler)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
