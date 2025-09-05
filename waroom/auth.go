package waroom

import (
	"context"
	gen "waroom-client-go/internal/ogen"
)

var _ gen.SecuritySource = (*WaroomBearerAuth)(nil)

type WaroomBearerAuth struct{}

func (s *WaroomBearerAuth) BearerAuth(ctx context.Context, operationName gen.OperationName, client *Client) (gen.BearerAuth, error) {
	return gen.BearerAuth{}, nil
}
