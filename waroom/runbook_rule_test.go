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

func TestRunbookRuleService_List_Errors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/svc/runbook_rules", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	rules, err := client.RunbookRules.List(ctx, ListServiceRunbookRulesParams{Name: "svc"})

	require.Error(t, err)
	assert.Nil(t, rules)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: list service runbook rules:")
}
