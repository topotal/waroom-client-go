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

func (auth *WaroomBearerAuth) BearerAuth(ctx context.Context, operationName gen.OperationName, client *gen.Client) (gen.BearerAuth, error) {
	println("run auth")
	return gen.BearerAuth{Token: os.Getenv(WAROOM_API_KEY)}, nil

	// curl -X GET "https://api.staging.waroom.com/api/v0/incidents?sort_order=asc" -H "Authorization: Bearer " | jq

}
