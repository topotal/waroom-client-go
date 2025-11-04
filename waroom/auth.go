package waroom

import (
	"context"
	gen "github.com/topotal/waroom-client-go/internal/ogen"
	"os"
)

var _ gen.SecuritySource = (*waroomBearerAuth)(nil)

type waroomBearerAuth struct {
}

const WAROOM_API_KEY = "WAROOM_API_KEY"

// ApiKeyAuth implements ogen.SecuritySource.
func (auth *waroomBearerAuth) ApiKeyAuth(ctx context.Context, operationName gen.OperationName, client *gen.Client) (gen.ApiKeyAuth, error) {
	// Implement API key authentication logic here
	return gen.ApiKeyAuth{APIKey: os.Getenv(WAROOM_API_KEY)}, nil
}
