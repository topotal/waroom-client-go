package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

// GET /incidents/{uuid}
func (c *Client) GetIncident(ctx context.Context, uuid uuid.UUID) (*gen.GetIncidentRes, error) {
	params := gen.GetIncidentParams{UUID: uuid}

	res, err := c.raw.GetIncident(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	return &res, nil
}

// GET /incidents
func (c *Client) ListIncidents(ctx context.Context) (*gen.ListIncidentsRes, error) {
	params := gen.ListIncidentsParams{}

	res, err := c.raw.ListIncidents(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	return &res, nil
}
