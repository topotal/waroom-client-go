package waroom

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestIncidentService_GetIncident_NotFound(t *testing.T) {
	testUUID := uuid.MustParse("d32fc8dd-e424-4974-b27e-c55bb323a6cb")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/d32fc8dd-e424-4974-b27e-c55bb323a6cb", testutil.NotFoundHandler(t))

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	inc, err := client.Incidents().GetIncident(context.Background(), testUUID)

	require.Error(t, err)
	assert.Nil(t, inc)
	assert.ErrorIs(t, err, ErrNotFound)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Equal(t, "getIncident", apiErr.OperationID)
	assert.Equal(t, []string{"Not found"}, apiErr.Messages)
}
