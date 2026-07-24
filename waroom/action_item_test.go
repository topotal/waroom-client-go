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

func TestActionItemService_List_Errors(t *testing.T) {
	incidentUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/550e8400-e29b-41d4-a716-446655440001/action_items", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	actionItems, err := client.ActionItems.List(ctx, ListIncidentActionItemsParams{UUID: incidentUUID})

	require.Error(t, err)
	assert.Nil(t, actionItems)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: list incident action items:")
}

func TestIncidentService_ListActionItems_Errors(t *testing.T) {
	incidentUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/550e8400-e29b-41d4-a716-446655440001/action_items", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	actionItems, err := client.Incidents.ListActionItems(ctx, ListIncidentActionItemsParams{UUID: incidentUUID})

	require.Error(t, err)
	assert.Nil(t, actionItems)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: list incident action items:")
}
