package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// ServiceService handles communication with the service related
// methods of the Waroom API. Use it via Client.Services.
type ServiceService service

// NewServiceService returns client.Services.
//
// Deprecated: Use Client.Services instead.
func NewServiceService(client *Client) *ServiceService {
	return client.Services
}

func (s *ServiceService) List(ctx context.Context, params ListServicesParams) ([]Service, error) {
	res, err := s.client.raw.ListServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("waroom: list services: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListServicesOK:
		return res.Services, nil
	case *gen.ListServicesUnauthorized:
		return nil, fmt.Errorf("waroom: list services: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListServicesPaymentRequired:
		return nil, fmt.Errorf("waroom: list services: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list services: unhandled response %T", res)
	}
}

func (s *ServiceService) UpdateIncidentService(ctx context.Context, req UpdateIncidentServiceReq, params UpdateIncidentServiceParams) (*IncidentServiceAssignment, error) {
	res, err := s.client.raw.UpdateIncidentService(ctx, &gen.UpdateIncidentServiceReq{Incident: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident service: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentServiceOK:
		return res, nil
	case *gen.UpdateIncidentServiceBadRequest:
		return nil, fmt.Errorf("waroom: update incident service: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentServiceNotFound:
		return nil, fmt.Errorf("waroom: update incident service: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentServiceUnauthorized:
		return nil, fmt.Errorf("waroom: update incident service: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentServicePaymentRequired:
		return nil, fmt.Errorf("waroom: update incident service: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateIncidentServiceUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update incident service: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update incident service: unhandled response %T", res)
	}
}

func (s *ServiceService) GetArchitectureContext(ctx context.Context, params GetServiceArchitectureContextParams) (*ServiceArchitectureContext, error) {
	res, err := s.client.raw.GetServiceArchitectureContext(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get service architecture context: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetServiceArchitectureContextOK:
		return res, nil
	case *gen.GetServiceArchitectureContextNotFound:
		return nil, fmt.Errorf("waroom: get service architecture context: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetServiceArchitectureContextUnauthorized:
		return nil, fmt.Errorf("waroom: get service architecture context: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetServiceArchitectureContextPaymentRequired:
		return nil, fmt.Errorf("waroom: get service architecture context: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get service architecture context: unhandled response %T", res)
	}
}

func (s *ServiceService) UpdateArchitectureContext(ctx context.Context, req UpdateArchitectureContextReq, params UpdateServiceArchitectureContextParams) (*UpdatedServiceArchitectureContext, error) {
	res, err := s.client.raw.UpdateServiceArchitectureContext(ctx, &gen.UpdateServiceArchitectureContextReq{ServiceArchitectureContext: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update service architecture context: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateServiceArchitectureContextOK:
		return res, nil
	case *gen.UpdateServiceArchitectureContextNotFound:
		return nil, fmt.Errorf("waroom: update service architecture context: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateServiceArchitectureContextUnauthorized:
		return nil, fmt.Errorf("waroom: update service architecture context: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateServiceArchitectureContextPaymentRequired:
		return nil, fmt.Errorf("waroom: update service architecture context: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: update service architecture context: unhandled response %T", res)
	}
}

func (s *ServiceService) GetNotificationSlackChannel(ctx context.Context, params GetServiceNotificationSlackChannelParams) (*ServiceNotificationSlackChannel, error) {
	res, err := s.client.raw.GetServiceNotificationSlackChannel(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get service notification slack channel: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetServiceNotificationSlackChannelOK:
		return res, nil
	case *gen.GetServiceNotificationSlackChannelNotFound:
		return nil, fmt.Errorf("waroom: get service notification slack channel: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetServiceNotificationSlackChannelUnauthorized:
		return nil, fmt.Errorf("waroom: get service notification slack channel: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetServiceNotificationSlackChannelPaymentRequired:
		return nil, fmt.Errorf("waroom: get service notification slack channel: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get service notification slack channel: unhandled response %T", res)
	}
}

func (s *ServiceService) PutNotificationSlackChannel(ctx context.Context, req PutNotificationSlackChannelReq, params PutServiceNotificationSlackChannelParams) (*UpdatedServiceNotificationSlackChannel, error) {
	res, err := s.client.raw.PutServiceNotificationSlackChannel(ctx, &gen.PutServiceNotificationSlackChannelReq{ServiceNotificationSlackChannel: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: put service notification slack channel: %w", err)
	}

	switch res := res.(type) {
	case *gen.PutServiceNotificationSlackChannelOK:
		return res, nil
	case *gen.PutServiceNotificationSlackChannelNotFound:
		return nil, fmt.Errorf("waroom: put service notification slack channel: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.PutServiceNotificationSlackChannelUnauthorized:
		return nil, fmt.Errorf("waroom: put service notification slack channel: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.PutServiceNotificationSlackChannelPaymentRequired:
		return nil, fmt.Errorf("waroom: put service notification slack channel: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.PutServiceNotificationSlackChannelUnprocessableEntity:
		return nil, fmt.Errorf("waroom: put service notification slack channel: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: put service notification slack channel: unhandled response %T", res)
	}
}

func (s *ServiceService) DeleteNotificationSlackChannel(ctx context.Context, params DeleteServiceNotificationSlackChannelParams) error {
	res, err := s.client.raw.DeleteServiceNotificationSlackChannel(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete service notification slack channel: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteServiceNotificationSlackChannelNoContent:
		return nil
	case *gen.DeleteServiceNotificationSlackChannelNotFound:
		return fmt.Errorf("waroom: delete service notification slack channel: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteServiceNotificationSlackChannelUnauthorized:
		return fmt.Errorf("waroom: delete service notification slack channel: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteServiceNotificationSlackChannelPaymentRequired:
		return fmt.Errorf("waroom: delete service notification slack channel: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete service notification slack channel: unhandled response %T", res)
	}
}

func (s *ServiceService) GetIncidentDocumentTemplate(ctx context.Context, params GetIncidentDocumentTemplateParams) (*IncidentDocumentTemplate, error) {
	res, err := s.client.raw.GetIncidentDocumentTemplate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get incident document template: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetIncidentDocumentTemplateOK:
		return res, nil
	case *gen.GetIncidentDocumentTemplateNotFound:
		return nil, fmt.Errorf("waroom: get incident document template: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetIncidentDocumentTemplateUnauthorized:
		return nil, fmt.Errorf("waroom: get incident document template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetIncidentDocumentTemplatePaymentRequired:
		return nil, fmt.Errorf("waroom: get incident document template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get incident document template: unhandled response %T", res)
	}
}

func (s *ServiceService) UpdateIncidentDocumentTemplate(ctx context.Context, req UpdateIncidentDocumentTemplateReq, params UpdateIncidentDocumentTemplateParams) (*UpdatedIncidentDocumentTemplate, error) {
	res, err := s.client.raw.UpdateIncidentDocumentTemplate(ctx, &gen.UpdateIncidentDocumentTemplateReq{IncidentDocumentTemplate: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident document template: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentDocumentTemplateOK:
		return res, nil
	case *gen.UpdateIncidentDocumentTemplateNotFound:
		return nil, fmt.Errorf("waroom: update incident document template: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentDocumentTemplateUnauthorized:
		return nil, fmt.Errorf("waroom: update incident document template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentDocumentTemplatePaymentRequired:
		return nil, fmt.Errorf("waroom: update incident document template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateIncidentDocumentTemplateUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update incident document template: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update incident document template: unhandled response %T", res)
	}
}

func (s *ServiceService) DeleteIncidentDocumentTemplate(ctx context.Context, params DeleteIncidentDocumentTemplateParams) error {
	res, err := s.client.raw.DeleteIncidentDocumentTemplate(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete incident document template: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteIncidentDocumentTemplateNoContent:
		return nil
	case *gen.DeleteIncidentDocumentTemplateNotFound:
		return fmt.Errorf("waroom: delete incident document template: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteIncidentDocumentTemplateUnauthorized:
		return fmt.Errorf("waroom: delete incident document template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteIncidentDocumentTemplatePaymentRequired:
		return fmt.Errorf("waroom: delete incident document template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete incident document template: unhandled response %T", res)
	}
}
