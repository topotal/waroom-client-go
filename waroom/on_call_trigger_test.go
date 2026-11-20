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
	testOnCallTriggerUUID = "770e8400-e29b-41d4-a716-446655440001"
	testOnCallTriggerPath = "/api/v0/on_call/triggers/" + testOnCallTriggerUUID
)

func newOnCallTriggerTestClient(t *testing.T, mux *http.ServeMux) (*Client, func()) {
	t.Helper()
	return testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		return NewClient(context.Background(), WithAPIKey("test-api-key"), WithBaseURL(baseURL))
	})
}

func onCallTriggerBadRequestHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return testutil.JSONFixtureHandler(t, "testdata/on_call_trigger_bad_request_error.json", testutil.WithStatus(http.StatusBadRequest))
}

func onCallTriggerUnprocessableEntityHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return testutil.JSONFixtureHandler(t, "testdata/on_call_trigger_unprocessable_entity_error.json", testutil.WithStatus(http.StatusUnprocessableEntity))
}

func TestOnCallTriggerService_List(t *testing.T) {
	fixture := testutil.JSONFixtureHandler(t, "testdata/on_call_triggers_list_response.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/triggers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "test-service", r.URL.Query().Get("service_name"))
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "50", r.URL.Query().Get("per_page"))
		fixture(w, r)
	})

	client, cleanup := newOnCallTriggerTestClient(t, mux)
	defer cleanup()

	triggers, err := client.OnCallTriggers.List(context.Background(), ListOnCallTriggersParams{
		ServiceName: NewOptString("test-service"),
		Page:        NewOptInt(2),
		PerPage:     NewOptInt(50),
	})

	require.NoError(t, err)
	require.Len(t, triggers, 2)

	assert.Equal(t, testOnCallTriggerUUID, triggers[0].UUID)
	assert.Equal(t, "critical incidents", triggers[0].Name)
	assert.True(t, triggers[0].Enabled)
	assert.Equal(t, "incident_created", string(triggers[0].EventType))
	require.Len(t, triggers[0].SeverityFilters, 2)
	assert.Equal(t, "critical", string(triggers[0].SeverityFilters[0]))
	assert.Equal(t, "high", string(triggers[0].SeverityFilters[1]))
	require.Len(t, triggers[0].SourceFilters, 1)
	assert.Equal(t, "integration", string(triggers[0].SourceFilters[0]))
	assert.Equal(t, 2, triggers[0].RetryCount)
	assert.Equal(t, "test-service", triggers[0].Service.Name)
	require.Len(t, triggers[0].EscalationPolicies, 1)
	assert.Equal(t, "880e8400-e29b-41d4-a716-446655440001", triggers[0].EscalationPolicies[0].UUID)
	assert.Equal(t, "primary", triggers[0].EscalationPolicies[0].Name)

	assert.Equal(t, "770e8400-e29b-41d4-a716-446655440002", triggers[1].UUID)
	assert.False(t, triggers[1].Enabled)
	assert.Len(t, triggers[1].SeverityFilters, 5)
	assert.Len(t, triggers[1].SourceFilters, 3)
	assert.Empty(t, triggers[1].EscalationPolicies)
}

func TestOnCallTriggerService_List_Errors(t *testing.T) {
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
			name:       "server error",
			handler:    testutil.InternalServerErrorHandler(),
			wantAPIErr: nil,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/on_call/triggers", tt.handler)

			client, cleanup := newOnCallTriggerTestClient(t, mux)
			defer cleanup()

			triggers, err := client.OnCallTriggers.List(context.Background(), ListOnCallTriggersParams{})

			require.Error(t, err)
			assert.Nil(t, triggers)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: list on-call triggers:")
		})
	}
}

func TestOnCallTriggerService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testOnCallTriggerPath, testutil.JSONFixtureHandler(t, "testdata/on_call_trigger_get_response.json"))

	client, cleanup := newOnCallTriggerTestClient(t, mux)
	defer cleanup()

	trigger, err := client.OnCallTriggers.Get(context.Background(), GetOnCallTriggerParams{UUID: testOnCallTriggerUUID})

	require.NoError(t, err)
	require.NotNil(t, trigger)
	assert.Equal(t, testOnCallTriggerUUID, trigger.UUID)
	assert.Equal(t, "critical incidents", trigger.Name)
	assert.True(t, trigger.Enabled)
	assert.Equal(t, "incident_created", string(trigger.EventType))
	assert.Len(t, trigger.SeverityFilters, 2)
	assert.Len(t, trigger.SourceFilters, 2)
	assert.Equal(t, 2, trigger.RetryCount)
	assert.Equal(t, "test-service", trigger.Service.Name)
	require.Len(t, trigger.EscalationPolicies, 1)
	assert.Equal(t, "primary", trigger.EscalationPolicies[0].Name)
}

func TestOnCallTriggerService_Get_Errors(t *testing.T) {
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
			mux.HandleFunc(testOnCallTriggerPath, tt.handler)

			client, cleanup := newOnCallTriggerTestClient(t, mux)
			defer cleanup()

			trigger, err := client.OnCallTriggers.Get(context.Background(), GetOnCallTriggerParams{UUID: testOnCallTriggerUUID})

			require.Error(t, err)
			assert.Nil(t, trigger)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: get on-call trigger:")
		})
	}
}

func newCreateOnCallTriggerReq() CreateOnCallTriggerReq {
	return CreateOnCallTriggerReq{
		Name:                  "critical incidents",
		ServiceName:           "test-service",
		EscalationPolicyUuids: []string{"880e8400-e29b-41d4-a716-446655440001"},
		SeverityFilters:       []CreateOnCallTriggerSeverityFilter{"critical", "high"},
		SourceFilters:         []CreateOnCallTriggerSourceFilter{"integration", "slack"},
		Enabled:               NewOptBool(true),
		EventType:             NewOptCreateOnCallTriggerEventType("incident_created"),
		RetryCount:            NewOptInt(1),
	}
}

func TestOnCallTriggerService_Create(t *testing.T) {
	fixture := testutil.JSONFixtureHandler(t, "testdata/on_call_trigger_create_response.json", testutil.WithStatus(http.StatusCreated))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/on_call/triggers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"trigger": {
				"name":                    "critical incidents",
				"service_name":            "test-service",
				"escalation_policy_uuids": []any{"880e8400-e29b-41d4-a716-446655440001"},
				"severity_filters":        []any{"critical", "high"},
				"source_filters":          []any{"integration", "slack"},
				"enabled":                 true,
				"event_type":              "incident_created",
				"retry_count":             float64(1),
			},
		}, body)

		fixture(w, r)
	})

	client, cleanup := newOnCallTriggerTestClient(t, mux)
	defer cleanup()

	trigger, err := client.OnCallTriggers.Create(context.Background(), newCreateOnCallTriggerReq())

	require.NoError(t, err)
	require.NotNil(t, trigger)
	assert.Equal(t, testOnCallTriggerUUID, trigger.UUID)
	assert.Equal(t, "critical incidents", trigger.Name)
	assert.True(t, trigger.Enabled)
	assert.Equal(t, 1, trigger.RetryCount)
	assert.Equal(t, "test-service", trigger.Service.Name)
	require.Len(t, trigger.EscalationPolicies, 1)
	assert.Equal(t, "880e8400-e29b-41d4-a716-446655440001", trigger.EscalationPolicies[0].UUID)
}

func TestOnCallTriggerService_Create_Errors(t *testing.T) {
	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantAPIErr *APIError
	}{
		{
			name:       "bad request error",
			handler:    onCallTriggerBadRequestHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusBadRequest, Messages: []string{"param is missing or the value is empty: trigger"}},
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
			name:       "unprocessable entity error",
			handler:    onCallTriggerUnprocessableEntityHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: []string{"Escalation policy not found"}},
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
			mux.HandleFunc("/api/v0/on_call/triggers", tt.handler)

			client, cleanup := newOnCallTriggerTestClient(t, mux)
			defer cleanup()

			trigger, err := client.OnCallTriggers.Create(context.Background(), newCreateOnCallTriggerReq())

			require.Error(t, err)
			assert.Nil(t, trigger)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: create on-call trigger:")
		})
	}
}

func TestOnCallTriggerService_Update(t *testing.T) {
	fixture := testutil.JSONFixtureHandler(t, "testdata/on_call_trigger_update_response.json")

	mux := http.NewServeMux()
	mux.HandleFunc(testOnCallTriggerPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)

		var body map[string]map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]map[string]any{
			"trigger": {
				"name":        "critical incidents updated",
				"enabled":     false,
				"retry_count": float64(3),
			},
		}, body)

		fixture(w, r)
	})

	client, cleanup := newOnCallTriggerTestClient(t, mux)
	defer cleanup()

	trigger, err := client.OnCallTriggers.Update(context.Background(), UpdateOnCallTriggerReq{
		Name:       NewOptString("critical incidents updated"),
		Enabled:    NewOptBool(false),
		RetryCount: NewOptInt(3),
	}, UpdateOnCallTriggerParams{UUID: testOnCallTriggerUUID})

	require.NoError(t, err)
	require.NotNil(t, trigger)
	assert.Equal(t, testOnCallTriggerUUID, trigger.UUID)
	assert.Equal(t, "critical incidents updated", trigger.Name)
	assert.False(t, trigger.Enabled)
	assert.Equal(t, 3, trigger.RetryCount)
}

func TestOnCallTriggerService_Update_Errors(t *testing.T) {
	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantAPIErr *APIError
	}{
		{
			name:       "bad request error",
			handler:    onCallTriggerBadRequestHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusBadRequest, Messages: []string{"param is missing or the value is empty: trigger"}},
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
			name:       "not found error",
			handler:    testutil.NotFoundHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusNotFound, Messages: []string{"Not found"}},
		},
		{
			name:       "unprocessable entity error",
			handler:    onCallTriggerUnprocessableEntityHandler(t),
			wantAPIErr: &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: []string{"Escalation policy not found"}},
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
			mux.HandleFunc(testOnCallTriggerPath, tt.handler)

			client, cleanup := newOnCallTriggerTestClient(t, mux)
			defer cleanup()

			trigger, err := client.OnCallTriggers.Update(context.Background(), UpdateOnCallTriggerReq{
				Name: NewOptString("critical incidents updated"),
			}, UpdateOnCallTriggerParams{UUID: testOnCallTriggerUUID})

			require.Error(t, err)
			assert.Nil(t, trigger)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: update on-call trigger:")
		})
	}
}

func TestOnCallTriggerService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testOnCallTriggerPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	client, cleanup := newOnCallTriggerTestClient(t, mux)
	defer cleanup()

	err := client.OnCallTriggers.Delete(context.Background(), DeleteOnCallTriggerParams{UUID: testOnCallTriggerUUID})

	require.NoError(t, err)
}

func TestOnCallTriggerService_Delete_Errors(t *testing.T) {
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
			mux.HandleFunc(testOnCallTriggerPath, tt.handler)

			client, cleanup := newOnCallTriggerTestClient(t, mux)
			defer cleanup()

			err := client.OnCallTriggers.Delete(context.Background(), DeleteOnCallTriggerParams{UUID: testOnCallTriggerUUID})

			require.Error(t, err)
			assertAPIError(t, err, tt.wantAPIErr, "waroom: delete on-call trigger:")
		})
	}
}
