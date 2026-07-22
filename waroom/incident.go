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

// ListActionItems is a convenience wrapper for Client.ActionItems.List.
func (s *IncidentService) ListActionItems(ctx context.Context, params ListIncidentActionItemsParams) ([]ActionItemListItem, error) {
	return s.client.ActionItems.List(ctx, params)
}

// CreateActionItem is a convenience wrapper for Client.ActionItems.Create.
func (s *IncidentService) CreateActionItem(ctx context.Context, req CreateActionItemReq, params CreateIncidentActionItemParams) (*CreatedActionItem, error) {
	return s.client.ActionItems.Create(ctx, req, params)
}

// GetActionItem is a convenience wrapper for Client.ActionItems.Get.
func (s *IncidentService) GetActionItem(ctx context.Context, params GetIncidentActionItemParams) (*ActionItem, error) {
	return s.client.ActionItems.Get(ctx, params)
}

// DeleteActionItem is a convenience wrapper for Client.ActionItems.Delete.
func (s *IncidentService) DeleteActionItem(ctx context.Context, params DeleteIncidentActionItemParams) error {
	return s.client.ActionItems.Delete(ctx, params)
}

// UpdateService is a convenience wrapper for Client.Services.UpdateIncidentService.
func (s *IncidentService) UpdateService(ctx context.Context, req UpdateIncidentServiceReq, params UpdateIncidentServiceParams) (*IncidentServiceAssignment, error) {
	return s.client.Services.UpdateIncidentService(ctx, req, params)
}

func (s *IncidentService) Create(ctx context.Context, req CreateIncidentReq) (*CreatedIncident, error) {
	res, err := s.client.raw.CreateIncident(ctx, &gen.CreateIncidentReq{Incident: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: create incident: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateIncidentCreated:
		return res, nil
	case *gen.CreateIncidentNotFound:
		return nil, fmt.Errorf("waroom: create incident: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.CreateIncidentUnauthorized:
		return nil, fmt.Errorf("waroom: create incident: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateIncidentPaymentRequired:
		return nil, fmt.Errorf("waroom: create incident: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreateIncidentUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create incident: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: create incident: unhandled response %T", res)
	}
}

func (s *IncidentService) Update(ctx context.Context, req UpdateIncidentReq, params UpdateIncidentParams) (*UpdatedIncident, error) {
	res, err := s.client.raw.UpdateIncident(ctx, &gen.UpdateIncidentReq{Incident: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentOK:
		return res, nil
	case *gen.UpdateIncidentNotFound:
		return nil, fmt.Errorf("waroom: update incident: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentUnauthorized:
		return nil, fmt.Errorf("waroom: update incident: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentPaymentRequired:
		return nil, fmt.Errorf("waroom: update incident: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateIncidentUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update incident: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update incident: unhandled response %T", res)
	}
}

func (s *IncidentService) Delete(ctx context.Context, params DeleteIncidentParams) error {
	res, err := s.client.raw.DeleteIncident(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete incident: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteIncidentNoContent:
		return nil
	case *gen.DeleteIncidentNotFound:
		return fmt.Errorf("waroom: delete incident: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteIncidentUnauthorized:
		return fmt.Errorf("waroom: delete incident: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteIncidentPaymentRequired:
		return fmt.Errorf("waroom: delete incident: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete incident: unhandled response %T", res)
	}
}

func (s *IncidentService) UpdateStatus(ctx context.Context, req UpdateIncidentStatusReq, params UpdateIncidentStatusParams) (*IncidentStatus, error) {
	res, err := s.client.raw.UpdateIncidentStatus(ctx, &gen.UpdateIncidentStatusReq{Status: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident status: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentStatusOK:
		return res, nil
	case *gen.UpdateIncidentStatusBadRequest:
		return nil, fmt.Errorf("waroom: update incident status: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentStatusNotFound:
		return nil, fmt.Errorf("waroom: update incident status: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentStatusUnauthorized:
		return nil, fmt.Errorf("waroom: update incident status: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentStatusPaymentRequired:
		return nil, fmt.Errorf("waroom: update incident status: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: update incident status: unhandled response %T", res)
	}
}

func (s *IncidentService) UpdateSeverity(ctx context.Context, req UpdateIncidentSeverityReq, params UpdateIncidentSeverityParams) (*IncidentSeverity, error) {
	res, err := s.client.raw.UpdateIncidentSeverity(ctx, &gen.UpdateIncidentSeverityReq{Severity: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident severity: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentSeverityOK:
		return res, nil
	case *gen.UpdateIncidentSeverityBadRequest:
		return nil, fmt.Errorf("waroom: update incident severity: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentSeverityNotFound:
		return nil, fmt.Errorf("waroom: update incident severity: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentSeverityUnauthorized:
		return nil, fmt.Errorf("waroom: update incident severity: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentSeverityPaymentRequired:
		return nil, fmt.Errorf("waroom: update incident severity: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: update incident severity: unhandled response %T", res)
	}
}

func (s *IncidentService) UpdateCommander(ctx context.Context, req UpdateIncidentCommanderReq, params UpdateIncidentCommanderParams) (*IncidentCommander, error) {
	res, err := s.client.raw.UpdateIncidentCommander(ctx, &gen.UpdateIncidentCommanderReq{Incident: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident commander: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentCommanderOK:
		return res, nil
	case *gen.UpdateIncidentCommanderBadRequest:
		return nil, fmt.Errorf("waroom: update incident commander: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentCommanderNotFound:
		return nil, fmt.Errorf("waroom: update incident commander: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentCommanderUnauthorized:
		return nil, fmt.Errorf("waroom: update incident commander: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentCommanderPaymentRequired:
		return nil, fmt.Errorf("waroom: update incident commander: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateIncidentCommanderUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update incident commander: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update incident commander: unhandled response %T", res)
	}
}

func (s *IncidentService) UpdateResponders(ctx context.Context, req UpdateIncidentRespondersReq, params UpdateIncidentRespondersParams) (*IncidentResponders, error) {
	res, err := s.client.raw.UpdateIncidentResponders(ctx, &gen.UpdateIncidentRespondersReq{Incident: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident responders: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentRespondersOK:
		return res, nil
	case *gen.UpdateIncidentRespondersBadRequest:
		return nil, fmt.Errorf("waroom: update incident responders: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentRespondersNotFound:
		return nil, fmt.Errorf("waroom: update incident responders: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentRespondersUnauthorized:
		return nil, fmt.Errorf("waroom: update incident responders: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentRespondersPaymentRequired:
		return nil, fmt.Errorf("waroom: update incident responders: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateIncidentRespondersUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update incident responders: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update incident responders: unhandled response %T", res)
	}
}

func (s *IncidentService) CreateMetrics(ctx context.Context, req CreateIncidentMetricsReq, params CreateIncidentMetricsParams) (*IncidentMetricsActivity, error) {
	res, err := s.client.raw.CreateIncidentMetrics(ctx, &req, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: create incident metrics: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateIncidentMetricsOK:
		return res, nil
	case *gen.CreateIncidentMetricsBadRequest:
		return nil, fmt.Errorf("waroom: create incident metrics: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.CreateIncidentMetricsNotFound:
		return nil, fmt.Errorf("waroom: create incident metrics: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.CreateIncidentMetricsUnauthorized:
		return nil, fmt.Errorf("waroom: create incident metrics: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateIncidentMetricsPaymentRequired:
		return nil, fmt.Errorf("waroom: create incident metrics: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: create incident metrics: unhandled response %T", res)
	}
}
