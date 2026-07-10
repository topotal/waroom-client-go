package waroom

import (
	"context"

	"github.com/go-faster/errors"

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
		return nil, errors.Wrap(err, "list sercice labels")
	}

	switch res := res.(type) {
	case *gen.ListServiceLabelsOK:
		return res.Labels, nil
	case *gen.ListServiceLabelsUnauthorized:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

func (s *LabelService) CreateServiceLabel(ctx context.Context, req CreateServiceLabelReq, params CreateServiceLabelParams) (*CreatedLabel, error) {
	res, err := s.client.raw.CreateServiceLabel(ctx, &gen.CreateServiceLabelReq{Label: req}, params)
	if err != nil {
		return nil, errors.Wrap(err, "create service label")
	}

	switch res := res.(type) {
	case *gen.CreateServiceLabelCreated:
		return &res.Label, nil
	case *gen.CreateServiceLabelUnauthorized:
		return nil, errors.Wrap(err, "unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

func (s *LabelService) UpdateServiceLabel(ctx context.Context, req UpdateServiceLabelReq, params UpdateServiceLabelParams) (*UpdatedLabel, error) {
	res, err := s.client.raw.UpdateServiceLabel(ctx, &gen.UpdateServiceLabelReq{Label: req}, params)
	if err != nil {
		return nil, errors.Wrap(err, "update service label")
	}

	switch res := res.(type) {
	case *gen.UpdateServiceLabelOK:
		return &res.Label, nil
	case *gen.UpdateServiceLabelUnauthorized:
		return nil, errors.Wrap(err, "unauthorized")
	case *gen.UpdateServiceLabelNotFound:
		return nil, errors.Errorf("not found: service=%s, uuid=%s", params.Name, params.UUID)
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

func (s *LabelService) DeleteServiceLabel(ctx context.Context, params DeleteServiceLabelParams) error {
	res, err := s.client.raw.DeleteServiceLabel(ctx, params)
	if err != nil {
		return errors.Wrap(err, "delete service label")
	}

	switch res := res.(type) {
	case *gen.DeleteServiceLabelNoContent:
		return nil
	case *gen.DeleteServiceLabelUnauthorized:
		return errors.New("unauthorized")
	case *gen.DeleteServiceLabelNotFound:
		return errors.Errorf("not found: service=%s, uuid=%s", params.Name, params.UUID)
	default:
		return errors.Errorf("unhandled response %T", res)
	}
}
