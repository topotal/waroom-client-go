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

func TestLabelService_ListServiceLabels(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels", testutil.JSONFixtureHandler(t, "testdata/labels_list_response.json"))

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	list, err := client.Labels().ListServiceLabels(context.Background(), "test-service", ListServiceLabelsOptions{})

	require.NoError(t, err)
	assert.Len(t, list.Items, 2)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), list.Items[0].UUID)
	assert.Equal(t, "production", list.Items[0].Name)
	assert.Equal(t, "FF0000", list.Items[0].Color)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"), list.Items[1].UUID)
	assert.Equal(t, "staging", list.Items[1].Name)
	assert.Equal(t, "00FF00", list.Items[1].Color)

	assert.Equal(t, 1, list.CurrentPage)
	assert.Equal(t, 1, list.TotalPages)
	assert.Equal(t, 2, list.TotalCount)
	assert.Nil(t, list.NextPage)
	assert.Nil(t, list.PrevPage)
}

func TestLabelService_ListServiceLabelsAll(t *testing.T) {
	page1 := testutil.LoadFixture(t, "testdata/labels_list_page1.json")
	page2 := testutil.LoadFixture(t, "testdata/labels_list_page2.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("page") == "2" {
			_, _ = w.Write(page2)
			return
		}
		_, _ = w.Write(page1)
	})

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	var names []string
	for label, err := range client.Labels().ListServiceLabelsAll(context.Background(), "test-service", ListServiceLabelsOptions{}) {
		require.NoError(t, err)
		names = append(names, label.Name)
	}

	assert.Equal(t, []string{"production", "staging", "backend"}, names)
}

func TestLabelService_ListServiceLabelsAll_Break(t *testing.T) {
	page1 := testutil.LoadFixture(t, "testdata/labels_list_page1.json")

	var requests int
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels", func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(page1)
	})

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	var names []string
	for label, err := range client.Labels().ListServiceLabelsAll(context.Background(), "test-service", ListServiceLabelsOptions{}) {
		require.NoError(t, err)
		names = append(names, label.Name)
		break // 1 件で打ち切り: 次ページは取得されないこと
	}

	assert.Equal(t, []string{"production"}, names)
	assert.Equal(t, 1, requests)
}

func TestLabelService_ListServiceLabelsAll_ErrorStops(t *testing.T) {
	page1 := testutil.LoadFixture(t, "testdata/labels_list_page1.json")

	mux := http.NewServeMux()
	unauthorized := testutil.UnauthorizedHandler(t)
	mux.HandleFunc("/api/v0/services/test-service/labels", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") == "2" {
			unauthorized(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(page1)
	})

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	var names []string
	var iterErr error
	for label, err := range client.Labels().ListServiceLabelsAll(context.Background(), "test-service", ListServiceLabelsOptions{}) {
		if err != nil {
			iterErr = err
			continue // エラー yield 後にイテレータが停止すること(continue しても再開しない)
		}
		names = append(names, label.Name)
	}

	assert.Equal(t, []string{"production", "staging"}, names)
	require.Error(t, iterErr)
	assert.ErrorIs(t, iterErr, ErrUnauthorized)
}

func TestLabelService_ListServiceLabels_Errors(t *testing.T) {
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
			mux.HandleFunc("/api/v0/services/test-service/labels", tt.handler)

			client, cleanup := setupTestClient(t, mux)
			defer cleanup()

			list, err := client.Labels().ListServiceLabels(context.Background(), "test-service", ListServiceLabelsOptions{})

			require.Error(t, err)
			assert.Nil(t, list)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestLabelService_CreateServiceLabel(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels", testutil.JSONFixtureHandler(
		t,
		"testdata/label_create_response.json",
		testutil.WithStatus(http.StatusCreated),
	))

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	label, err := client.Labels().CreateServiceLabel(context.Background(), "test-service", CreateServiceLabelRequest{
		Name:  "production",
		Color: "FF0000",
	})

	require.NoError(t, err)
	require.NotNil(t, label)
	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), label.UUID)
	assert.Equal(t, "production", label.Name)
	assert.Equal(t, "FF0000", label.Color)
}

func TestLabelService_CreateServiceLabel_Errors(t *testing.T) {
	unprocessableHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error": "Validation failed"}`))
	}

	errorTests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr error
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), ErrUnauthorized},
		{"unprocessable entity error", unprocessableHandler, ErrUnprocessableEntity},
		{"server error", testutil.InternalServerErrorHandler(), nil},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/services/test-service/labels", tt.handler)

			client, cleanup := setupTestClient(t, mux)
			defer cleanup()

			label, err := client.Labels().CreateServiceLabel(context.Background(), "test-service", CreateServiceLabelRequest{
				Name:  "production",
				Color: "FF0000",
			})

			require.Error(t, err)
			assert.Nil(t, label)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestLabelService_CreateServiceLabel_UnprocessableEntityMessage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error": "Validation failed: Name has already been taken"}`))
	})

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	_, err := client.Labels().CreateServiceLabel(context.Background(), "test-service", CreateServiceLabelRequest{
		Name:  "production",
		Color: "FF0000",
	})

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 422, apiErr.StatusCode)
	assert.Equal(t, []string{"Validation failed: Name has already been taken"}, apiErr.Messages)
}

func TestLabelService_UpdateServiceLabel(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", testutil.JSONFixtureHandler(
		t,
		"testdata/label_update_response.json",
	))

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	label, err := client.Labels().UpdateServiceLabel(context.Background(), "test-service", testUUID, UpdateServiceLabelRequest{
		Name:  NewOptString("production-updated"),
		Color: NewOptString("0000FF"),
	})

	require.NoError(t, err)
	require.NotNil(t, label)
	assert.Equal(t, testUUID, label.UUID)
	assert.Equal(t, "production-updated", label.Name)
	assert.Equal(t, "0000FF", label.Color)
}

func TestLabelService_UpdateServiceLabel_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	errorTests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr error
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), ErrUnauthorized},
		{"not found error", testutil.NotFoundHandler(t), ErrNotFound},
		{"server error", testutil.InternalServerErrorHandler(), nil},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", tt.handler)

			client, cleanup := setupTestClient(t, mux)
			defer cleanup()

			label, err := client.Labels().UpdateServiceLabel(context.Background(), "test-service", testUUID, UpdateServiceLabelRequest{
				Name: NewOptString("production-updated"),
			})

			require.Error(t, err)
			assert.Nil(t, label)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestLabelService_DeleteServiceLabel(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	client, cleanup := setupTestClient(t, mux)
	defer cleanup()

	err := client.Labels().DeleteServiceLabel(context.Background(), "test-service", testUUID)

	require.NoError(t, err)
}

func TestLabelService_DeleteServiceLabel_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	errorTests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr error
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), ErrUnauthorized},
		{"not found error", testutil.NotFoundHandler(t), ErrNotFound},
		{"server error", testutil.InternalServerErrorHandler(), nil},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v0/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", tt.handler)

			client, cleanup := setupTestClient(t, mux)
			defer cleanup()

			err := client.Labels().DeleteServiceLabel(context.Background(), "test-service", testUUID)

			require.Error(t, err)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
