package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDKey = "request_id"

func RequestContext(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := c.Get(fiber.HeaderXRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(fiber.HeaderXRequestID, requestID)
		c.Locals(requestIDKey, requestID)

		err := c.Next()

		log.Info("http request",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", time.Since(start)),
		)

		return err
	}
}

func Recover(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				requestID, _ := c.Locals(requestIDKey).(string)
				log.Error("panic recovered",
					zap.String("request_id", requestID),
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

		requestID, _ := c.Locals(requestIDKey).(string)
		log.Error("http request failed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", code),
			zap.String("error", message),
		)

		return c.Status(code).JSON(fiber.Map{"error": message, "requestId": requestID})
	}
}

func LogStartup(log *zap.Logger, addr string) {
	log.Info(fmt.Sprintf("server listening on %s", addr))
}
