package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// PostmortemService handles communication with the postmortem related
// methods of the Waroom API. Use it via Client.Postmortems.
type PostmortemService service

// NewPostmortemService returns client.Postmortems.
//
// Deprecated: Use Client.Postmortems instead.
func NewPostmortemService(client *Client) *PostmortemService {
	return client.Postmortems
}

func (s *PostmortemService) List(ctx context.Context, params ListPostmortemsParams) ([]PostmortemListItem, error) {
	res, err := s.client.raw.ListPostmortems(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list postmortems: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListPostmortemsOK:
		return res.Postmortems, nil
	case *gen.ListPostmortemsUnauthorized:
		return nil, fmt.Errorf("waroom: list postmortems: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListPostmortemsPaymentRequired:
		return nil, fmt.Errorf("waroom: list postmortems: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list postmortems: unhandled response %T", res)
	}
}

func (s *PostmortemService) Create(ctx context.Context, req CreatePostmortemReq) (*CreatedPostmortem, error) {
	res, err := s.client.raw.CreatePostmortem(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("waroom: create postmortem: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreatePostmortemOK:
		return res, nil
	case *gen.CreatePostmortemUnauthorized:
		return nil, fmt.Errorf("waroom: create postmortem: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreatePostmortemPaymentRequired:
		return nil, fmt.Errorf("waroom: create postmortem: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreatePostmortemUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create postmortem: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: create postmortem: unhandled response %T", res)
	}
}

func (s *PostmortemService) Get(ctx context.Context, params GetPostmortemParams) (*Postmortem, error) {
	res, err := s.client.raw.GetPostmortem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get postmortem: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetPostmortemOK:
		return &res.Postmortem, nil
	case *gen.GetPostmortemNotFound:
		return nil, fmt.Errorf("waroom: get postmortem: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetPostmortemUnauthorized:
		return nil, fmt.Errorf("waroom: get postmortem: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetPostmortemPaymentRequired:
		return nil, fmt.Errorf("waroom: get postmortem: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get postmortem: unhandled response %T", res)
	}
}

func (s *PostmortemService) Update(ctx context.Context, req UpdatePostmortemReq, params UpdatePostmortemParams) (*UpdatedPostmortem, error) {
	res, err := s.client.raw.UpdatePostmortem(ctx, &req, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update postmortem: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdatePostmortemOK:
		return &res.Postmortem, nil
	case *gen.UpdatePostmortemNotFound:
		return nil, fmt.Errorf("waroom: update postmortem: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdatePostmortemUnauthorized:
		return nil, fmt.Errorf("waroom: update postmortem: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdatePostmortemPaymentRequired:
		return nil, fmt.Errorf("waroom: update postmortem: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdatePostmortemUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update postmortem: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update postmortem: unhandled response %T", res)
	}
}

func (s *PostmortemService) Delete(ctx context.Context, params DeletePostmortemParams) error {
	res, err := s.client.raw.DeletePostmortem(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete postmortem: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeletePostmortemNoContent:
		return nil
	case *gen.DeletePostmortemNotFound:
		return fmt.Errorf("waroom: delete postmortem: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeletePostmortemUnauthorized:
		return fmt.Errorf("waroom: delete postmortem: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeletePostmortemPaymentRequired:
		return fmt.Errorf("waroom: delete postmortem: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete postmortem: unhandled response %T", res)
	}
}

func (s *PostmortemService) GetTemplate(ctx context.Context) (*PostmortemDefaultTemplate, error) {
	res, err := s.client.raw.GetPostmortemTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("waroom: get postmortem template: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetPostmortemTemplateOK:
		return res, nil
	case *gen.GetPostmortemTemplateUnauthorized:
		return nil, fmt.Errorf("waroom: get postmortem template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetPostmortemTemplatePaymentRequired:
		return nil, fmt.Errorf("waroom: get postmortem template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get postmortem template: unhandled response %T", res)
	}
}

func (s *PostmortemService) UpdateTemplate(ctx context.Context, req UpdatePostmortemTemplateReq) (*UpdatedPostmortemDefaultTemplate, error) {
	res, err := s.client.raw.UpdatePostmortemTemplate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("waroom: update postmortem template: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdatePostmortemTemplateOK:
		return res, nil
	case *gen.UpdatePostmortemTemplateUnauthorized:
		return nil, fmt.Errorf("waroom: update postmortem template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdatePostmortemTemplatePaymentRequired:
		return nil, fmt.Errorf("waroom: update postmortem template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: update postmortem template: unhandled response %T", res)
	}
}
