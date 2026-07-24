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

func TestIncidentDescriptionTemplateService_Get_Errors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incident_description_template", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	template, err := client.IncidentDescriptionTemplate.Get(ctx)

	require.Error(t, err)
	assert.Nil(t, template)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: get incident description template:")
}
