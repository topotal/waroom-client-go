package waroom

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

const (
	testAbsenceSub  = "auth0|1234567890"
	testAbsenceUUID = "770e8400-e29b-41d4-a716-446655440001"
	// The sub contains "|", which must be percent-encoded in the path.
	testAbsencesPath = "/api/v0/on_call/users/auth0%7C1234567890/absences"
	testAbsencePath  = testAbsencesPath + "/" + testAbsenceUUID
)

func newOnCallUserAbsenceTestClient(t *testing.T, mux *http.ServeMux) (*Client, func()) {
	t.Helper()
	return testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		return NewClient(context.Background(), WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
}

// onCallUserAbsenceHandler checks the method and the escaped path before
// delegating to next.
func onCallUserAbsenceHandler(t *testing.T, method, escapedPath string, next http.HandlerFunc) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method)
		assert.Equal(t, escapedPath, r.URL.EscapedPath())
		assert.Equal(t, testAbsenceSub, r.PathValue("sub"))
		next(w, r)
	}
}

func onCallUserAbsenceBadRequestHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return testutil.JSONFixtureHandler(t, "testdata/on_call_user_absence_bad_request_error.json", testutil.WithStatus(http.StatusBadRequest))
}

func onCallUserAbsenceUnprocessableEntityHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return testutil.JSONFixtureHandler(t, "testdata/on_call_user_absence_unprocessable_entity_error.json", testutil.WithStatus(http.StatusUnprocessableEntity))
}

type onCallUserAbsenceErrorTest struct {
	name       string
	handler    http.HandlerFunc
	wantAPIErr *APIError
}

func onCallUserAbsenceCommonErrorTests(t *testing.T) []onCallUserAbsenceErrorTest {
	return []onCallUserAbsenceErrorTest{
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
}

func onCallUserAbsenceWriteErrorTests(t *testing.T) []onCallUserAbsenceErrorTest {
	return append(onCallUserAbsenceCommonErrorTests(t),
		onCallUserAbsenceErrorTest{
			name:       "bad request error",
			handler:    onCallUserAbsenceBadRequestHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusBadRequest, Messages: []string{"param is missing or the value is empty: absence"}},
		},
		onCallUserAbsenceErrorTest{
			name:       "unprocessable entity error",
			handler:    onCallUserAbsenceUnprocessableEntityHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: []string{"End at must be after start at"}},
		},
	)
}

func TestOnCallUserAbsenceService_List(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/users/{sub}/absences", onCallUserAbsenceHandler(t, http.MethodGet, testAbsencesPath, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "true", q.Get("include_expired"))
		assert.Equal(t, "2", q.Get("page"))
		assert.Equal(t, "50", q.Get("per_page"))
		testutil.JSONFixtureHandler(t, "testdata/on_call_user_absence_list_response.json")(w, r)
	}))

	client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
	defer cleanup()

	params := ListOnCallUserAbsencesParams{
		Sub:     testAbsenceSub,
		Page:    NewOptInt(2),
		PerPage: NewOptInt(50),
	}
	params.IncludeExpired.SetTo(true)

	absences, err := client.OnCallUserAbsences.List(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, absences, 2)
	assert.Equal(t, testAbsenceUUID, absences[0].UUID)
	assert.Equal(t, time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), absences[0].StartAt)
	assert.Equal(t, time.Date(2026, 12, 3, 0, 0, 0, 0, time.UTC), absences[0].EndAt)
	assert.Equal(t, "none", string(absences[0].RepeatType))
	reason, ok := absences[0].Reason.Get()
	require.True(t, ok)
	assert.Equal(t, "Vacation", reason)
	assert.Equal(t, "weekly", string(absences[1].RepeatType))
	assert.True(t, absences[1].Reason.IsNull())
}

func TestOnCallUserAbsenceService_List_Errors(t *testing.T) {
	for _, tt := range onCallUserAbsenceCommonErrorTests(t) {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/users/{sub}/absences", tt.handler)

			client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
			defer cleanup()

			absences, err := client.OnCallUserAbsences.List(context.Background(), ListOnCallUserAbsencesParams{Sub: testAbsenceSub})

			require.Error(t, err)
			assert.Nil(t, absences)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: list on-call user absences:")
		})
	}
}

func TestOnCallUserAbsenceService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/users/{sub}/absences/{uuid}", onCallUserAbsenceHandler(t, http.MethodGet, testAbsencePath,
		testutil.JSONFixtureHandler(t, "testdata/on_call_user_absence_get_response.json")))

	client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
	defer cleanup()

	absence, err := client.OnCallUserAbsences.Get(context.Background(), GetOnCallUserAbsenceParams{Sub: testAbsenceSub, UUID: testAbsenceUUID})

	require.NoError(t, err)
	require.NotNil(t, absence)
	assert.Equal(t, testAbsenceUUID, absence.UUID)
	assert.Equal(t, time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), absence.StartAt)
	assert.Equal(t, time.Date(2026, 12, 3, 0, 0, 0, 0, time.UTC), absence.EndAt)
	assert.Equal(t, "monthly", string(absence.RepeatType))
	reason, ok := absence.Reason.Get()
	require.True(t, ok)
	assert.Equal(t, "Vacation", reason)
}

func TestOnCallUserAbsenceService_Get_Errors(t *testing.T) {
	for _, tt := range onCallUserAbsenceCommonErrorTests(t) {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/users/{sub}/absences/{uuid}", tt.handler)

			client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
			defer cleanup()

			absence, err := client.OnCallUserAbsences.Get(context.Background(), GetOnCallUserAbsenceParams{Sub: testAbsenceSub, UUID: testAbsenceUUID})

			require.Error(t, err)
			assert.Nil(t, absence)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: get on-call user absence:")
		})
	}
}

func TestOnCallUserAbsenceService_Create(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/users/{sub}/absences", onCallUserAbsenceHandler(t, http.MethodPost, testAbsencesPath, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"absence": {
				"start_at":    "2026-12-01T00:00:00Z",
				"end_at":      "2026-12-03T00:00:00Z",
				"repeat_type": "monthly",
				"reason":      "Vacation",
			},
		}, body)

		testutil.JSONFixtureHandler(t, "testdata/on_call_user_absence_create_response.json", testutil.WithStatus(http.StatusCreated))(w, r)
	}))

	client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
	defer cleanup()

	req := CreateOnCallUserAbsenceReq{
		StartAt: "2026-12-01T00:00:00Z",
		EndAt:   "2026-12-03T00:00:00Z",
		Reason:  NewOptNilString("Vacation"),
	}
	req.RepeatType.SetTo("monthly")

	absence, err := client.OnCallUserAbsences.Create(context.Background(), req, CreateOnCallUserAbsenceParams{Sub: testAbsenceSub})

	require.NoError(t, err)
	require.NotNil(t, absence)
	assert.Equal(t, testAbsenceUUID, absence.UUID)
	assert.Equal(t, time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), absence.StartAt)
	assert.Equal(t, "monthly", string(absence.RepeatType))
}

func TestOnCallUserAbsenceService_Create_Errors(t *testing.T) {
	for _, tt := range onCallUserAbsenceWriteErrorTests(t) {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/users/{sub}/absences", tt.handler)

			client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
			defer cleanup()

			absence, err := client.OnCallUserAbsences.Create(context.Background(), CreateOnCallUserAbsenceReq{
				StartAt: "2026-12-03T00:00:00Z",
				EndAt:   "2026-12-01T00:00:00Z",
			}, CreateOnCallUserAbsenceParams{Sub: testAbsenceSub})

			require.Error(t, err)
			assert.Nil(t, absence)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: create on-call user absence:")
		})
	}
}

func TestOnCallUserAbsenceService_Update(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/users/{sub}/absences/{uuid}", onCallUserAbsenceHandler(t, http.MethodPatch, testAbsencePath, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"absence": {
				"end_at":      "2026-12-04T00:00:00Z",
				"repeat_type": "none",
				"reason":      nil,
			},
		}, body)

		testutil.JSONFixtureHandler(t, "testdata/on_call_user_absence_update_response.json")(w, r)
	}))

	client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
	defer cleanup()

	req := UpdateOnCallUserAbsenceReq{
		EndAt: NewOptString("2026-12-04T00:00:00Z"),
	}
	req.RepeatType.SetTo("none")
	req.Reason.SetToNull()

	absence, err := client.OnCallUserAbsences.Update(context.Background(), req, UpdateOnCallUserAbsenceParams{Sub: testAbsenceSub, UUID: testAbsenceUUID})

	require.NoError(t, err)
	require.NotNil(t, absence)
	assert.Equal(t, testAbsenceUUID, absence.UUID)
	assert.Equal(t, time.Date(2026, 12, 4, 0, 0, 0, 0, time.UTC), absence.EndAt)
	assert.Equal(t, "none", string(absence.RepeatType))
	assert.True(t, absence.Reason.IsNull())
}

func TestOnCallUserAbsenceService_Update_Errors(t *testing.T) {
	for _, tt := range onCallUserAbsenceWriteErrorTests(t) {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/users/{sub}/absences/{uuid}", tt.handler)

			client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
			defer cleanup()

			absence, err := client.OnCallUserAbsences.Update(context.Background(), UpdateOnCallUserAbsenceReq{
				EndAt: NewOptString("2026-11-01T00:00:00Z"),
			}, UpdateOnCallUserAbsenceParams{Sub: testAbsenceSub, UUID: testAbsenceUUID})

			require.Error(t, err)
			assert.Nil(t, absence)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: update on-call user absence:")
		})
	}
}

func TestOnCallUserAbsenceService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/users/{sub}/absences/{uuid}", onCallUserAbsenceHandler(t, http.MethodDelete, testAbsencePath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
	defer cleanup()

	err := client.OnCallUserAbsences.Delete(context.Background(), DeleteOnCallUserAbsenceParams{Sub: testAbsenceSub, UUID: testAbsenceUUID})

	require.NoError(t, err)
}

func TestOnCallUserAbsenceService_Delete_Errors(t *testing.T) {
	for _, tt := range onCallUserAbsenceCommonErrorTests(t) {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/users/{sub}/absences/{uuid}", tt.handler)

			client, cleanup := newOnCallUserAbsenceTestClient(t, mux)
			defer cleanup()

			err := client.OnCallUserAbsences.Delete(context.Background(), DeleteOnCallUserAbsenceParams{Sub: testAbsenceSub, UUID: testAbsenceUUID})

			require.Error(t, err)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: delete on-call user absence:")
		})
	}
}
