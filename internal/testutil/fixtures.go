package testutil

import (
	"net/http"
	"os"
	"testing"
)

// FixtureOption is a function type that allows customization of HTTP responses
// when loading JSON fixtures. It provides a hook to modify response headers,
// status codes, or other aspects of the http.ResponseWriter before the fixture
// content is written.
type FixtureOption func(w http.ResponseWriter, r *http.Request)

func LoadFixture(t *testing.T, path string) []byte {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load fixture: %s (%v)", path, err)
	}
	return b
}

func WithHeader(k, v string) FixtureOption {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(k, v)
	}
}

func WithStatus(code int) FixtureOption {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}
}

func JSONFixtureHandler(t *testing.T, path string, opts ...FixtureOption) http.HandlerFunc {
	b := LoadFixture(t, path)

	return func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Type first before any WriteHeader calls
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}

		// Apply options (headers, status code, etc.)
		for _, opt := range opts {
			opt(w, r)
		}

		// Write response body
		_, _ = w.Write(b)
	}
}

// UnauthorizedHandler returns an HTTP handler that responds with a 401 Unauthorized
// error using the unauthorized_error.json fixture.
func UnauthorizedHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return JSONFixtureHandler(
		t,
		"testdata/unauthorized_error.json",
		WithStatus(http.StatusUnauthorized),
	)
}

// InternalServerErrorHandler returns an HTTP handler that responds with a 500 Internal Server Error.
func InternalServerErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// NotFoundHandler returns an HTTP handler that responds with a 404 Not Found
// error using the not_found_error.json fixture.
func NotFoundHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return JSONFixtureHandler(
		t,
		"testdata/not_found_error.json",
		WithStatus(http.StatusNotFound),
	)
}

// PaymentRequiredHandler returns an HTTP handler that responds with a 402 Payment Required
// error using the payment_required_error.json fixture.
func PaymentRequiredHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return JSONFixtureHandler(
		t,
		"testdata/payment_required_error.json",
		WithStatus(http.StatusPaymentRequired),
	)
}
