package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

type IncidentService struct {
	client *WaroomClient
}

func NewIncidentService(client *WaroomClient) *IncidentService {
	return &IncidentService{client: client}
}

func (s *IncidentService) ListIncidents(ctx context.Context, params gen.ListIncidentsParams) ([]gen.Incident, error) {
	res, err := s.client.raw.ListIncidents(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.IncidentListResponse:
		return res.Incidents, nil
	case *gen.UnauthorizedError:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

// GET /incidents/{uuid}
func (s *IncidentService) GetIncident(ctx context.Context, uuid uuid.UUID) (*gen.GetIncidentRes, error) {
	params := gen.GetIncidentParams{UUID: uuid}

	res, err := s.client.raw.GetIncident(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.GetIncidentRes:
		return res.Incidents, nil
	case *gen.UnauthorizedError:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

	return &res, nil
}
