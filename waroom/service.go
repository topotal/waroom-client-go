package waroom

import (
	"context"

	gen "github.com/topotal/waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
)

type ServiceService struct {
	client *WaroomClient
}

func NewServiceService(client *WaroomClient) *ServiceService {
	return &ServiceService{client: client}
}

func (s *ServiceService) ListServices(ctx context.Context, params gen.ListServicesParams) ([]gen.Service, error) {
	res, err := s.client.raw.ListServices(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err") //TODO ちゃんとする
	}

	switch res := res.(type) {
	case *gen.ServicesListResponse:
		return res.Services, nil
	case *gen.UnauthorizedError:
		return nil, errors.New("unauthorized")
	default:
		return nil, errors.Errorf("unhandled response %T", res)
	}
}
