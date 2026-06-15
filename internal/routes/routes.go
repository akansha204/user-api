package routes

import (
	"github.com/gofiber/fiber/v2"

	"user-api/internal/handler"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {
	users := app.Group("/users")
	users.Post("", userHandler.CreateUser)
	users.Get("", userHandler.ListUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
