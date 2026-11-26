package httpx

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithRequestDump_RedactsCredentials(t *testing.T) {
	const secret = "super-secret-api-key"
	const body = `{"name":"example"}`

	var gotHeader http.Header
	var gotBody string
	next := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		gotHeader = r.Header
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		gotBody = string(b)
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Request: r}, nil
	})

	req, err := http.NewRequest(http.MethodPost, "https://api.example.com/api/v0/services", strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Proxy-Authorization", "Basic "+secret)
	req.Header.Set("Cookie", "session="+secret)
	req.Header.Set("X-Api-Key", secret)

	var out bytes.Buffer
	_, err = WithRequestDump(&out)(next).RoundTrip(req)
	require.NoError(t, err)

	dump := out.String()
	assert.NotContains(t, dump, secret)
	assert.Contains(t, dump, "Authorization: Bearer [REDACTED]")
	assert.Contains(t, dump, "Proxy-Authorization: Basic [REDACTED]")
	assert.Contains(t, dump, "Cookie: [REDACTED]")
	assert.Contains(t, dump, "X-Api-Key: [REDACTED]")
	assert.Contains(t, dump, body)

	// The request sent upstream keeps its real credentials and body.
	assert.Equal(t, "Bearer "+secret, gotHeader.Get("Authorization"))
	assert.Equal(t, body, gotBody)
}
