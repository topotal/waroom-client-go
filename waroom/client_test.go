package waroom

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestNewClient_DefaultBaseURL(t *testing.T) {
	assert.Equal(t, "https://api.app.waroom.com", DefaultBaseURL)
}

func TestWithRequestDump_DoesNotLeakAPIKey(t *testing.T) {
	const apiKey = "test-secret-api-key"

	var gotAuth string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/users/{sub}", func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		testutil.JSONFixtureHandler(t, "testdata/user_get_response.json")(w, r)
	})

	var dump bytes.Buffer
	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		return NewClient(context.Background(),
			WithAPIKey(apiKey),
			WithBaseURL(baseURL),
			WithRequestDump(),
			func(c *config) { c.dumpWriter = &dump },
		)
	})
	defer cleanup()

	_, err := client.Users.Get(context.Background(), GetUserParams{Sub: "auth0|abc123"})
	require.NoError(t, err)

	assert.Equal(t, "Bearer "+apiKey, gotAuth)
	assert.Contains(t, dump.String(), "Authorization: Bearer [REDACTED]")
	assert.NotContains(t, dump.String(), apiKey)
}
