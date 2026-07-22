package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// IncidentDescriptionTemplateService handles communication with the incident description template
// related methods of the Waroom API. Use it via Client.IncidentDescriptionTemplate.
type IncidentDescriptionTemplateService service

// NewIncidentDescriptionTemplateService returns client.IncidentDescriptionTemplate.
//
// Deprecated: Use Client.IncidentDescriptionTemplate instead.
func NewIncidentDescriptionTemplateService(client *Client) *IncidentDescriptionTemplateService {
	return client.IncidentDescriptionTemplate
}

func (s *IncidentDescriptionTemplateService) Get(ctx context.Context) (*IncidentDescriptionTemplate, error) {
	res, err := s.client.raw.GetIncidentDescriptionTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("waroom: get incident description template: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetIncidentDescriptionTemplateOK:
		return res, nil
	case *gen.GetIncidentDescriptionTemplateUnauthorized:
		return nil, fmt.Errorf("waroom: get incident description template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetIncidentDescriptionTemplatePaymentRequired:
		return nil, fmt.Errorf("waroom: get incident description template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get incident description template: unhandled response %T", res)
	}
}

func (s *IncidentDescriptionTemplateService) Update(ctx context.Context, req UpdateIncidentDescriptionTemplateReq) (*UpdatedIncidentDescriptionTemplate, error) {
	res, err := s.client.raw.UpdateIncidentDescriptionTemplate(ctx, &gen.UpdateIncidentDescriptionTemplateReq{IncidentDescriptionTemplate: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident description template: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentDescriptionTemplateOK:
		return res, nil
	case *gen.UpdateIncidentDescriptionTemplateUnauthorized:
		return nil, fmt.Errorf("waroom: update incident description template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentDescriptionTemplatePaymentRequired:
		return nil, fmt.Errorf("waroom: update incident description template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateIncidentDescriptionTemplateUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update incident description template: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update incident description template: unhandled response %T", res)
	}
}

func (s *IncidentDescriptionTemplateService) Delete(ctx context.Context) error {
	res, err := s.client.raw.DeleteIncidentDescriptionTemplate(ctx)
	if err != nil {
		return fmt.Errorf("waroom: delete incident description template: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteIncidentDescriptionTemplateNoContent:
		return nil
	case *gen.DeleteIncidentDescriptionTemplateUnauthorized:
		return fmt.Errorf("waroom: delete incident description template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteIncidentDescriptionTemplatePaymentRequired:
		return fmt.Errorf("waroom: delete incident description template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.DeleteIncidentDescriptionTemplateUnprocessableEntity:
		return fmt.Errorf("waroom: delete incident description template: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return fmt.Errorf("waroom: delete incident description template: unhandled response %T", res)
	}
}
