package waroom

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topotal/waroom-client-go/internal/testutil"
)

const (
	testOnCallTeamUUID = "550e8400-e29b-41d4-a716-446655440001"
	testOnCallTeamPath = "/api/v0/on_call/teams/" + testOnCallTeamUUID
)

func newOnCallTeamTestClient(t *testing.T, mux *http.ServeMux) (*Client, func()) {
	t.Helper()
	return testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
}

func onCallTeamBadRequestHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return testutil.JSONFixtureHandler(t, "testdata/on_call_team_bad_request_error.json", testutil.WithStatus(http.StatusBadRequest))
}

func onCallTeamUnprocessableEntityHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return testutil.JSONFixtureHandler(t, "testdata/on_call_team_unprocessable_entity_error.json", testutil.WithStatus(http.StatusUnprocessableEntity))
}

var (
	onCallTeamUnauthorizedErr        = &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}
	onCallTeamPaymentRequiredErr     = &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{"Payment Required"}}
	onCallTeamNotFoundErr            = &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}}
	onCallTeamBadRequestErr          = &APIError{StatusCode: http.StatusBadRequest, Messages: []string{"param is missing or the value is empty: team"}}
	onCallTeamUnprocessableEntityErr = &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: []string{"Participants user must have a verified phone number"}}
)

type onCallTeamErrorCase struct {
	name       string
	handler    func(t *testing.T) http.HandlerFunc
	wantAPIErr *APIError
}

func onCallTeamServerErrorHandler(*testing.T) http.HandlerFunc {
	return testutil.InternalServerErrorHandler()
}

func TestOnCallTeamService_List(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/teams", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "50", r.URL.Query().Get("per_page"))
		testutil.JSONFixtureHandler(t, "testdata/on_call_team_list_response.json")(w, r)
	})

	client, cleanup := newOnCallTeamTestClient(t, mux)
	defer cleanup()

	teams, err := client.OnCallTeams.List(context.Background(), ListOnCallTeamsParams{
		Page:    NewOptInt(2),
		PerPage: NewOptInt(50),
	})

	require.NoError(t, err)
	require.Len(t, teams, 2)

	assert.Equal(t, testOnCallTeamUUID, teams[0].UUID)
	assert.Equal(t, "SRE", teams[0].Name)
	assert.Equal(t, "FF0000", teams[0].Color)
	require.Len(t, teams[0].Rotations, 1)
	assert.Equal(t, "Asia/Tokyo", teams[0].Rotations[0].Timezone)
	assert.Len(t, teams[0].Rotations[0].Participants, 2)
	onDuty, ok := teams[0].OnDutyParticipant.Get()
	require.True(t, ok)
	assert.Equal(t, "alice", onDuty.User.Nickname)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440002", teams[1].UUID)
	assert.Equal(t, "Backend", teams[1].Name)
	assert.True(t, teams[1].Rotation.Null)
	assert.Empty(t, teams[1].Rotations)
	assert.True(t, teams[1].OnDutyParticipant.Null)
}

func TestOnCallTeamService_List_Errors(t *testing.T) {
	tests := []onCallTeamErrorCase{
		{name: "unauthorized error", handler: testutil.UnauthorizedHandler, wantAPIErr: onCallTeamUnauthorizedErr},
		{name: "payment required error", handler: testutil.PaymentRequiredHandler, wantAPIErr: onCallTeamPaymentRequiredErr},
		{name: "server error", handler: onCallTeamServerErrorHandler, wantAPIErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/teams", tt.handler(t))

			client, cleanup := newOnCallTeamTestClient(t, mux)
			defer cleanup()

			teams, err := client.OnCallTeams.List(context.Background(), ListOnCallTeamsParams{})

			require.Error(t, err)
			assert.Nil(t, teams)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: list on-call teams:")
		})
	}
}

func TestOnCallTeamService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testOnCallTeamPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		testutil.JSONFixtureHandler(t, "testdata/on_call_team_get_response.json")(w, r)
	})

	client, cleanup := newOnCallTeamTestClient(t, mux)
	defer cleanup()

	team, err := client.OnCallTeams.Get(context.Background(), GetOnCallTeamParams{UUID: testOnCallTeamUUID})

	require.NoError(t, err)
	require.NotNil(t, team)
	assert.Equal(t, testOnCallTeamUUID, team.UUID)
	assert.Equal(t, "SRE", team.Name)
	assert.Equal(t, "FF0000", team.Color)

	rotation, ok := team.Rotation.Get()
	require.True(t, ok)
	assert.Equal(t, "660e8400-e29b-41d4-a716-446655440001", rotation.UUID)
	assert.Equal(t, 7, rotation.IntervalDays)
	wday, ok := rotation.ShiftHandoffWday.Get()
	require.True(t, ok)
	assert.Equal(t, 1, wday)
	assert.Equal(t, "10:00", rotation.ShiftHandoffTime)
	require.Len(t, rotation.TimeSlots, 1)
	assert.Equal(t, "09:00", rotation.TimeSlots[0].StartTime)
	require.Len(t, rotation.Participants, 2)
	assert.Equal(t, "auth0|alice", rotation.Participants[0].User.Sub)
	slackUser, ok := rotation.Participants[0].User.SlackUser.Get()
	require.True(t, ok)
	assert.Equal(t, "U00000001", slackUser.Key)
	assert.True(t, rotation.Participants[1].Absent)
	assert.True(t, rotation.Participants[1].User.SlackUser.Null)

	require.Len(t, team.Rotations, 1)
	onDuty, ok := team.OnDutyParticipant.Get()
	require.True(t, ok)
	assert.Equal(t, "alice", onDuty.User.Nickname)
	assert.EqualValues(t, "member", onDuty.Role)
}

func TestOnCallTeamService_Get_Errors(t *testing.T) {
	tests := []onCallTeamErrorCase{
		{name: "unauthorized error", handler: testutil.UnauthorizedHandler, wantAPIErr: onCallTeamUnauthorizedErr},
		{name: "payment required error", handler: testutil.PaymentRequiredHandler, wantAPIErr: onCallTeamPaymentRequiredErr},
		{name: "not found error", handler: testutil.NotFoundHandler, wantAPIErr: onCallTeamNotFoundErr},
		{name: "server error", handler: onCallTeamServerErrorHandler, wantAPIErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testOnCallTeamPath, tt.handler(t))

			client, cleanup := newOnCallTeamTestClient(t, mux)
			defer cleanup()

			team, err := client.OnCallTeams.Get(context.Background(), GetOnCallTeamParams{UUID: testOnCallTeamUUID})

			require.Error(t, err)
			assert.Nil(t, team)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: get on-call team:")
		})
	}
}

func newCreateOnCallTeamReq() CreateOnCallTeamReq {
	return CreateOnCallTeamReq{
		Name:  "SRE",
		Color: "FF0000",
		Rotations: []CreateOnCallTeamRotationReq{
			{
				Timezone:         "Asia/Tokyo",
				StartedAt:        "2026-11-02T10:00:00+09:00",
				ShiftHandoffTime: "10:00",
				ShiftHandoffWday: NewOptNilInt(1),
				IntervalDays:     NewOptInt(7),
				Participants: []CreateOnCallTeamParticipantReq{
					{UserSub: "auth0|alice", Position: 0},
					{UserSub: "auth0|bob", Position: 1},
				},
				TimeSlots: []CreateOnCallTeamTimeSlotReq{
					{DayOfWeek: 1, StartTime: "09:00", EndTime: "18:00"},
				},
			},
		},
	}
}

func TestOnCallTeamService_Create(t *testing.T) {
	fixture := testutil.JSONFixtureHandler(t, "testdata/on_call_team_create_response.json", testutil.WithStatus(http.StatusCreated))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/teams", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{
			"team": map[string]any{
				"name":  "SRE",
				"color": "FF0000",
				"rotations": []any{
					map[string]any{
						"timezone":           "Asia/Tokyo",
						"started_at":         "2026-11-02T10:00:00+09:00",
						"shift_handoff_time": "10:00",
						"shift_handoff_wday": float64(1),
						"interval_days":      float64(7),
						"participants": []any{
							map[string]any{"user_sub": "auth0|alice", "position": float64(0)},
							map[string]any{"user_sub": "auth0|bob", "position": float64(1)},
						},
						"time_slots": []any{
							map[string]any{"day_of_week": float64(1), "start_time": "09:00", "end_time": "18:00"},
						},
					},
				},
			},
		}, body)

		fixture(w, r)
	})

	client, cleanup := newOnCallTeamTestClient(t, mux)
	defer cleanup()

	team, err := client.OnCallTeams.Create(context.Background(), newCreateOnCallTeamReq())

	require.NoError(t, err)
	require.NotNil(t, team)
	assert.Equal(t, testOnCallTeamUUID, team.UUID)
	assert.Equal(t, "SRE", team.Name)
	assert.Equal(t, "FF0000", team.Color)
	require.Len(t, team.Rotations, 1)
	assert.Len(t, team.Rotations[0].Participants, 2)
}

func TestOnCallTeamService_Create_Errors(t *testing.T) {
	tests := []onCallTeamErrorCase{
		{name: "bad request error", handler: onCallTeamBadRequestHandler, wantAPIErr: onCallTeamBadRequestErr},
		{name: "unauthorized error", handler: testutil.UnauthorizedHandler, wantAPIErr: onCallTeamUnauthorizedErr},
		{name: "payment required error", handler: testutil.PaymentRequiredHandler, wantAPIErr: onCallTeamPaymentRequiredErr},
		{name: "unprocessable entity error", handler: onCallTeamUnprocessableEntityHandler, wantAPIErr: onCallTeamUnprocessableEntityErr},
		{name: "server error", handler: onCallTeamServerErrorHandler, wantAPIErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/teams", tt.handler(t))

			client, cleanup := newOnCallTeamTestClient(t, mux)
			defer cleanup()

			team, err := client.OnCallTeams.Create(context.Background(), newCreateOnCallTeamReq())

			require.Error(t, err)
			assert.Nil(t, team)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: create on-call team:")
		})
	}
}

func TestOnCallTeamService_Update(t *testing.T) {
	fixture := testutil.JSONFixtureHandler(t, "testdata/on_call_team_update_response.json")

	mux := http.NewServeMux()
	mux.HandleFunc(testOnCallTeamPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		// rotations is omitted when not set, so the existing rotations are kept.
		assert.Equal(t, map[string]any{
			"team": map[string]any{
				"name":  "SRE Primary",
				"color": "00FF00",
			},
		}, body)

		fixture(w, r)
	})

	client, cleanup := newOnCallTeamTestClient(t, mux)
	defer cleanup()

	team, err := client.OnCallTeams.Update(context.Background(), UpdateOnCallTeamReq{
		Name:  NewOptString("SRE Primary"),
		Color: NewOptString("00FF00"),
	}, UpdateOnCallTeamParams{UUID: testOnCallTeamUUID})

	require.NoError(t, err)
	require.NotNil(t, team)
	assert.Equal(t, testOnCallTeamUUID, team.UUID)
	assert.Equal(t, "SRE Primary", team.Name)
	assert.Equal(t, "00FF00", team.Color)
}

func TestOnCallTeamService_Update_Rotations(t *testing.T) {
	fixture := testutil.JSONFixtureHandler(t, "testdata/on_call_team_update_response.json")

	mux := http.NewServeMux()
	mux.HandleFunc(testOnCallTeamPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{
			"team": map[string]any{
				"rotations": []any{
					map[string]any{
						"uuid":               "660e8400-e29b-41d4-a716-446655440001",
						"timezone":           "Asia/Tokyo",
						"started_at":         "2026-11-02T10:00:00+09:00",
						"shift_handoff_time": "10:00",
						"shift_handoff_wday": nil,
						"interval_days":      float64(1),
						"participants":       []any{},
					},
				},
			},
		}, body)

		fixture(w, r)
	})

	client, cleanup := newOnCallTeamTestClient(t, mux)
	defer cleanup()

	rotation := UpdateOnCallTeamRotationReq{
		UUID:             NewOptString("660e8400-e29b-41d4-a716-446655440001"),
		Timezone:         "Asia/Tokyo",
		StartedAt:        "2026-11-02T10:00:00+09:00",
		ShiftHandoffTime: "10:00",
		IntervalDays:     NewOptInt(1),
		Participants:     []UpdateOnCallTeamParticipantReq{},
	}
	rotation.ShiftHandoffWday.SetToNull()

	team, err := client.OnCallTeams.Update(context.Background(), UpdateOnCallTeamReq{
		Rotations: []UpdateOnCallTeamRotationReq{rotation},
	}, UpdateOnCallTeamParams{UUID: testOnCallTeamUUID})

	require.NoError(t, err)
	require.NotNil(t, team)
}

func TestOnCallTeamService_Update_Errors(t *testing.T) {
	tests := []onCallTeamErrorCase{
		{name: "bad request error", handler: onCallTeamBadRequestHandler, wantAPIErr: onCallTeamBadRequestErr},
		{name: "unauthorized error", handler: testutil.UnauthorizedHandler, wantAPIErr: onCallTeamUnauthorizedErr},
		{name: "payment required error", handler: testutil.PaymentRequiredHandler, wantAPIErr: onCallTeamPaymentRequiredErr},
		{name: "not found error", handler: testutil.NotFoundHandler, wantAPIErr: onCallTeamNotFoundErr},
		{name: "unprocessable entity error", handler: onCallTeamUnprocessableEntityHandler, wantAPIErr: onCallTeamUnprocessableEntityErr},
		{name: "server error", handler: onCallTeamServerErrorHandler, wantAPIErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testOnCallTeamPath, tt.handler(t))

			client, cleanup := newOnCallTeamTestClient(t, mux)
			defer cleanup()

			team, err := client.OnCallTeams.Update(context.Background(), UpdateOnCallTeamReq{
				Name: NewOptString("SRE Primary"),
			}, UpdateOnCallTeamParams{UUID: testOnCallTeamUUID})

			require.Error(t, err)
			assert.Nil(t, team)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: update on-call team:")
		})
	}
}

func TestOnCallTeamService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testOnCallTeamPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	client, cleanup := newOnCallTeamTestClient(t, mux)
	defer cleanup()

	err := client.OnCallTeams.Delete(context.Background(), DeleteOnCallTeamParams{UUID: testOnCallTeamUUID})

	require.NoError(t, err)
}

func TestOnCallTeamService_Delete_Errors(t *testing.T) {
	tests := []onCallTeamErrorCase{
		{name: "unauthorized error", handler: testutil.UnauthorizedHandler, wantAPIErr: onCallTeamUnauthorizedErr},
		{name: "payment required error", handler: testutil.PaymentRequiredHandler, wantAPIErr: onCallTeamPaymentRequiredErr},
		{name: "not found error", handler: testutil.NotFoundHandler, wantAPIErr: onCallTeamNotFoundErr},
		{name: "unprocessable entity error", handler: onCallTeamUnprocessableEntityHandler, wantAPIErr: onCallTeamUnprocessableEntityErr},
		{name: "server error", handler: onCallTeamServerErrorHandler, wantAPIErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testOnCallTeamPath, tt.handler(t))

			client, cleanup := newOnCallTeamTestClient(t, mux)
			defer cleanup()

			err := client.OnCallTeams.Delete(context.Background(), DeleteOnCallTeamParams{UUID: testOnCallTeamUUID})

			require.Error(t, err)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: delete on-call team:")
		})
	}
}
