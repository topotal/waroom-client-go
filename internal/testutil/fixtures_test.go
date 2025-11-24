package testutil

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFixtureHandler_ContentType(t *testing.T) {
	// Create a temporary fixture file for testing
	tmpFile, err := os.CreateTemp("", "test-fixture-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`{"test": "data"}`)
	require.NoError(t, err)
	tmpFile.Close()

	t.Run("default content-type is application/json", func(t *testing.T) {
		handler := JSONFixtureHandler(t, tmpFile.Name())

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("WithStatus preserves default content-type", func(t *testing.T) {
		handler := JSONFixtureHandler(t, tmpFile.Name(), WithStatus(http.StatusCreated))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("WithHeader can override content-type", func(t *testing.T) {
		handler := JSONFixtureHandler(t, tmpFile.Name(), WithHeader("Content-Type", "text/plain"))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("WithHeader after WithStatus - header wins", func(t *testing.T) {
		handler := JSONFixtureHandler(
			t,
			tmpFile.Name(),
			WithStatus(http.StatusAccepted),
			WithHeader("Content-Type", "text/xml"),
		)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, "text/xml", rec.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("WithHeader before WithStatus - header wins", func(t *testing.T) {
		handler := JSONFixtureHandler(
			t,
			tmpFile.Name(),
			WithHeader("Content-Type", "application/xml"),
			WithStatus(http.StatusNoContent),
		)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, "application/xml", rec.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("multiple WithHeader calls - last one wins", func(t *testing.T) {
		handler := JSONFixtureHandler(
			t,
			tmpFile.Name(),
			WithHeader("Content-Type", "text/plain"),
			WithHeader("Content-Type", "text/html"),
			WithHeader("Content-Type", "application/xml"),
		)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, "application/xml", rec.Header().Get("Content-Type"))
	})

	t.Run("WithHeader with custom header plus WithStatus", func(t *testing.T) {
		handler := JSONFixtureHandler(
			t,
			tmpFile.Name(),
			WithHeader("X-Custom-Header", "custom-value"),
			WithStatus(http.StatusCreated),
		)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Equal(t, "custom-value", rec.Header().Get("X-Custom-Header"))
		assert.Equal(t, http.StatusCreated, rec.Code)

		body, err := io.ReadAll(rec.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"test": "data"}`, string(body))
	})
}
