package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// OnCallEscalationPolicyService handles communication with the on-call
// escalation policy related methods of the Waroom API. Use it via
// Client.OnCallEscalationPolicies.
type OnCallEscalationPolicyService service

func (s *OnCallEscalationPolicyService) List(ctx context.Context, params ListOnCallEscalationPoliciesParams) ([]OnCallEscalationPolicyListItem, error) {
	res, err := s.client.raw.ListOnCallEscalationPolicies(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list on-call escalation policies: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListOnCallEscalationPoliciesOK:
		return res.EscalationPolicies, nil
	case *gen.ListOnCallEscalationPoliciesUnauthorized:
		return nil, fmt.Errorf("waroom: list on-call escalation policies: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListOnCallEscalationPoliciesPaymentRequired:
		return nil, fmt.Errorf("waroom: list on-call escalation policies: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list on-call escalation policies: unhandled response %T", res)
	}
}

func (s *OnCallEscalationPolicyService) Get(ctx context.Context, params GetOnCallEscalationPolicyParams) (*OnCallEscalationPolicy, error) {
	res, err := s.client.raw.GetOnCallEscalationPolicy(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get on-call escalation policy: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetOnCallEscalationPolicyOK:
		return &res.EscalationPolicy, nil
	case *gen.GetOnCallEscalationPolicyUnauthorized:
		return nil, fmt.Errorf("waroom: get on-call escalation policy: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetOnCallEscalationPolicyPaymentRequired:
		return nil, fmt.Errorf("waroom: get on-call escalation policy: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.GetOnCallEscalationPolicyNotFound:
		return nil, fmt.Errorf("waroom: get on-call escalation policy: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: get on-call escalation policy: unhandled response %T", res)
	}
}

func (s *OnCallEscalationPolicyService) Create(ctx context.Context, req CreateOnCallEscalationPolicyReq) (*CreatedOnCallEscalationPolicy, error) {
	res, err := s.client.raw.CreateOnCallEscalationPolicy(ctx, &gen.CreateOnCallEscalationPolicyReq{EscalationPolicy: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: create on-call escalation policy: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateOnCallEscalationPolicyCreated:
		return &res.EscalationPolicy, nil
	case *gen.CreateOnCallEscalationPolicyBadRequest:
		return nil, fmt.Errorf("waroom: create on-call escalation policy: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.CreateOnCallEscalationPolicyUnauthorized:
		return nil, fmt.Errorf("waroom: create on-call escalation policy: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateOnCallEscalationPolicyPaymentRequired:
		return nil, fmt.Errorf("waroom: create on-call escalation policy: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreateOnCallEscalationPolicyUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create on-call escalation policy: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: create on-call escalation policy: unhandled response %T", res)
	}
}

func (s *OnCallEscalationPolicyService) Update(ctx context.Context, req UpdateOnCallEscalationPolicyReq, params UpdateOnCallEscalationPolicyParams) (*UpdatedOnCallEscalationPolicy, error) {
	res, err := s.client.raw.UpdateOnCallEscalationPolicy(ctx, &gen.UpdateOnCallEscalationPolicyReq{EscalationPolicy: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update on-call escalation policy: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateOnCallEscalationPolicyOK:
		return &res.EscalationPolicy, nil
	case *gen.UpdateOnCallEscalationPolicyBadRequest:
		return nil, fmt.Errorf("waroom: update on-call escalation policy: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallEscalationPolicyUnauthorized:
		return nil, fmt.Errorf("waroom: update on-call escalation policy: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateOnCallEscalationPolicyPaymentRequired:
		return nil, fmt.Errorf("waroom: update on-call escalation policy: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateOnCallEscalationPolicyNotFound:
		return nil, fmt.Errorf("waroom: update on-call escalation policy: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallEscalationPolicyUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update on-call escalation policy: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: update on-call escalation policy: unhandled response %T", res)
	}
}

func (s *OnCallEscalationPolicyService) Delete(ctx context.Context, params DeleteOnCallEscalationPolicyParams) error {
	res, err := s.client.raw.DeleteOnCallEscalationPolicy(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete on-call escalation policy: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteOnCallEscalationPolicyNoContent:
		return nil
	case *gen.DeleteOnCallEscalationPolicyUnauthorized:
		return fmt.Errorf("waroom: delete on-call escalation policy: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteOnCallEscalationPolicyPaymentRequired:
		return fmt.Errorf("waroom: delete on-call escalation policy: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.DeleteOnCallEscalationPolicyNotFound:
		return fmt.Errorf("waroom: delete on-call escalation policy: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteOnCallEscalationPolicyUnprocessableEntity:
		return fmt.Errorf("waroom: delete on-call escalation policy: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return fmt.Errorf("waroom: delete on-call escalation policy: unhandled response %T", res)
	}
}

// GetSchedule returns the day-by-day on-call schedule of an escalation
// policy. StartDate / EndDate are optional; the API fills in the missing side
// so that the range covers 7 days.
func (s *OnCallEscalationPolicyService) GetSchedule(ctx context.Context, params GetOnCallEscalationPolicyScheduleParams) (*OnCallEscalationPolicySchedule, error) {
	res, err := s.client.raw.GetOnCallEscalationPolicySchedule(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get on-call escalation policy schedule: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetOnCallEscalationPolicyScheduleOK:
		return &res.Schedule, nil
	case *gen.GetOnCallEscalationPolicyScheduleBadRequest:
		return nil, fmt.Errorf("waroom: get on-call escalation policy schedule: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.GetOnCallEscalationPolicyScheduleUnauthorized:
		return nil, fmt.Errorf("waroom: get on-call escalation policy schedule: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetOnCallEscalationPolicySchedulePaymentRequired:
		return nil, fmt.Errorf("waroom: get on-call escalation policy schedule: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.GetOnCallEscalationPolicyScheduleNotFound:
		return nil, fmt.Errorf("waroom: get on-call escalation policy schedule: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: get on-call escalation policy schedule: unhandled response %T", res)
	}
}
