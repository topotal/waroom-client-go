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
	mux.HandleFunc("/services/test-service/labels", testutil.JSONFixtureHandler(t, "testdata/labels_list_response.json"))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewLabelService(client)
	labels, err := svc.ListServiceLabels(ctx, ListServiceLabelsParams{
		Name: "test-service",
	})

	require.NoError(t, err)
	assert.Len(t, labels, 2)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), labels[0].UUID)
	assert.Equal(t, "production", labels[0].Name)
	assert.Equal(t, "FF0000", labels[0].Color.Value)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"), labels[1].UUID)
	assert.Equal(t, "staging", labels[1].Name)
	assert.Equal(t, "00FF00", labels[1].Color.Value)
}

func TestLabelService_ListServiceLabels_Errors(t *testing.T) {
	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErrMsg string
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), "unauthorized"},
		{"server error", testutil.InternalServerErrorHandler(), ""},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/services/test-service/labels", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			svc := NewLabelService(client)
			labels, err := svc.ListServiceLabels(ctx, ListServiceLabelsParams{
				Name: "test-service",
			})

			require.Error(t, err)
			assert.Nil(t, labels)
			if tt.wantErrMsg != "" {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestLabelService_CreateServiceLabel(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/services/test-service/labels", testutil.JSONFixtureHandler(
		t,
		"testdata/label_create_response.json",
		testutil.WithStatus(http.StatusCreated),
	))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewLabelService(client)
	label, err := svc.CreateServiceLabel(ctx, CreateServiceLabelReq{
		Name:  "production",
		Color: "FF0000",
	}, CreateServiceLabelParams{
		Name: "test-service",
	})

	require.NoError(t, err)
	require.NotNil(t, label)
	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), label.UUID)
	assert.Equal(t, "production", label.Name)
	assert.Equal(t, "FF0000", label.Color.Value)
}

func TestLabelService_CreateServiceLabel_Errors(t *testing.T) {
	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErrMsg string
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), "unauthorized"},
		{"server error", testutil.InternalServerErrorHandler(), ""},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/services/test-service/labels", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			svc := NewLabelService(client)
			label, err := svc.CreateServiceLabel(ctx, CreateServiceLabelReq{
				Name:  "production",
				Color: "FF0000",
			}, CreateServiceLabelParams{
				Name: "test-service",
			})

			require.Error(t, err)
			assert.Nil(t, label)
			if tt.wantErrMsg != "" {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestLabelService_UpdateServiceLabel(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", testutil.JSONFixtureHandler(
		t,
		"testdata/label_update_response.json",
	))

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewLabelService(client)
	label, err := svc.UpdateServiceLabel(ctx, UpdateServiceLabelReq{
		Name:  NewOptString("production-updated"),
		Color: NewOptString("0000FF"),
	}, UpdateServiceLabelParams{
		Name: "test-service",
		UUID: testUUID,
	})

	require.NoError(t, err)
	require.NotNil(t, label)
	assert.Equal(t, testUUID, label.UUID)
	assert.Equal(t, "production-updated", label.Name)
	assert.Equal(t, "0000FF", label.Color.Value)
}

func TestLabelService_UpdateServiceLabel_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErrMsg string
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), "unauthorized"},
		{"not found error", testutil.NotFoundHandler(t), "not found"},
		{"server error", testutil.InternalServerErrorHandler(), ""},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			svc := NewLabelService(client)
			label, err := svc.UpdateServiceLabel(ctx, UpdateServiceLabelReq{
				Name: NewOptString("production-updated"),
			}, UpdateServiceLabelParams{
				Name: "test-service",
				UUID: testUUID,
			})

			require.Error(t, err)
			assert.Nil(t, label)
			if tt.wantErrMsg != "" {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestLabelService_DeleteServiceLabel(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	mux := http.NewServeMux()
	mux.HandleFunc("/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
		ctx := context.Background()
		return NewClient(ctx, WithBaseURL(baseURL))
	})
	defer cleanup()

	ctx := context.Background()
	svc := NewLabelService(client)
	err := svc.DeleteServiceLabel(ctx, DeleteServiceLabelParams{
		Name: "test-service",
		UUID: testUUID,
	})

	require.NoError(t, err)
}

func TestLabelService_DeleteServiceLabel_Errors(t *testing.T) {
	testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	errorTests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErrMsg string
	}{
		{"unauthorized error", testutil.UnauthorizedHandler(t), "unauthorized"},
		{"not found error", testutil.NotFoundHandler(t), "not found"},
		{"server error", testutil.InternalServerErrorHandler(), ""},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/services/test-service/labels/550e8400-e29b-41d4-a716-446655440001", tt.handler)

			client, cleanup := testutil.SetupTestClient(t, mux, func(baseURL string) (*Client, error) {
				ctx := context.Background()
				return NewClient(ctx, WithBaseURL(baseURL))
			})
			defer cleanup()

			ctx := context.Background()
			svc := NewLabelService(client)
			err := svc.DeleteServiceLabel(ctx, DeleteServiceLabelParams{
				Name: "test-service",
				UUID: testUUID,
			})

			require.Error(t, err)
			if tt.wantErrMsg != "" {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}
