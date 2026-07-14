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

func TestIncidentService_GetIncident(t *testing.T) {
	testUUID := uuid.MustParse("d32fc8dd-e424-4974-b27e-c55bb323a6cb")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents/d32fc8dd-e424-4974-b27e-c55bb323a6cb", testutil.JSONFixtureHandler(t, "testdata/incident_get_response.json"))

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	inc, err := client.Incidents().GetIncident(context.Background(), testUUID)

	require.NoError(t, err)
	require.NotNil(t, inc)
	assert.Equal(t, testUUID, inc.UUID)
	assert.Equal(t, "Database outage", inc.Title)
	assert.Equal(t, "SEV1", inc.Severity)
	assert.Equal(t, "Test Service 1", inc.Service.Name)

	ttd, ok := inc.Metrics.Ttd.Get()
	require.True(t, ok)
	assert.Equal(t, 60, ttd)
	_, ok = inc.Metrics.Tti.Get()
	assert.False(t, ok, "tti should be null")
}

func TestIncidentService_ListIncidents(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents", func(w http.ResponseWriter, r *http.Request) {
		// Options struct のフィールドがクエリパラメータへ転写されること
		q := r.URL.Query()
		assert.Equal(t, "2", q.Get("page"))
		assert.Equal(t, "10", q.Get("per_page"))
		assert.Equal(t, "asc", q.Get("sort_order"))
		testutil.JSONFixtureHandler(t, "testdata/incidents_list_response.json")(w, r)
	})

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	list, err := client.Incidents().ListIncidents(context.Background(), ListIncidentsOptions{
		Page:      NewOptInt(2),
		PerPage:   NewOptInt(10),
		SortOrder: NewOptListIncidentsSortOrder(ListIncidentsSortOrderAsc),
	})

	require.NoError(t, err)
	assert.Len(t, list.Items, 2)
	assert.Equal(t, "Database outage", list.Items[0].Title)
	assert.Equal(t, "API latency spike", list.Items[1].Title)
	assert.Equal(t, 1, list.CurrentPage)
	assert.Equal(t, 2, list.TotalCount)
}

func TestIncidentService_ListIncidents_PaymentRequired(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/incidents", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)
		// 402 のボディは required な `error: string`(401 の errors 配列、
		// 404/422 の optional string とは別の第 3 形)
		_, _ = w.Write([]byte(`{"error": "The V0 Public API is available only on the Team or Pro plan."}`))
	})

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	list, err := client.Incidents().ListIncidents(context.Background(), ListIncidentsOptions{})

	require.Error(t, err)
	assert.Nil(t, list)
	assert.ErrorIs(t, err, ErrPaymentRequired)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 402, apiErr.StatusCode)
	assert.Equal(t, []string{"The V0 Public API is available only on the Team or Pro plan."}, apiErr.Messages)
}

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
