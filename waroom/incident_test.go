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

func TestIncidentService_ListIncidents_Errors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewIncidentService(client)
	incidents, err := svc.ListIncidents(ctx, ListIncidentsParams{})

	require.Error(t, err)
	assert.Nil(t, incidents)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: list incidents:")
}

func TestIncidentService_GetIncident_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/"+testUUID.String(), testutil.NotFoundHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewIncidentService(client)
	incident, err := svc.GetIncident(ctx, testUUID)

	require.Error(t, err)
	assert.Nil(t, incident)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: get incident:")
}
