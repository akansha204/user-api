package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recover(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					zap.Any("panic", r),
					zap.ByteString("stack", debug.Stack()),
				)
				err = fiber.NewError(fiber.StatusInternalServerError, "internal server error")
			}
		}()
		return c.Next()
	}
}

func LogErrors(log *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := err.Error()

		if fe, ok := err.(*fiber.Error); ok {
			code = fe.Code
			message = fe.Message
		}

		log.Error("http request failed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", code),
			zap.String("error", message),
		)

		return c.Status(code).JSON(fiber.Map{"error": message})
	}
}

func LogStartup(log *zap.Logger, addr string) {
	log.Info(fmt.Sprintf("server listening on %s", addr))
}
