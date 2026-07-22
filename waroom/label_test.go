package waroom

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestLabelService_List(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels", testutil.JSONFixtureHandler(t, "testdata/labels_list_response.json"))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	labels, err := client.Labels.List(ctx, ListServiceLabelsParams{
		Name: "test-service",
	})

	require.NoError(t, err)
	assert.Len(t, labels, 2)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), labels[0].UUID)
	assert.Equal(t, "production", labels[0].Name)
	assert.Equal(t, "FF0000", labels[0].Color)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"), labels[1].UUID)
	assert.Equal(t, "staging", labels[1].Name)
	assert.Equal(t, "00FF00", labels[1].Color)
}

func TestLabelService_List_Errors(t *testing.T) {
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
			name:       "not found error",
			handler:    testutil.NotFoundHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}},
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
			mux.HandleFunc("/api/v0/services/test-service/labels", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			labels, err := client.Labels.List(ctx, ListServiceLabelsParams{
				Name: "test-service",
			})

			require.Error(t, err)
			assert.Nil(t, labels)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: list service labels:")
		})
	}
}

func TestLabelService_Create(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels", testutil.JSONFixtureHandler(
		t,
		"testdata/label_create_response.json",
		testutil.WithStatus(http.StatusCreated),
	))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	label, err := client.Labels.Create(ctx, CreateServiceLabelReq{
		Name:  "production",
		Color: "FF0000",
	}, CreateServiceLabelParams{
		Name: "test-service",
	})

	require.NoError(t, err)
	require.NotNil(t, label)
	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), label.UUID)
	assert.Equal(t, "production", label.Name)
	assert.Equal(t, "FF0000", label.Color)
}

func TestLabelService_Create_Errors(t *testing.T) {
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
			mux.HandleFunc("/api/v0/services/test-service/labels", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			label, err := client.Labels.Create(ctx, CreateServiceLabelReq{
				Name:  "production",
				Color: "FF0000",
			}, CreateServiceLabelParams{
				Name: "test-service",
			})

			require.Error(t, err)
			assert.Nil(t, label)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: create service label:")
		})
	}
}

func TestLabelService_Update(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", testutil.JSONFixtureHandler(
		t,
		"testdata/label_update_response.json",
	))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	label, err := client.Labels.Update(ctx, UpdateServiceLabelReq{
		Name:  NewOptString("production-updated"),
		Color: NewOptString("0000FF"),
	}, UpdateServiceLabelParams{
		Name: "test-service",
		UUID: testUUID,
	})

	require.NoError(t, err)
	require.NotNil(t, label)
	assert.Equal(t, testUUID, label.UUID)
	assert.Equal(t, "production-updated", label.Name)
	assert.Equal(t, "0000FF", label.Color)
}

func TestLabelService_Update_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

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
			name:       "not found error",
			handler:    testutil.NotFoundHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}},
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
			mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			label, err := client.Labels.Update(ctx, UpdateServiceLabelReq{
				Name: NewOptString("production-updated"),
			}, UpdateServiceLabelParams{
				Name: "test-service",
				UUID: testUUID,
			})

			require.Error(t, err)
			assert.Nil(t, label)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: update service label:")
		})
	}
}

func TestLabelService_Delete(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	err := client.Labels.Delete(ctx, DeleteServiceLabelParams{
		Name: "test-service",
		UUID: testUUID,
	})

	require.NoError(t, err)
}

func TestLabelService_Delete_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

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
			name:       "not found error",
			handler:    testutil.NotFoundHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}},
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
			mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			err := client.Labels.Delete(ctx, DeleteServiceLabelParams{
				Name: "test-service",
				UUID: testUUID,
			})

			require.Error(t, err)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: delete service label:")
		})
	}
}

func TestLabelService_UpdateIncidentLabels_Errors(t *testing.T) {
	incidentUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

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
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/incidents/550e8400-e29b-41d4-a716-446655440001/labels", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			labels, err := client.Labels.UpdateIncidentLabels(ctx, UpdateIncidentLabelsReq{}, UpdateIncidentLabelsParams{UUID: incidentUUID})

			require.Error(t, err)
			assert.Nil(t, labels)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: update incident labels:")
		})
	}
}

// assertAPIError asserts that err unwraps to the expected *APIError and that
// its message carries the operation prefix. A nil want asserts that err does
// not unwrap to *APIError (i.e. it is a transport-layer error).
func assertAPIError(t *testing.T, err error, want *APIError, wantPrefix string) {
	t.Helper()

	var apiErr *APIError
	if want == nil {
		assert.False(t, errors.As(err, &apiErr))
		return
	}
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, want, apiErr)
	assert.Contains(t, err.Error(), wantPrefix)
}
