package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"user-api/config"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/middleware"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"
)

func main() {
	cfg := config.Load()

	zapLogger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = zapLogger.Sync() }()

	db, err := repository.OpenDB(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.LogErrors(zapLogger),
	})
	app.Use(middleware.RequestContext(zapLogger))
	app.Use(middleware.Recover(zapLogger))

	routes.Register(app, userHandler)

	middleware.LogStartup(zapLogger, cfg.Address())
	if err := app.Listen(cfg.Address()); err != nil {
		log.Fatal(err)
	}
}
