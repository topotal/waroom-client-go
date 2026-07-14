package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

var _ gen.SecuritySource = (*staticTokenSource)(nil)

// staticTokenSource is a SecuritySource that always returns the same API
// key. Additional auth methods can be added later as separate
// SecuritySource implementations behind new With* options.
type staticTokenSource struct {
	apiKey string
}

// BearerAuth implements gen.SecuritySource.
func (s *staticTokenSource) BearerAuth(ctx context.Context, operationName gen.OperationName, client *gen.Client) (gen.BearerAuth, error) {
	return gen.BearerAuth{Token: s.apiKey}, nil
}
