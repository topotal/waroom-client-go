package waroom

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestServiceService_List(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services", testutil.JSONFixtureHandler(t, "testdata/services_list_response.json"))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	services, err := client.Services.List(ctx, ListServicesParams{})

	require.NoError(t, err)
	assert.Len(t, services, 2)

	assert.Equal(t, "Test Service 1", services[0].Name)
	assert.Equal(t, "Test Service 2", services[1].Name)
}

func TestServiceService_List_Errors(t *testing.T) {
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
			services, err := client.Services.List(ctx, ListServicesParams{})

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

func TestServiceService_Create(t *testing.T) {
	fixture := testutil.JSONFixtureHandler(t, "testdata/service_create_response.json", testutil.WithStatus(http.StatusCreated))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body struct {
			Service struct {
				Name string `json:"name"`
			} `json:"service"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "test-service", body.Service.Name)

		fixture(w, r)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	service, err := client.Services.Create(ctx, CreateServiceReq{Name: "test-service"})

	require.NoError(t, err)
	require.NotNil(t, service)
	assert.Equal(t, "test-service", service.Name)
}

func TestServiceService_Create_Errors(t *testing.T) {
	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantAPIErr *APIError
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
			name:       "unprocessable entity error",
			handler:    testutil.UnprocessableEntityHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: []string{"Name has already been taken"}},
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
			service, err := client.Services.Create(ctx, CreateServiceReq{Name: "test-service"})

			require.Error(t, err)
			assert.Nil(t, service)

			var apiErr *APIError
			if tt.wantAPIErr != nil {
				require.True(t, errors.As(err, &apiErr))
				assert.Equal(t, tt.wantAPIErr, apiErr)
				assert.Contains(t, err.Error(), "waroom: create service:")
			} else {
				assert.False(t, errors.As(err, &apiErr))
			}
		})
	}
}

func TestServiceService_GetArchitectureContext_Errors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/svc/service_architecture_context", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	result, err := client.Services.GetArchitectureContext(ctx, GetServiceArchitectureContextParams{Name: "svc"})

	require.Error(t, err)
	assert.Nil(t, result)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: get service architecture context:")
}
