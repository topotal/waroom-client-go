package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
)

type LabelService struct {
	client *WaroomClient
}

func NewLabelService(client *WaroomClient) *LabelService {
	return &LabelService{client: client}
}

func (s *LabelService) ListServiceLabels(ctx context.Context, params gen.ListServiceLabelsParams) ([]gen.Label, error) {
	res, err := s.client.raw.ListServiceLabels(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "List Service Error") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.ListServiceLabelsOK:
		return res.Labels, nil
	case *gen.UnauthorizedError:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

func (s *LabelService) CreateServiceLabel(ctx context.Context, req gen.Label, params gen.CreateServiceLabelParams) (*gen.Label, error) {
	res, err := s.client.raw.CreateServiceLabel(ctx, &req, params)
	if err != nil {
		return nil, errors.Wrap(err, "Create service error")
	}

	switch res := res.(type) {
	case *gen.CreateServiceLabelCreated:
		return &req, nil
	case *gen.UnauthorizedError:
		return nil, errors.Wrap(err, "unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}

func (s *LabelService) DeleteServiceLabel(ctx context.Context, params gen.DeleteServiceLabelParams) error {
	res, err := s.client.raw.DeleteServiceLabel(ctx, params)
	if err != nil {
		return errors.Wrap(err, "err") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.DeleteServiceLabelNoContent:
		return nil
	case *gen.UnauthorizedError:
		return errors.New("unauthorized")
	case *gen.NotFoundError:
		return errors.Errorf("not found: name=%s, uuid=%s", params.Name, params.UUID)
	default:
		return errors.Errorf("unhandled response %T", res)
	}
}
