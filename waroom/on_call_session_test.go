package waroom

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

func TestOnCallSessionService_List(t *testing.T) {
	mux := http.NewServeMux()
	fixture := testutil.JSONFixtureHandler(t, "testdata/on_call_session_list_response.json")
	mux.HandleFunc("/api/v0/on_call/sessions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		q := r.URL.Query()
		assert.Equal(t, "acked", q.Get("status"))
		assert.Equal(t, "6f1c2b3a-0000-4000-8000-0000000000a1", q.Get("incident_uuid"))
		assert.Equal(t, "2", q.Get("page"))
		assert.Equal(t, "50", q.Get("per_page"))
		fixture(w, r)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	sessions, err := client.OnCallSessions.List(ctx, ListOnCallSessionsParams{
		Status:       NewOptOnCallSessionStatusFilter("acked"),
		IncidentUUID: NewOptString("6f1c2b3a-0000-4000-8000-0000000000a1"),
		Page:         NewOptInt(2),
		PerPage:      NewOptInt(50),
	})

	require.NoError(t, err)
	require.Len(t, sessions, 2)

	first := sessions[0]
	assert.Equal(t, "6f1c2b3a-0000-4000-8000-000000000001", first.UUID)
	assert.Equal(t, "acked", string(first.Status))
	assert.Equal(t, "integration", string(first.Source))
	assert.Equal(t, 1, first.CurrentStepPosition)
	assert.Equal(t, "6f1c2b3a-0000-4000-8000-0000000000a1", first.Incident.UUID)
	policy, ok := first.EscalationPolicy.Get()
	require.True(t, ok)
	assert.Equal(t, "Primary escalation", policy.Name)
	assert.True(t, first.TargetUser.IsNull())
	trigger, ok := first.Trigger.Get()
	require.True(t, ok)
	assert.Equal(t, "Datadog trigger", trigger.Name)
	assert.True(t, first.CreatedBy.IsNull())
	ackedBy, ok := first.AckedBy.Get()
	require.True(t, ok)
	assert.Equal(t, "alice", ackedBy.Nickname)
	ackedAt, ok := first.AckedAt.Get()
	require.True(t, ok)
	assert.Equal(t, time.Date(2026, 11, 20, 10, 5, 0, 0, time.UTC), ackedAt)

	second := sessions[1]
	assert.Equal(t, "running", string(second.Status))
	assert.True(t, second.EscalationPolicy.IsNull())
	targetUser, ok := second.TargetUser.Get()
	require.True(t, ok)
	assert.Equal(t, "bob", targetUser.Nickname)
	createdBy, ok := second.CreatedBy.Get()
	require.True(t, ok)
	assert.Equal(t, "auth0|user-1", createdBy.Sub)
	assert.True(t, second.AckedBy.IsNull())
	assert.True(t, second.AckedAt.IsNull())
}

func TestOnCallSessionService_List_Errors(t *testing.T) {
	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantAPIErr *APIError
	}{
		{
			name: "bad request error",
			handler: testutil.JSONFixtureHandler(t, "testdata/on_call_session_list_bad_request_error.json",
				testutil.WithStatus(http.StatusBadRequest)),
			wantAPIErr: &APIError{StatusCode: http.StatusBadRequest, Messages: []string{"status is invalid"}},
		},
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
			name:       "server error",
			handler:    testutil.InternalServerErrorHandler(),
			wantAPIErr: nil,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/sessions", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			sessions, err := client.OnCallSessions.List(ctx, ListOnCallSessionsParams{})

			require.Error(t, err)
			assert.Nil(t, sessions)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: list on-call sessions:")
		})
	}
}

func TestOnCallSessionService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/sessions/6f1c2b3a-0000-4000-8000-000000000001", testutil.JSONFixtureHandler(t, "testdata/on_call_session_get_response.json"))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	session, err := client.OnCallSessions.Get(ctx, GetOnCallSessionParams{
		UUID: "6f1c2b3a-0000-4000-8000-000000000001",
	})

	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "6f1c2b3a-0000-4000-8000-000000000001", session.UUID)
	assert.Equal(t, "acked", string(session.Status))
	assert.Equal(t, "6f1c2b3a-0000-4000-8000-0000000000a1", session.Incident.UUID)
	assert.Equal(t, time.Date(2026, 11, 20, 10, 0, 0, 0, time.UTC), session.CreatedAt)

	require.Len(t, session.Activities, 3)
	var activities []OnCallSessionActivity = session.Activities

	assert.Equal(t, "session_created", string(activities[0].ActivityType))
	assert.True(t, activities[0].User.IsNull())
	assert.True(t, activities[0].Team.IsNull())

	assert.Equal(t, "user_skipped", string(activities[1].ActivityType))
	skipReason, ok := activities[1].SkipReason.Get()
	require.True(t, ok)
	assert.Equal(t, "phone_number_unverified", string(skipReason))
	team, ok := activities[1].Team.Get()
	require.True(t, ok)
	assert.Equal(t, "SRE", team.Name)

	assert.Equal(t, "dtmf_received", string(activities[2].ActivityType))
	dtmf, ok := activities[2].DtmfAction.Get()
	require.True(t, ok)
	assert.Equal(t, "acknowledged", string(dtmf))
	user, ok := activities[2].User.Get()
	require.True(t, ok)
	assert.Equal(t, "alice", user.Nickname)
}

func TestOnCallSessionService_Get_Errors(t *testing.T) {
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
			mux.HandleFunc("/api/v0/on_call/sessions/6f1c2b3a-0000-4000-8000-000000000001", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			session, err := client.OnCallSessions.Get(ctx, GetOnCallSessionParams{
				UUID: "6f1c2b3a-0000-4000-8000-000000000001",
			})

			require.Error(t, err)
			assert.Nil(t, session)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: get on-call session:")
		})
	}
}
