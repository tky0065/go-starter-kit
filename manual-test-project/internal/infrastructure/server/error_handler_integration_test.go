package server

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"manual-test-project/internal/adapters/middleware"
	"manual-test-project/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Details any    `json:"details"`
}

func TestErrorHandler_404NotFound(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/exists", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Request a route that doesn't exist
	req := httptest.NewRequest("GET", "/does-not-exist", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var errResp ErrorResponse
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err)

	// Verify JSON structure matches requirements
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "NOT_FOUND", errResp.Code)
	assert.NotEmpty(t, errResp.Message)
	assert.Nil(t, errResp.Details)
}

func TestErrorHandler_405MethodNotAllowed(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Request with wrong method
	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, fiber.StatusMethodNotAllowed, resp.StatusCode)

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var errResp ErrorResponse
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err)

	// Verify JSON structure
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "METHOD_NOT_ALLOWED", errResp.Code)
	assert.NotEmpty(t, errResp.Message)
	assert.Nil(t, errResp.Details)
}

func TestErrorHandler_DomainAppError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/test-not-found", func(c *fiber.Ctx) error {
		return domain.NewNotFoundError("User not found", "USER_NOT_FOUND")
	})

	app.Get("/test-bad-request", func(c *fiber.Ctx) error {
		details := map[string]string{"field": "email", "issue": "invalid format"}
		return domain.NewBadRequestError("Invalid input", "INVALID_INPUT", details)
	})

	t.Run("NotFoundError", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test-not-found", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errResp ErrorResponse
		err = json.Unmarshal(body, &errResp)
		require.NoError(t, err)

		assert.Equal(t, "error", errResp.Status)
		assert.Equal(t, "USER_NOT_FOUND", errResp.Code)
		assert.Equal(t, "User not found", errResp.Message)
		assert.Nil(t, errResp.Details)
	})

	t.Run("BadRequestErrorWithDetails", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test-bad-request", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errResp ErrorResponse
		err = json.Unmarshal(body, &errResp)
		require.NoError(t, err)

		assert.Equal(t, "error", errResp.Status)
		assert.Equal(t, "INVALID_INPUT", errResp.Code)
		assert.Equal(t, "Invalid input", errResp.Message)
		assert.NotNil(t, errResp.Details)

		// Verify details structure
		detailsMap, ok := errResp.Details.(map[string]any)
		require.True(t, ok, "Details should be a map")
		assert.Equal(t, "email", detailsMap["field"])
		assert.Equal(t, "invalid format", detailsMap["issue"])
	})
}

func TestErrorHandler_InternalServerError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		// Simulate an unexpected error
		return fiber.NewError(fiber.StatusInternalServerError, "Database connection failed")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var errResp ErrorResponse
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err)

	// Verify JSON structure
	assert.Equal(t, "error", errResp.Status)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", errResp.Code)
	assert.NotEmpty(t, errResp.Message)
	// In production, the message should be generic, but for this test we accept any message
}

func TestErrorHandler_JSONStructureCompliance(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return domain.NewBadRequestError("Test error", "TEST_ERROR", nil)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Verify that the response is valid JSON
	var jsonData map[string]any
	err = json.Unmarshal(body, &jsonData)
	require.NoError(t, err)

	// Verify all required fields are present
	assert.Contains(t, jsonData, "status")
	assert.Contains(t, jsonData, "message")
	assert.Contains(t, jsonData, "code")
	assert.Contains(t, jsonData, "details")

	// Verify field types
	assert.IsType(t, "", jsonData["status"])
	assert.IsType(t, "", jsonData["message"])
	assert.IsType(t, "", jsonData["code"])
	// details can be nil or any type

	// Verify status value
	assert.Equal(t, "error", jsonData["status"])
}

func TestErrorHandler_ProductionMasking(t *testing.T) {
	// Set environment to production
	os.Setenv("APP_ENV", "production")
	defer os.Setenv("APP_ENV", "development")

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/test-internal-error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "Sensitive database error")
	})

	req := httptest.NewRequest("GET", "/test-internal-error", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var errResp ErrorResponse
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err)

	// In production, the message MUST be generic
	assert.Equal(t, "Internal server error", errResp.Message)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", errResp.Code)
}
