package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// PostmortemTemplateService handles communication with the postmortem template related
// methods of the Waroom API. Use it via Client.PostmortemTemplates.
type PostmortemTemplateService service

// NewPostmortemTemplateService returns client.PostmortemTemplates.
//
// Deprecated: Use Client.PostmortemTemplates instead.
func NewPostmortemTemplateService(client *Client) *PostmortemTemplateService {
	return client.PostmortemTemplates
}

func (s *PostmortemTemplateService) List(ctx context.Context, params ListPostmortemTemplatesParams) ([]PostmortemTemplateListItem, error) {
	res, err := s.client.raw.ListPostmortemTemplates(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list postmortem templates: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListPostmortemTemplatesOK:
		return res.PostmortemTemplates, nil
	case *gen.ListPostmortemTemplatesUnauthorized:
		return nil, fmt.Errorf("waroom: list postmortem templates: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListPostmortemTemplatesPaymentRequired:
		return nil, fmt.Errorf("waroom: list postmortem templates: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list postmortem templates: unhandled response %T", res)
	}
}

func (s *PostmortemTemplateService) Create(ctx context.Context, req CreatePostmortemTemplateReq) (*CreatedPostmortemTemplate, error) {
	res, err := s.client.raw.CreatePostmortemTemplate(ctx, &gen.CreatePostmortemTemplateReq{PostmortemTemplate: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: create postmortem template: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreatePostmortemTemplateCreated:
		return &res.PostmortemTemplate, nil
	case *gen.CreatePostmortemTemplateUnauthorized:
		return nil, fmt.Errorf("waroom: create postmortem template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreatePostmortemTemplatePaymentRequired:
		return nil, fmt.Errorf("waroom: create postmortem template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreatePostmortemTemplateUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create postmortem template: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: create postmortem template: unhandled response %T", res)
	}
}

func (s *PostmortemTemplateService) Get(ctx context.Context, params GetPostmortemTemplateByUuidParams) (*PostmortemTemplate, error) {
	res, err := s.client.raw.GetPostmortemTemplateByUuid(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get postmortem template: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetPostmortemTemplateByUuidOK:
		return &res.PostmortemTemplate, nil
	case *gen.GetPostmortemTemplateByUuidNotFound:
		return nil, fmt.Errorf("waroom: get postmortem template: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetPostmortemTemplateByUuidUnauthorized:
		return nil, fmt.Errorf("waroom: get postmortem template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetPostmortemTemplateByUuidPaymentRequired:
		return nil, fmt.Errorf("waroom: get postmortem template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get postmortem template: unhandled response %T", res)
	}
}

func (s *PostmortemTemplateService) Update(ctx context.Context, req UpdatePostmortemTemplateByUuidReq, params UpdatePostmortemTemplateByUuidParams) (*UpdatedPostmortemTemplate, error) {
	res, err := s.client.raw.UpdatePostmortemTemplateByUuid(ctx, &gen.UpdatePostmortemTemplateByUuidReq{PostmortemTemplate: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update postmortem template: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdatePostmortemTemplateByUuidOK:
		return &res.PostmortemTemplate, nil
	case *gen.UpdatePostmortemTemplateByUuidNotFound:
		return nil, fmt.Errorf("waroom: update postmortem template: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdatePostmortemTemplateByUuidUnauthorized:
		return nil, fmt.Errorf("waroom: update postmortem template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdatePostmortemTemplateByUuidPaymentRequired:
		return nil, fmt.Errorf("waroom: update postmortem template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdatePostmortemTemplateByUuidUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update postmortem template: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: update postmortem template: unhandled response %T", res)
	}
}

func (s *PostmortemTemplateService) Delete(ctx context.Context, params DeletePostmortemTemplateByUuidParams) error {
	res, err := s.client.raw.DeletePostmortemTemplateByUuid(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete postmortem template: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeletePostmortemTemplateByUuidNoContent:
		return nil
	case *gen.DeletePostmortemTemplateByUuidNotFound:
		return fmt.Errorf("waroom: delete postmortem template: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeletePostmortemTemplateByUuidUnauthorized:
		return fmt.Errorf("waroom: delete postmortem template: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeletePostmortemTemplateByUuidPaymentRequired:
		return fmt.Errorf("waroom: delete postmortem template: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return fmt.Errorf("waroom: delete postmortem template: unhandled response %T", res)
	}
}
