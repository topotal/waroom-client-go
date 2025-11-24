package waroom

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestServiceService_ListServices(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/services", testutil.JSONFixtureHandler(t, "testdata/services_list_response.json"))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewServiceService(client)
	services, err := svc.ListServices(ctx, ListServicesParams{})

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
}

func TestServiceService_ListServices_Errors(t *testing.T) {
	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErrMsg string
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), "unauthorized"},
		{"server error", testutil.InternalServerErrorHandler(), ""},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/services", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			svc := NewServiceService(client)
			services, err := svc.ListServices(ctx, ListServicesParams{})

			require.Error(t, err)
			assert.Nil(t, services)
			if tt.wantErrMsg != "" {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}
