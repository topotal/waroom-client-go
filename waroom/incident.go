package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

type IncidentService struct {
	client *Client
}

func NewIncidentService(client *Client) *IncidentService {
	return &IncidentService{client: client}
}

func (s *IncidentService) ListIncidents(ctx context.Context, params ListIncidentsParams) ([]IncidentListItem, error) {
	res, err := s.client.raw.ListIncidents(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.ListIncidentsOK:
		return res.Incidents, nil
	case *gen.ListIncidentsUnauthorized:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

// GET /incidents/{uuid}
func (s *IncidentService) GetIncident(ctx context.Context, uuid uuid.UUID) (*Incident, error) {
	params := gen.GetIncidentParams{UUID: uuid}

	res, err := s.client.raw.GetIncident(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.GetIncidentOK:
		return res, nil
	case *gen.GetIncidentUnauthorized:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}
