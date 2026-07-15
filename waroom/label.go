package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

type LabelService struct {
	client *Client
}

func NewLabelService(client *Client) *LabelService {
	return &LabelService{client: client}
}

func (s *LabelService) ListServiceLabels(ctx context.Context, params ListServiceLabelsParams) ([]Label, error) {
	res, err := s.client.raw.ListServiceLabels(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list service labels: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListServiceLabelsOK:
		return res.Labels, nil
	case *gen.ListServiceLabelsUnauthorized:
		return nil, fmt.Errorf("waroom: list service labels: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListServiceLabelsNotFound:
		return nil, fmt.Errorf("waroom: list service labels: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.ListServiceLabelsPaymentRequired:
		return nil, fmt.Errorf("waroom: list service labels: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list service labels: unhandled response %T", res)
	}
}

func (s *LabelService) CreateServiceLabel(ctx context.Context, req CreateServiceLabelReq, params CreateServiceLabelParams) (*CreatedLabel, error) {
	res, err := s.client.raw.CreateServiceLabel(ctx, &gen.CreateServiceLabelReq{Label: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: create service label: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateServiceLabelCreated:
		return &res.Label, nil
	case *gen.CreateServiceLabelUnauthorized:
		return nil, fmt.Errorf("waroom: create service label: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateServiceLabelNotFound:
		return nil, fmt.Errorf("waroom: create service label: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.CreateServiceLabelUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create service label: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	case *gen.CreateServiceLabelPaymentRequired:
		return nil, fmt.Errorf("waroom: create service label: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: create service label: unhandled response %T", res)
	}
}

func (s *LabelService) UpdateServiceLabel(ctx context.Context, req UpdateServiceLabelReq, params UpdateServiceLabelParams) (*UpdatedLabel, error) {
	res, err := s.client.raw.UpdateServiceLabel(ctx, &gen.UpdateServiceLabelReq{Label: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update service label: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateServiceLabelOK:
		return &res.Label, nil
	case *gen.UpdateServiceLabelUnauthorized:
		return nil, fmt.Errorf("waroom: update service label: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateServiceLabelNotFound:
		return nil, fmt.Errorf("waroom: update service label: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateServiceLabelUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update service label: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	case *gen.UpdateServiceLabelPaymentRequired:
		return nil, fmt.Errorf("waroom: update service label: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: update service label: unhandled response %T", res)
	}
}

func (s *LabelService) DeleteServiceLabel(ctx context.Context, params DeleteServiceLabelParams) error {
	res, err := s.client.raw.DeleteServiceLabel(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete service label: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteServiceLabelNoContent:
		return nil
	case *gen.DeleteServiceLabelUnauthorized:
		return fmt.Errorf("waroom: delete service label: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteServiceLabelNotFound:
		return fmt.Errorf("waroom: delete service label: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteServiceLabelPaymentRequired:
		return fmt.Errorf("waroom: delete service label: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete service label: unhandled response %T", res)
	}
}
