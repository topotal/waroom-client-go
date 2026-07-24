package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// WebhookConfigurationService handles communication with the webhook configuration related
// methods of the Waroom API. Use it via Client.WebhookConfigurations.
type WebhookConfigurationService service

// NewWebhookConfigurationService returns client.WebhookConfigurations.
//
// Deprecated: Use Client.WebhookConfigurations instead.
func NewWebhookConfigurationService(client *Client) *WebhookConfigurationService {
	return client.WebhookConfigurations
}

func (s *WebhookConfigurationService) List(ctx context.Context, params ListServiceWebhookConfigurationsParams) ([]WebhookConfigurationListItem, error) {
	res, err := s.client.raw.ListServiceWebhookConfigurations(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list service webhook configurations: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListServiceWebhookConfigurationsOK:
		return res.WebhookConfigurations, nil
	case *gen.ListServiceWebhookConfigurationsNotFound:
		return nil, fmt.Errorf("waroom: list service webhook configurations: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.ListServiceWebhookConfigurationsUnauthorized:
		return nil, fmt.Errorf("waroom: list service webhook configurations: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListServiceWebhookConfigurationsPaymentRequired:
		return nil, fmt.Errorf("waroom: list service webhook configurations: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list service webhook configurations: unhandled response %T", res)
	}
}

func (s *WebhookConfigurationService) Get(ctx context.Context, params GetServiceWebhookConfigurationParams) (*WebhookConfiguration, error) {
	res, err := s.client.raw.GetServiceWebhookConfiguration(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get service webhook configuration: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetServiceWebhookConfigurationOK:
		return &res.WebhookConfiguration, nil
	case *gen.GetServiceWebhookConfigurationNotFound:
		return nil, fmt.Errorf("waroom: get service webhook configuration: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetServiceWebhookConfigurationUnauthorized:
		return nil, fmt.Errorf("waroom: get service webhook configuration: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetServiceWebhookConfigurationPaymentRequired:
		return nil, fmt.Errorf("waroom: get service webhook configuration: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get service webhook configuration: unhandled response %T", res)
	}
}
