package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/google/uuid"
)

// IncidentService handles communication with the incident related
// methods of the Waroom API. Use it via Client.Incidents.
type IncidentService service

// NewIncidentService returns client.Incidents.
//
// Deprecated: Use Client.Incidents instead.
func NewIncidentService(client *Client) *IncidentService {
	return client.Incidents
}

func (s *IncidentService) List(ctx context.Context, params ListIncidentsParams) ([]IncidentListItem, error) {
	res, err := s.client.raw.ListIncidents(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list incidents: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListIncidentsOK:
		return res.Incidents, nil
	case *gen.ListIncidentsUnauthorized:
		return nil, fmt.Errorf("waroom: list incidents: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListIncidentsPaymentRequired:
		return nil, fmt.Errorf("waroom: list incidents: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list incidents: unhandled response %T", res)
	}
}

func (s *IncidentService) Get(ctx context.Context, uuid uuid.UUID) (*Incident, error) {
	params := gen.GetIncidentParams{UUID: uuid}

	res, err := s.client.raw.GetIncident(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get incident: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetIncidentOK:
		return res, nil
	case *gen.GetIncidentUnauthorized:
		return nil, fmt.Errorf("waroom: get incident: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetIncidentNotFound:
		return nil, fmt.Errorf("waroom: get incident: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetIncidentPaymentRequired:
		return nil, fmt.Errorf("waroom: get incident: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get incident: unhandled response %T", res)
	}
}
