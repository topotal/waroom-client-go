package waroom

import (
	"context"
	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

var _ gen.SecuritySource = (*waroomBearerAuth)(nil)

type waroomBearerAuth struct {
	apiKey string
}

// ApiKeyAuth implements ogen.SecuritySource.
func (auth *waroomBearerAuth) ApiKeyAuth(ctx context.Context, operationName gen.OperationName, client *gen.Client) (gen.ApiKeyAuth, error) {
	return gen.ApiKeyAuth{APIKey: auth.apiKey}, nil
}
