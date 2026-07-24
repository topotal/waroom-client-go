package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// RunbookRuleService handles communication with the runbook rule related
// methods of the Waroom API. Use it via Client.RunbookRules.
type RunbookRuleService service

// NewRunbookRuleService returns client.RunbookRules.
//
// Deprecated: Use Client.RunbookRules instead.
func NewRunbookRuleService(client *Client) *RunbookRuleService {
	return client.RunbookRules
}

func (s *RunbookRuleService) List(ctx context.Context, params ListServiceRunbookRulesParams) ([]RunbookRuleListItem, error) {
	res, err := s.client.raw.ListServiceRunbookRules(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list service runbook rules: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListServiceRunbookRulesOK:
		return res.RunbookRules, nil
	case *gen.ListServiceRunbookRulesUnauthorized:
		return nil, fmt.Errorf("waroom: list service runbook rules: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListServiceRunbookRulesNotFound:
		return nil, fmt.Errorf("waroom: list service runbook rules: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.ListServiceRunbookRulesPaymentRequired:
		return nil, fmt.Errorf("waroom: list service runbook rules: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list service runbook rules: unhandled response %T", res)
	}
}

func (s *RunbookRuleService) Create(ctx context.Context, req CreateRunbookRuleReq, params CreateServiceRunbookRuleParams) (*CreatedRunbookRule, error) {
	res, err := s.client.raw.CreateServiceRunbookRule(ctx, &gen.CreateServiceRunbookRuleReq{RunbookRule: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: create service runbook rule: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateServiceRunbookRuleCreated:
		return res, nil
	case *gen.CreateServiceRunbookRuleUnauthorized:
		return nil, fmt.Errorf("waroom: create service runbook rule: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateServiceRunbookRuleNotFound:
		return nil, fmt.Errorf("waroom: create service runbook rule: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.CreateServiceRunbookRuleUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create service runbook rule: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	case *gen.CreateServiceRunbookRulePaymentRequired:
		return nil, fmt.Errorf("waroom: create service runbook rule: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: create service runbook rule: unhandled response %T", res)
	}
}

func (s *RunbookRuleService) Update(ctx context.Context, req UpdateRunbookRuleReq, params UpdateServiceRunbookRuleParams) (*UpdatedRunbookRule, error) {
	res, err := s.client.raw.UpdateServiceRunbookRule(ctx, &gen.UpdateServiceRunbookRuleReq{RunbookRule: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update service runbook rule: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateServiceRunbookRuleOK:
		return res, nil
	case *gen.UpdateServiceRunbookRuleUnauthorized:
		return nil, fmt.Errorf("waroom: update service runbook rule: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateServiceRunbookRuleNotFound:
		return nil, fmt.Errorf("waroom: update service runbook rule: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateServiceRunbookRuleUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update service runbook rule: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	case *gen.UpdateServiceRunbookRulePaymentRequired:
		return nil, fmt.Errorf("waroom: update service runbook rule: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: update service runbook rule: unhandled response %T", res)
	}
}

func (s *RunbookRuleService) Delete(ctx context.Context, params DeleteServiceRunbookRuleParams) error {
	res, err := s.client.raw.DeleteServiceRunbookRule(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete service runbook rule: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteServiceRunbookRuleNoContent:
		return nil
	case *gen.DeleteServiceRunbookRuleUnauthorized:
		return fmt.Errorf("waroom: delete service runbook rule: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteServiceRunbookRuleNotFound:
		return fmt.Errorf("waroom: delete service runbook rule: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteServiceRunbookRulePaymentRequired:
		return fmt.Errorf("waroom: delete service runbook rule: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete service runbook rule: unhandled response %T", res)
	}
}
