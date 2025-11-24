package waroom

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestServiceService_ListServices(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock server
		mux := http.NewServeMux()
		mux.HandleFunc("/services", testutil.JSONFixtureHandler(t, "testdata/services_list_response.json"))

		server := httptest.NewServer(mux)
		defer server.Close()

		// Create client with test server URL
		ctx := context.Background()
		client, err := NewClient(ctx, WithBaseURL(server.URL))
		require.NoError(t, err)

		// Execute test
		svc := NewServiceService(client)
		services, err := svc.ListServices(ctx, ListServicesParams{})

		// Assertions
		require.NoError(t, err)
		assert.Len(t, services, 2)

		assert.Equal(t, "Test Service 1", services[0].Name)
		assert.Len(t, services[0].Labels, 1)
		assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), services[0].Labels[0].UUID)
		assert.Equal(t, "production", services[0].Labels[0].Name)
		assert.Equal(t, "FF0000", services[0].Labels[0].Color.Value)

		assert.Equal(t, "Test Service 2", services[1].Name)
		assert.Len(t, services[1].Labels, 2)
		assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"), services[1].Labels[0].UUID)
		assert.Equal(t, "staging", services[1].Labels[0].Name)
		assert.Equal(t, "00FF00", services[1].Labels[0].Color.Value)
		assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"), services[1].Labels[1].UUID)
		assert.Equal(t, "backend", services[1].Labels[1].Name)
		assert.Equal(t, "0000FF", services[1].Labels[1].Color.Value)
	})

	t.Run("unauthorized error", func(t *testing.T) {
		// Setup mock server with error response
		mux := http.NewServeMux()
		mux.HandleFunc("/services", testutil.JSONFixtureHandler(
			t,
			"testdata/unauthorized_error.json",
			testutil.WithStatus(http.StatusUnauthorized),
		))

		server := httptest.NewServer(mux)
		defer server.Close()

		// Create client with test server URL
		ctx := context.Background()
		client, err := NewClient(ctx, WithBaseURL(server.URL))
		require.NoError(t, err)

		// Execute test
		svc := NewServiceService(client)
		services, err := svc.ListServices(ctx, ListServicesParams{})

		// Assertions
		require.Error(t, err)
		assert.Nil(t, services)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("server error", func(t *testing.T) {
		// Setup mock server with 500 error
		mux := http.NewServeMux()
		mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		// Create client with test server URL
		ctx := context.Background()
		client, err := NewClient(ctx, WithBaseURL(server.URL))
		require.NoError(t, err)

		// Execute test
		svc := NewServiceService(client)
		services, err := svc.ListServices(ctx, ListServicesParams{})

		// Assertions
		require.Error(t, err)
		assert.Nil(t, services)
	})
}
