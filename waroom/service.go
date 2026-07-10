package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
)

type ServiceService struct {
	client *Client
}

func NewServiceService(client *Client) *ServiceService {
	return &ServiceService{client: client}
}

func (s *ServiceService) ListServices(ctx context.Context, params ListServicesParams) ([]Service, error) {
	res, err := s.client.raw.ListServices(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.ListServicesOK:
		return res.Services, nil
	case *gen.ListServicesUnauthorized:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}
