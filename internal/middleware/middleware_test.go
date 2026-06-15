package middleware

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogErrorsReturnsJSONErrorHandler(t *testing.T) {
	log := zap.NewNop()
	handler := LogErrors(log)
	if handler == nil {
		t.Fatal("expected error handler")
	}
}
