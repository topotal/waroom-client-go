package waroom

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestUserService_List_Errors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/users", testutil.UnauthorizedHandler(t))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	users, err := client.Users.List(ctx, ListUsersParams{})

	require.Error(t, err)
	assert.Nil(t, users)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}, apiErr)
	assert.Contains(t, err.Error(), "waroom: list users:")
}

func TestUserService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/users/{sub}", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v0/users/auth0%7Cabc123", r.URL.EscapedPath())
		testutil.JSONFixtureHandler(t, "testdata/user_get_response.json")(w, r)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	user, err := client.Users.Get(ctx, GetUserParams{Sub: "auth0|abc123"})

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "auth0|abc123", user.Sub)
	assert.Equal(t, "alice", user.Nickname)
	assert.Equal(t, "https://example.com/alice.png", user.Picture)
	assert.Equal(t, "active", string(user.Status))
	require.False(t, user.SlackUser.Null)
	assert.Equal(t, "U0123456789", user.SlackUser.Value.Key)
	assert.True(t, user.PhoneVerified)
	assert.Equal(t, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), user.CreatedAt)
}

func TestUserService_Get_Errors(t *testing.T) {
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
			name:       "server error",
			handler:    testutil.InternalServerErrorHandler(),
			wantAPIErr: nil,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/users/{sub}", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			user, err := client.Users.Get(ctx, GetUserParams{Sub: "auth0|abc123"})

			require.Error(t, err)
			assert.Nil(t, user)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: get user:")
		})
	}
}

func TestUserService_Invite(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/users/invitations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body struct {
			Invitation struct {
				Emails []string `json:"emails"`
			} `json:"invitation"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []string{"bob@example.com", "carol@example.com"}, body.Invitation.Emails)

		testutil.JSONFixtureHandler(t, "testdata/user_invite_response.json", testutil.WithStatus(http.StatusAccepted))(w, r)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	users, err := client.Users.Invite(ctx, InviteUsersReq{
		Emails: []string{"bob@example.com", "carol@example.com"},
	})

	require.NoError(t, err)
	require.Len(t, users, 2)

	assert.Equal(t, "auth0|invited001", users[0].Sub)
	assert.Equal(t, "bob", users[0].Nickname)
	assert.Equal(t, "invitation_pending", string(users[0].Status))
	assert.True(t, users[0].SlackUser.Null)
	assert.False(t, users[0].PhoneVerified)

	assert.Equal(t, "auth0|invited002", users[1].Sub)
	assert.Equal(t, "carol", users[1].Nickname)
	assert.Equal(t, "invitation_pending", string(users[1].Status))
}

func TestUserService_Invite_Errors(t *testing.T) {
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
			name:       "conflict error",
			handler:    testutil.JSONFixtureHandler(t, "testdata/user_invite_conflict_error.json", testutil.WithStatus(http.StatusConflict)),
			wantAPIErr: &APIError{StatusCode: http.StatusConflict, Messages: []string{"Nickname bob has already been taken"}},
		},
		{
			name:       "unprocessable entity error",
			handler:    testutil.JSONFixtureHandler(t, "testdata/user_invite_unprocessable_entity_error.json", testutil.WithStatus(http.StatusUnprocessableEntity)),
			wantAPIErr: &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: []string{"Emails must be 50 or fewer"}},
		},
		{
			name: "too many requests error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
			},
			wantAPIErr: &APIError{StatusCode: http.StatusTooManyRequests},
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
			mux.HandleFunc("/api/v0/users/invitations", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			users, err := client.Users.Invite(ctx, InviteUsersReq{
				Emails: []string{"bob@example.com"},
			})

			require.Error(t, err)
			assert.Nil(t, users)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: invite users:")
		})
	}
}
