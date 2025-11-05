package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
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

func (s *LabelService) CreateServiceLabel(ctx context.Context, req Label, params CreateServiceLabelParams) (*Label, error) {
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

func (s *LabelService) DeleteServiceLabel(ctx context.Context, params DeleteServiceLabelParams) error {
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
