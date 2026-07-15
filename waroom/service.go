package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

type ServiceService struct {
	client *Client
}

func NewServiceService(client *Client) *ServiceService {
	return &ServiceService{client: client}
}

func (s *ServiceService) ListServices(ctx context.Context, params ListServicesParams) ([]Service, error) {
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
