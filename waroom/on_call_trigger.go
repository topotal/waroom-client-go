package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// OnCallTriggerService handles communication with the on-call trigger related
// methods of the Waroom API. Use it via Client.OnCallTriggers.
type OnCallTriggerService service

func (s *OnCallTriggerService) List(ctx context.Context, params ListOnCallTriggersParams) ([]OnCallTriggerListItem, error) {
	res, err := s.client.raw.ListOnCallTriggers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list on-call triggers: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListOnCallTriggersOK:
		return res.Triggers, nil
	case *gen.ListOnCallTriggersUnauthorized:
		return nil, fmt.Errorf("waroom: list on-call triggers: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListOnCallTriggersPaymentRequired:
		return nil, fmt.Errorf("waroom: list on-call triggers: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list on-call triggers: unhandled response %T", res)
	}
}

func (s *OnCallTriggerService) Get(ctx context.Context, params GetOnCallTriggerParams) (*OnCallTrigger, error) {
	res, err := s.client.raw.GetOnCallTrigger(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get on-call trigger: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetOnCallTriggerOK:
		return &res.Trigger, nil
	case *gen.GetOnCallTriggerNotFound:
		return nil, fmt.Errorf("waroom: get on-call trigger: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetOnCallTriggerUnauthorized:
		return nil, fmt.Errorf("waroom: get on-call trigger: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetOnCallTriggerPaymentRequired:
		return nil, fmt.Errorf("waroom: get on-call trigger: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get on-call trigger: unhandled response %T", res)
	}
}

func (s *OnCallTriggerService) Create(ctx context.Context, req CreateOnCallTriggerReq) (*CreatedOnCallTrigger, error) {
	res, err := s.client.raw.CreateOnCallTrigger(ctx, &gen.CreateOnCallTriggerReq{Trigger: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: create on-call trigger: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateOnCallTriggerCreated:
		return &res.Trigger, nil
	case *gen.CreateOnCallTriggerBadRequest:
		return nil, fmt.Errorf("waroom: create on-call trigger: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.CreateOnCallTriggerUnauthorized:
		return nil, fmt.Errorf("waroom: create on-call trigger: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateOnCallTriggerPaymentRequired:
		return nil, fmt.Errorf("waroom: create on-call trigger: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreateOnCallTriggerUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create on-call trigger: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: create on-call trigger: unhandled response %T", res)
	}
}

func (s *OnCallTriggerService) Update(ctx context.Context, req UpdateOnCallTriggerReq, params UpdateOnCallTriggerParams) (*UpdatedOnCallTrigger, error) {
	res, err := s.client.raw.UpdateOnCallTrigger(ctx, &gen.UpdateOnCallTriggerReq{Trigger: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update on-call trigger: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateOnCallTriggerOK:
		return &res.Trigger, nil
	case *gen.UpdateOnCallTriggerBadRequest:
		return nil, fmt.Errorf("waroom: update on-call trigger: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallTriggerNotFound:
		return nil, fmt.Errorf("waroom: update on-call trigger: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallTriggerUnauthorized:
		return nil, fmt.Errorf("waroom: update on-call trigger: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateOnCallTriggerPaymentRequired:
		return nil, fmt.Errorf("waroom: update on-call trigger: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateOnCallTriggerUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update on-call trigger: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: update on-call trigger: unhandled response %T", res)
	}
}

func (s *OnCallTriggerService) Delete(ctx context.Context, params DeleteOnCallTriggerParams) error {
	res, err := s.client.raw.DeleteOnCallTrigger(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete on-call trigger: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteOnCallTriggerNoContent:
		return nil
	case *gen.DeleteOnCallTriggerNotFound:
		return fmt.Errorf("waroom: delete on-call trigger: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteOnCallTriggerUnauthorized:
		return fmt.Errorf("waroom: delete on-call trigger: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteOnCallTriggerPaymentRequired:
		return fmt.Errorf("waroom: delete on-call trigger: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete on-call trigger: unhandled response %T", res)
	}
}
