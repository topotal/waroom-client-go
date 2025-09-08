package waroom

import (
	"context"
	"os"
	gen "waroom-client-go/internal/ogen"
)

var _ gen.SecuritySource = (*WaroomBearerAuth)(nil)

type WaroomBearerAuth struct {
}

const WAROOM_API_KEY = "WAROOM_API_KEY"

// ApiKeyAuth implements ogen.SecuritySource.
func (auth *WaroomBearerAuth) ApiKeyAuth(ctx context.Context, operationName gen.OperationName, client *gen.Client) (gen.ApiKeyAuth, error) {
	// Implement API key authentication logic here
	return gen.ApiKeyAuth{APIKey: os.Getenv(WAROOM_API_KEY)}, nil
}
