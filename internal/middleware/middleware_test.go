package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func TestLogErrorsReturnsJSONErrorHandler(t *testing.T) {
	log := zap.NewNop()
	handler := LogErrors(log)
	if handler == nil {
		t.Fatal("expected error handler")
	}
}

func TestRequestContextInjectsRequestID(t *testing.T) {
	log := zap.NewNop()
	app := fiber.New()
	app.Use(RequestContext(log))
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString(c.Get(fiber.HeaderXRequestID))
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", bytes.NewBuffer(nil))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if resp.Header.Get(fiber.HeaderXRequestID) == "" {
		t.Fatal("expected X-Request-Id header to be set")
	}
}
