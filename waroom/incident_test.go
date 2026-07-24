package waroom

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gen "github.com/topotal/waroom-client-go/internal/ogen"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestIncidentService_List_Errors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	incidents, err := client.Incidents.List(ctx, ListIncidentsParams{})

	require.Error(t, err)
	assert.Nil(t, incidents)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: list incidents:")
}

func TestIncidentService_Get_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/"+testUUID.String(), testutil.NotFoundHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	incident, err := client.Incidents.Get(ctx, testUUID)

	require.Error(t, err)
	assert.Nil(t, incident)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: get incident:")
}

func TestIncidentService_Create_Errors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	// Severity must be a valid enum value or client-side request validation
	// rejects the request before it ever reaches the (mocked) server.
	incident, err := client.Incidents.Create(ctx, CreateIncidentReq{Severity: gen.CreateIncidentReqIncidentSeverityLow})

	require.Error(t, err)
	assert.Nil(t, incident)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: create incident:")
}
