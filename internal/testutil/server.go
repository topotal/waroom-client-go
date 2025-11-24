package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// SetupTestServer creates a test HTTP server with the provided mux and returns
// the server URL along with a cleanup function. This is useful for creating
// test clients without introducing circular dependencies.
//
// Example usage:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/services", testutil.JSONFixtureHandler(t, "testdata/services.json"))
//	serverURL, cleanup := testutil.SetupTestServer(t, mux)
//	defer cleanup()
//
//	client, err := NewClient(ctx, WithBaseURL(serverURL))
//	require.NoError(t, err)
func SetupTestServer(t *testing.T, mux *http.ServeMux) (string, func()) {
	t.Helper()

	server := httptest.NewServer(mux)

	cleanup := func() {
		server.Close()
	}

	return server.URL, cleanup
}

// SetupTestClient creates a test HTTP server with the provided mux and returns
// a client created by the provided factory function, along with a cleanup function.
// This avoids circular dependencies by taking a factory function that creates the client.
//
// Example usage:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/services", testutil.JSONFixtureHandler(t, "testdata/services.json"))
//	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*waroom.Client, error) {
//	    ctx := context.Background()
//	    return waroom.NewClient(ctx, waroom.WithBaseURL(baseURL))
//	})
//	defer cleanup()
func SetupTestClient[T any](t *testing.T, mux *http.ServeMux, factory func(baseURL string) (T, error)) (T, func()) {
	t.Helper()

	serverURL, cleanup := SetupTestServer(t, mux)

	client, err := factory(serverURL)
	if err != nil {
		cleanup()
		var zero T
		t.Fatalf("failed to create client: %v", err)
		return zero, func() {}
	}

	return client, cleanup
}
