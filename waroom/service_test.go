package waroom

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func setupTestClient(t *testing.T, mux *http.ServeMux) (*Client, func()) {
	t.Helper()
	return testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		return NewClient(context.Background(), WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
}

func TestServiceService_ListServices(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services", testutil.JSONFixtureHandler(t, "testdata/services_list_response.json"))

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	list, err := client.Services().ListServices(context.Background())

	require.NoError(t, err)
	assert.Len(t, list.Items, 2)

	// 新specではserviceはnameのみを返す(labelsフィールドは無くなった)
	assert.Equal(t, "Test Service 1", list.Items[0].Name)
	assert.Equal(t, "Test Service 2", list.Items[1].Name)

	assert.Equal(t, 1, list.CurrentPage)
	assert.Nil(t, list.NextPage)
	assert.Equal(t, 2, list.TotalCount)
}

func TestServiceService_ListServices_Errors(t *testing.T) {
	errorTests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr error
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), ErrUnauthorized},
		{"server error", testutil.InternalServerErrorHandler(), nil},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/services", tt.handler)

			client, cleanup := setupTestClient(t, mux)
			defer cleanup()

			list, err := client.Services().ListServices(context.Background())

			require.Error(t, err)
			assert.Nil(t, list)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestAPIError_Fields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services", testutil.UnauthorizedHandler(t))

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	_, err := client.Services().ListServices(context.Background())

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 401, apiErr.StatusCode)
	assert.Equal(t, "listServices", apiErr.OperationID)
	assert.Equal(t, []string{"Not Authenticated"}, apiErr.Messages)
	assert.Contains(t, apiErr.Error(), "401 Unauthorized")
	assert.NotErrorIs(t, err, ErrNotFound)
}
