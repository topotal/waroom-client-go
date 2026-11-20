package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// ActionItemService handles communication with the action item related
// methods of the Waroom API. Use it via Client.ActionItems.
type ActionItemService service

// NewActionItemService returns client.ActionItems.
//
// Deprecated: Use Client.ActionItems instead.
func NewActionItemService(client *Client) *ActionItemService {
	return client.ActionItems
}

func (s *ActionItemService) List(ctx context.Context, params ListIncidentActionItemsParams) ([]ActionItemListItem, error) {
	res, err := s.client.raw.ListIncidentActionItems(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list incident action items: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListIncidentActionItemsOK:
		return res.ActionItems, nil
	case *gen.ListIncidentActionItemsNotFound:
		return nil, fmt.Errorf("waroom: list incident action items: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.ListIncidentActionItemsUnauthorized:
		return nil, fmt.Errorf("waroom: list incident action items: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListIncidentActionItemsPaymentRequired:
		return nil, fmt.Errorf("waroom: list incident action items: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list incident action items: unhandled response %T", res)
	}
}

func (s *ActionItemService) Create(ctx context.Context, req CreateActionItemReq, params CreateIncidentActionItemParams) (*CreatedActionItem, error) {
	res, err := s.client.raw.CreateIncidentActionItem(ctx, &gen.CreateIncidentActionItemReq{ActionItem: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: create incident action item: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateIncidentActionItemCreated:
		return &res.ActionItem, nil
	case *gen.CreateIncidentActionItemNotFound:
		return nil, fmt.Errorf("waroom: create incident action item: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.CreateIncidentActionItemUnauthorized:
		return nil, fmt.Errorf("waroom: create incident action item: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateIncidentActionItemPaymentRequired:
		return nil, fmt.Errorf("waroom: create incident action item: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreateIncidentActionItemUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create incident action item: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: create incident action item: unhandled response %T", res)
	}
}

func (s *ActionItemService) Get(ctx context.Context, params GetIncidentActionItemParams) (*ActionItem, error) {
	res, err := s.client.raw.GetIncidentActionItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get incident action item: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetIncidentActionItemOK:
		return &res.ActionItem, nil
	case *gen.GetIncidentActionItemNotFound:
		return nil, fmt.Errorf("waroom: get incident action item: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetIncidentActionItemUnauthorized:
		return nil, fmt.Errorf("waroom: get incident action item: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetIncidentActionItemPaymentRequired:
		return nil, fmt.Errorf("waroom: get incident action item: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get incident action item: unhandled response %T", res)
	}
}

func (s *ActionItemService) Update(ctx context.Context, req UpdateActionItemReq, params UpdateIncidentActionItemParams) (*UpdatedActionItem, error) {
	res, err := s.client.raw.UpdateIncidentActionItem(ctx, &gen.UpdateIncidentActionItemReq{ActionItem: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update incident action item: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateIncidentActionItemOK:
		return &res.ActionItem, nil
	case *gen.UpdateIncidentActionItemNotFound:
		return nil, fmt.Errorf("waroom: update incident action item: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateIncidentActionItemUnauthorized:
		return nil, fmt.Errorf("waroom: update incident action item: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateIncidentActionItemPaymentRequired:
		return nil, fmt.Errorf("waroom: update incident action item: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateIncidentActionItemUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update incident action item: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update incident action item: unhandled response %T", res)
	}
}

func (s *ActionItemService) Delete(ctx context.Context, params DeleteIncidentActionItemParams) error {
	res, err := s.client.raw.DeleteIncidentActionItem(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete incident action item: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteIncidentActionItemNoContent:
		return nil
	case *gen.DeleteIncidentActionItemNotFound:
		return fmt.Errorf("waroom: delete incident action item: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteIncidentActionItemUnauthorized:
		return fmt.Errorf("waroom: delete incident action item: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteIncidentActionItemPaymentRequired:
		return fmt.Errorf("waroom: delete incident action item: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete incident action item: unhandled response %T", res)
	}
}
