package waroom

import (
	"context"
	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

var _ gen.SecuritySource = (*waroomBearerAuth)(nil)

type waroomBearerAuth struct {
	apiKey string
}

// BearerAuth implements ogen.SecuritySource.
func (auth *waroomBearerAuth) BearerAuth(ctx context.Context, operationName gen.OperationName, client *gen.Client) (gen.BearerAuth, error) {
	return gen.BearerAuth{Token: auth.apiKey}, nil
}
