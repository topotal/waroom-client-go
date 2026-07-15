package waroom

import (
	"context"
	"errors"
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
		wantAPIErr *APIError // nil expects a transport-layer error (not an APIError)
	}{
		{
			name:       "unauthorized error",
			handler:    testutil.UnauthorizedHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}},
		},
		{
			name:       "payment required error",
			handler:    testutil.PaymentRequiredHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{"Payment Required"}},
		},
		{
			name:       "server error",
			handler:    testutil.InternalServerErrorHandler(),
			wantAPIErr: nil,
		},
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

			var apiErr *APIError
			if tt.wantAPIErr != nil {
				require.True(t, errors.As(err, &apiErr))
				assert.Equal(t, tt.wantAPIErr, apiErr)
				assert.Contains(t, err.Error(), "waroom: list services:")
			} else {
				assert.False(t, errors.As(err, &apiErr))
			}
		})
	}
}
