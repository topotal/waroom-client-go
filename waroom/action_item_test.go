package waroom

import (
	"context"
	"encoding/json"
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

func TestActionItemService_Update(t *testing.T) {
	incidentUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	actionItemUUID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	fixture := testutil.JSONFixtureHandler(t, "testdata/action_item_update_response.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/550e8400-e29b-41d4-a716-446655440001/action_items/660e8400-e29b-41d4-a716-446655440002", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)

		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"action_item": {
				"title":             "Rotate the leaked credentials",
				"status":            "closed",
				"assignee_nickname": "alice",
			},
		}, body)

		fixture(w, r)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	req := UpdateActionItemReq{
		Title:            NewOptString("Rotate the leaked credentials"),
		AssigneeNickname: NewOptNilString("alice"),
	}
	req.Status.SetTo("closed")

	ctx := context.Background()
	actionItem, err := client.ActionItems.Update(ctx, req, UpdateIncidentActionItemParams{
		UUID:           incidentUUID,
		ActionItemUUID: actionItemUUID,
	})

	require.NoError(t, err)
	require.NotNil(t, actionItem)
	assert.Equal(t, actionItemUUID, actionItem.UUID)
	assert.Equal(t, "Rotate the leaked credentials", actionItem.Title)
	assert.Equal(t, "closed", string(actionItem.Status))
	assert.Equal(t, 1, actionItem.Position)
	assignee, ok := actionItem.Assignee.Get()
	require.True(t, ok)
	assert.Equal(t, "alice", assignee.Nickname)
}

func TestActionItemService_Update_Unassign(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/550e8400-e29b-41d4-a716-446655440001/action_items/660e8400-e29b-41d4-a716-446655440002", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"action_item": {"assignee_nickname": nil},
		}, body)

		testutil.JSONFixtureHandler(t, "testdata/action_item_update_response.json")(w, r)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	var req UpdateActionItemReq
	req.AssigneeNickname.SetToNull()

	ctx := context.Background()
	_, err := client.ActionItems.Update(ctx, req, UpdateIncidentActionItemParams{
		UUID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		ActionItemUUID: uuid.MustParse("660e8400-e29b-41d4-a716-446655440002"),
	})

	require.NoError(t, err)
}

func TestActionItemService_Update_Errors(t *testing.T) {
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
			mux.HandleFunc("/api/v0/incidents/550e8400-e29b-41d4-a716-446655440001/action_items/660e8400-e29b-41d4-a716-446655440002", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			actionItem, err := client.ActionItems.Update(ctx, UpdateActionItemReq{
				Title: NewOptString("Rotate the leaked credentials"),
			}, UpdateIncidentActionItemParams{
				UUID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
				ActionItemUUID: uuid.MustParse("660e8400-e29b-41d4-a716-446655440002"),
			})

			require.Error(t, err)
			assert.Nil(t, actionItem)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: update incident action item:")
		})
	}
}
