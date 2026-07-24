package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// RunbookService handles communication with the runbook related
// methods of the Waroom API. Use it via Client.Runbooks.
type RunbookService service

// NewRunbookService returns client.Runbooks.
//
// Deprecated: Use Client.Runbooks instead.
func NewRunbookService(client *Client) *RunbookService {
	return client.Runbooks
}

func (s *RunbookService) List(ctx context.Context, params ListRunbooksParams) ([]RunbookListItem, error) {
	res, err := s.client.raw.ListRunbooks(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list runbooks: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListRunbooksOK:
		return res.Runbooks, nil
	case *gen.ListRunbooksUnauthorized:
		return nil, fmt.Errorf("waroom: list runbooks: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListRunbooksPaymentRequired:
		return nil, fmt.Errorf("waroom: list runbooks: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list runbooks: unhandled response %T", res)
	}
}

func (s *RunbookService) Get(ctx context.Context, params GetRunbookParams) (*Runbook, error) {
	res, err := s.client.raw.GetRunbook(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get runbook: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetRunbookOK:
		return &res.Runbook, nil
	case *gen.GetRunbookNotFound:
		return nil, fmt.Errorf("waroom: get runbook: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetRunbookUnauthorized:
		return nil, fmt.Errorf("waroom: get runbook: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetRunbookPaymentRequired:
		return nil, fmt.Errorf("waroom: get runbook: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get runbook: unhandled response %T", res)
	}
}

func (s *RunbookService) Create(ctx context.Context, req CreateRunbookReq) (*CreatedRunbook, error) {
	res, err := s.client.raw.CreateRunbook(ctx, &gen.CreateRunbookReq{Runbook: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: create runbook: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateRunbookCreated:
		return &res.Runbook, nil
	case *gen.CreateRunbookUnauthorized:
		return nil, fmt.Errorf("waroom: create runbook: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateRunbookPaymentRequired:
		return nil, fmt.Errorf("waroom: create runbook: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreateRunbookUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create runbook: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: create runbook: unhandled response %T", res)
	}
}

func (s *RunbookService) Update(ctx context.Context, req UpdateRunbookReq, params UpdateRunbookParams) (*UpdatedRunbook, error) {
	res, err := s.client.raw.UpdateRunbook(ctx, &gen.UpdateRunbookReq{Runbook: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update runbook: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateRunbookOK:
		return &res.Runbook, nil
	case *gen.UpdateRunbookNotFound:
		return nil, fmt.Errorf("waroom: update runbook: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateRunbookUnauthorized:
		return nil, fmt.Errorf("waroom: update runbook: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateRunbookPaymentRequired:
		return nil, fmt.Errorf("waroom: update runbook: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateRunbookUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update runbook: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update runbook: unhandled response %T", res)
	}
}

func (s *RunbookService) Delete(ctx context.Context, params DeleteRunbookParams) error {
	res, err := s.client.raw.DeleteRunbook(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete runbook: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteRunbookNoContent:
		return nil
	case *gen.DeleteRunbookNotFound:
		return fmt.Errorf("waroom: delete runbook: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteRunbookUnauthorized:
		return fmt.Errorf("waroom: delete runbook: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteRunbookPaymentRequired:
		return fmt.Errorf("waroom: delete runbook: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.DeleteRunbookUnprocessableEntity:
		return fmt.Errorf("waroom: delete runbook: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return fmt.Errorf("waroom: delete runbook: unhandled response %T", res)
	}
}
