package waroom

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestServiceService_ListServices(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services", testutil.JSONFixtureHandler(t, "testdata/services_list_response.json"))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewServiceService(client)
	services, err := svc.ListServices(ctx, ListServicesParams{})

	require.NoError(t, err)
	assert.Len(t, services, 2)

	// 新specではserviceはnameのみを返す(labelsフィールドは無くなった)
	assert.Equal(t, "Test Service 1", services[0].Name)
	assert.Equal(t, "Test Service 2", services[1].Name)
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
			mux.HandleFunc("/api/v0/services", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
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
