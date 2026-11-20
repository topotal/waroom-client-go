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
	testEscalationPolicyUUID = "770e8400-e29b-41d4-a716-446655440001"
	testEscalationPolicyPath = "/api/v0/on_call/escalation_policies/" + testEscalationPolicyUUID
)

func setupOnCallEscalationPolicyTestClient(t *testing.T, path string, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc(path, handler)
	return testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
}

func onCallEscalationPolicyBadRequestHandler(t *testing.T) http.HandlerFunc {
	return testutil.JSONFixtureHandler(t, "testdata/on_call_escalation_policy_bad_request_error.json", testutil.WithStatus(http.StatusBadRequest))
}

func onCallEscalationPolicyUnprocessableEntityHandler(t *testing.T) http.HandlerFunc {
	return testutil.JSONFixtureHandler(t, "testdata/on_call_escalation_policy_unprocessable_entity_error.json", testutil.WithStatus(http.StatusUnprocessableEntity))
}

type onCallEscalationPolicyErrorTest struct {
	name       string
	handler    http.HandlerFunc
	wantAPIErr *APIError
}

func onCallEscalationPolicyErrorTests(t *testing.T, statuses ...int) []onCallEscalationPolicyErrorTest {
	tests := []onCallEscalationPolicyErrorTest{}
	for _, status := range statuses {
		switch status {
		case http.StatusBadRequest:
			tests = append(tests, onCallEscalationPolicyErrorTest{"bad request error", onCallEscalationPolicyBadRequestHandler(t), &APIError{StatusCode: http.StatusBadRequest, Messages: []string{"Bad request"}}})
		case http.StatusUnauthorized:
			tests = append(tests, onCallEscalationPolicyErrorTest{"unauthorized error", testutil.UnauthorizedHandler(t), &APIError{StatusCode: http.StatusUnauthorized, Messages: []string{"Not Authenticated"}}})
		case http.StatusPaymentRequired:
			tests = append(tests, onCallEscalationPolicyErrorTest{"payment required error", testutil.PaymentRequiredHandler(t), &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{"Payment Required"}}})
		case http.StatusNotFound:
			tests = append(tests, onCallEscalationPolicyErrorTest{"not found error", testutil.NotFoundHandler(t), &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}}})
		case http.StatusUnprocessableEntity:
			tests = append(tests, onCallEscalationPolicyErrorTest{"unprocessable entity error", onCallEscalationPolicyUnprocessableEntityHandler(t), &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: []string{"Name can't be blank"}}})
		}
	}
	return append(tests, onCallEscalationPolicyErrorTest{"server error", testutil.InternalServerErrorHandler(), nil})
}

func TestOnCallEscalationPolicyService_List(t *testing.T) {
	client, cleanup := setupOnCallEscalationPolicyTestClient(t, "/api/v0/on_call/escalation_policies", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "50", r.URL.Query().Get("per_page"))
		testutil.JSONFixtureHandler(t, "testdata/on_call_escalation_policy_list_response.json")(w, r)
	})
	defer cleanup()

	policies, err := client.OnCallEscalationPolicies.List(context.Background(), ListOnCallEscalationPoliciesParams{
		Page:    NewOptInt(2),
		PerPage: NewOptInt(50),
	})

	require.NoError(t, err)
	require.Len(t, policies, 2)

	assert.Equal(t, testEscalationPolicyUUID, policies[0].UUID)
	assert.Equal(t, "Primary escalation", policies[0].Name)
	assert.Equal(t, "Asia/Tokyo", policies[0].Timezone)
	assert.Equal(t, "all_members", string(policies[0].TeamCallMode))
	require.Len(t, policies[0].Steps, 2)
	assert.Equal(t, 1, policies[0].Steps[1].Position)
	assert.Equal(t, "fallback", string(policies[0].Steps[1].StepType))
	assert.True(t, policies[0].Steps[1].SimultaneousCall)
	assert.Equal(t, "Backend", policies[0].Steps[1].Team.Name)
	require.Len(t, policies[0].TimeSlots, 1)
	assert.Equal(t, 1, policies[0].TimeSlots[0].DayOfWeek)
	assert.Equal(t, "09:00", policies[0].TimeSlots[0].StartTime)
	assert.Equal(t, "18:00", policies[0].TimeSlots[0].EndTime)

	assert.Equal(t, "Night shift", policies[1].Name)
	assert.Equal(t, "on_duty_only", string(policies[1].TeamCallMode))
	assert.Empty(t, policies[1].Steps)
}

func TestOnCallEscalationPolicyService_List_Errors(t *testing.T) {
	for _, tt := range onCallEscalationPolicyErrorTests(t, http.StatusUnauthorized, http.StatusPaymentRequired) {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupOnCallEscalationPolicyTestClient(t, "/api/v0/on_call/escalation_policies", tt.handler)
			defer cleanup()

			policies, err := client.OnCallEscalationPolicies.List(context.Background(), ListOnCallEscalationPoliciesParams{})

			require.Error(t, err)
			assert.Nil(t, policies)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: list on-call escalation policies:")
		})
	}
}

func TestOnCallEscalationPolicyService_Get(t *testing.T) {
	client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		testutil.JSONFixtureHandler(t, "testdata/on_call_escalation_policy_get_response.json")(w, r)
	})
	defer cleanup()

	policy, err := client.OnCallEscalationPolicies.Get(context.Background(), GetOnCallEscalationPolicyParams{UUID: testEscalationPolicyUUID})

	require.NoError(t, err)
	require.NotNil(t, policy)
	assert.Equal(t, testEscalationPolicyUUID, policy.UUID)
	assert.Equal(t, "Primary escalation", policy.Name)
	assert.Equal(t, "all_members", string(policy.TeamCallMode))
	require.Len(t, policy.Steps, 2)
	assert.Equal(t, "primary", string(policy.Steps[0].StepType))
	assert.Equal(t, "880e8400-e29b-41d4-a716-446655440001", policy.Steps[0].Team.UUID)
	assert.Equal(t, "FF0000", policy.Steps[0].Team.Color)
	require.Len(t, policy.TimeSlots, 1)
	assert.Equal(t, "770e8400-e29b-41d4-a716-446655440021", policy.TimeSlots[0].UUID)
}

func TestOnCallEscalationPolicyService_Get_Errors(t *testing.T) {
	for _, tt := range onCallEscalationPolicyErrorTests(t, http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusNotFound) {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath, tt.handler)
			defer cleanup()

			policy, err := client.OnCallEscalationPolicies.Get(context.Background(), GetOnCallEscalationPolicyParams{UUID: testEscalationPolicyUUID})

			require.Error(t, err)
			assert.Nil(t, policy)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: get on-call escalation policy:")
		})
	}
}

func TestOnCallEscalationPolicyService_Create(t *testing.T) {
	client, cleanup := setupOnCallEscalationPolicyTestClient(t, "/api/v0/on_call/escalation_policies", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"escalation_policy": {
				"name":           "Primary escalation",
				"timezone":       "Asia/Tokyo",
				"team_call_mode": "all_members",
				"steps": []any{
					map[string]any{"team_uuid": "880e8400-e29b-41d4-a716-446655440001"},
					map[string]any{"team_uuid": "880e8400-e29b-41d4-a716-446655440002", "step_type": "fallback", "simultaneous_call": true},
				},
				"time_slots": []any{
					map[string]any{"day_of_week": float64(1), "start_time": "09:00", "end_time": "18:00"},
				},
			},
		}, body)

		testutil.JSONFixtureHandler(t, "testdata/on_call_escalation_policy_create_response.json", testutil.WithStatus(http.StatusCreated))(w, r)
	})
	defer cleanup()

	fallback := CreateOnCallEscalationPolicyStep{TeamUUID: "880e8400-e29b-41d4-a716-446655440002"}
	stepType := "fallback"
	fallback.StepType.SetTo(CreateOnCallEscalationPolicyStepType(stepType))
	fallback.SimultaneousCall.SetTo(true)
	req := CreateOnCallEscalationPolicyReq{
		Name:     "Primary escalation",
		Timezone: NewOptString("Asia/Tokyo"),
		Steps: []CreateOnCallEscalationPolicyStep{
			{TeamUUID: "880e8400-e29b-41d4-a716-446655440001"},
			fallback,
		},
		TimeSlots: []CreateOnCallEscalationPolicyTimeSlot{
			{DayOfWeek: 1, StartTime: "09:00", EndTime: "18:00"},
		},
	}
	req.TeamCallMode.SetTo("all_members")

	policy, err := client.OnCallEscalationPolicies.Create(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, policy)
	assert.Equal(t, testEscalationPolicyUUID, policy.UUID)
	assert.Equal(t, "Primary escalation", policy.Name)
	require.Len(t, policy.Steps, 2)
	assert.Equal(t, "fallback", string(policy.Steps[1].StepType))
}

func TestOnCallEscalationPolicyService_Create_Errors(t *testing.T) {
	for _, tt := range onCallEscalationPolicyErrorTests(t, http.StatusBadRequest, http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusUnprocessableEntity) {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupOnCallEscalationPolicyTestClient(t, "/api/v0/on_call/escalation_policies", tt.handler)
			defer cleanup()

			policy, err := client.OnCallEscalationPolicies.Create(context.Background(), CreateOnCallEscalationPolicyReq{
				Name:      "Primary escalation",
				TimeSlots: []CreateOnCallEscalationPolicyTimeSlot{{DayOfWeek: 1, StartTime: "09:00", EndTime: "18:00"}},
			})

			require.Error(t, err)
			assert.Nil(t, policy)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: create on-call escalation policy:")
		})
	}
}

func TestOnCallEscalationPolicyService_Update(t *testing.T) {
	client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)

		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"escalation_policy": {
				"name":           "Primary escalation",
				"team_call_mode": "on_duty_only",
				"steps": []any{
					map[string]any{"team_uuid": "880e8400-e29b-41d4-a716-446655440001", "step_type": "primary"},
				},
			},
		}, body)

		testutil.JSONFixtureHandler(t, "testdata/on_call_escalation_policy_update_response.json")(w, r)
	})
	defer cleanup()

	step := UpdateOnCallEscalationPolicyStep{TeamUUID: "880e8400-e29b-41d4-a716-446655440001"}
	step.StepType.SetTo("primary")
	req := UpdateOnCallEscalationPolicyReq{
		Name:  NewOptString("Primary escalation"),
		Steps: []UpdateOnCallEscalationPolicyStep{step},
	}
	req.TeamCallMode.SetTo("on_duty_only")

	policy, err := client.OnCallEscalationPolicies.Update(context.Background(), req, UpdateOnCallEscalationPolicyParams{UUID: testEscalationPolicyUUID})

	require.NoError(t, err)
	require.NotNil(t, policy)
	assert.Equal(t, testEscalationPolicyUUID, policy.UUID)
	assert.Equal(t, "Primary escalation", policy.Name)
	require.Len(t, policy.TimeSlots, 1)
	assert.Equal(t, "09:00", policy.TimeSlots[0].StartTime)
}

func TestOnCallEscalationPolicyService_Update_Errors(t *testing.T) {
	for _, tt := range onCallEscalationPolicyErrorTests(t, http.StatusBadRequest, http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusNotFound, http.StatusUnprocessableEntity) {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath, tt.handler)
			defer cleanup()

			policy, err := client.OnCallEscalationPolicies.Update(context.Background(), UpdateOnCallEscalationPolicyReq{
				Name: NewOptString("Primary escalation"),
			}, UpdateOnCallEscalationPolicyParams{UUID: testEscalationPolicyUUID})

			require.Error(t, err)
			assert.Nil(t, policy)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: update on-call escalation policy:")
		})
	}
}

func TestOnCallEscalationPolicyService_Delete(t *testing.T) {
	client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	defer cleanup()

	err := client.OnCallEscalationPolicies.Delete(context.Background(), DeleteOnCallEscalationPolicyParams{UUID: testEscalationPolicyUUID})

	require.NoError(t, err)
}

func TestOnCallEscalationPolicyService_Delete_Errors(t *testing.T) {
	for _, tt := range onCallEscalationPolicyErrorTests(t, http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusNotFound, http.StatusUnprocessableEntity) {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath, tt.handler)
			defer cleanup()

			err := client.OnCallEscalationPolicies.Delete(context.Background(), DeleteOnCallEscalationPolicyParams{UUID: testEscalationPolicyUUID})

			require.Error(t, err)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: delete on-call escalation policy:")
		})
	}
}

func TestOnCallEscalationPolicyService_GetSchedule(t *testing.T) {
	client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath+"/schedule", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "2026-11-23", r.URL.Query().Get("start_date"))
		assert.Equal(t, "2026-11-24", r.URL.Query().Get("end_date"))
		testutil.JSONFixtureHandler(t, "testdata/on_call_escalation_policy_schedule_response.json")(w, r)
	})
	defer cleanup()

	params := GetOnCallEscalationPolicyScheduleParams{UUID: testEscalationPolicyUUID}
	params.StartDate.SetTo(time.Date(2026, 11, 23, 0, 0, 0, 0, time.UTC))
	params.EndDate.SetTo(time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC))

	schedule, err := client.OnCallEscalationPolicies.GetSchedule(context.Background(), params)

	require.NoError(t, err)
	require.NotNil(t, schedule)
	assert.Equal(t, testEscalationPolicyUUID, schedule.EscalationPolicy.UUID)
	assert.Equal(t, "Asia/Tokyo", schedule.EscalationPolicy.Timezone)
	require.Len(t, schedule.Days, 2)

	day := schedule.Days[0]
	assert.Equal(t, "2026-11-23", day.Date.Format(time.DateOnly))
	assert.Equal(t, 1, day.DayOfWeek)
	assert.True(t, day.Active)
	require.Len(t, day.Entries, 1)

	entry := day.Entries[0]
	assert.Equal(t, 0, entry.EscalationPolicyStep.Position)
	assert.Equal(t, "SRE", entry.EscalationPolicyStep.Team.Name)
	onDuty, ok := entry.OnDutyUser.Get()
	require.True(t, ok)
	assert.Equal(t, "alice", onDuty.Nickname)
	assert.True(t, onDuty.PhoneVerified)
	assert.True(t, entry.Override.IsNull())
	jst := time.FixedZone("", 9*60*60)
	assert.True(t, time.Date(2026, 11, 23, 9, 0, 0, 0, jst).Equal(entry.StartAt))
	assert.True(t, time.Date(2026, 11, 23, 18, 0, 0, 0, jst).Equal(entry.EndAt))

	assert.False(t, schedule.Days[1].Active)
	assert.Empty(t, schedule.Days[1].Entries)
}

func TestOnCallEscalationPolicyService_GetSchedule_Errors(t *testing.T) {
	for _, tt := range onCallEscalationPolicyErrorTests(t, http.StatusBadRequest, http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusNotFound) {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupOnCallEscalationPolicyTestClient(t, testEscalationPolicyPath+"/schedule", tt.handler)
			defer cleanup()

			schedule, err := client.OnCallEscalationPolicies.GetSchedule(context.Background(), GetOnCallEscalationPolicyScheduleParams{UUID: testEscalationPolicyUUID})

			require.Error(t, err)
			assert.Nil(t, schedule)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: get on-call escalation policy schedule:")
		})
	}
}
