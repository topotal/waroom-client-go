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
