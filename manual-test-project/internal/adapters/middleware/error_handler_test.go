package middleware

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"manual-test-project/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler_FiberError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "Resource not found")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Verify JSON contains expected structure
	assert.Contains(t, string(body), `"status":"error"`)
	assert.Contains(t, string(body), `"message":"Resource not found"`)
	assert.Contains(t, string(body), `"code":"NOT_FOUND"`)
	assert.Contains(t, string(body), `"details":null`)
}

func TestErrorHandler_AppError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	details := map[string]string{"field": "email", "issue": "invalid format"}
	app.Get("/test", func(c *fiber.Ctx) error {
		return domain.NewBadRequestError("Invalid input", "INVALID_INPUT", details)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Contains(t, string(body), `"status":"error"`)
	assert.Contains(t, string(body), `"message":"Invalid input"`)
	assert.Contains(t, string(body), `"code":"INVALID_INPUT"`)
	assert.Contains(t, string(body), `"details":`)
	assert.Contains(t, string(body), `"email"`)
}

func TestErrorHandler_UnknownError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return errors.New("unexpected error")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Contains(t, string(body), `"status":"error"`)
	assert.Contains(t, string(body), `"code":"INTERNAL_SERVER_ERROR"`)
	assert.Contains(t, string(body), `"message":"Internal server error"`)
	assert.Contains(t, string(body), `"details":null`)
}

func TestErrorHandler_MethodNotAllowed(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusMethodNotAllowed, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Contains(t, string(body), `"status":"error"`)
	assert.Contains(t, string(body), `"code":"METHOD_NOT_ALLOWED"`)
	assert.Contains(t, string(body), `"details":null`)
}
