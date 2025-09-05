package waroom

import (
	"context"

	gen "waroom-client-go/internal/ogen"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

type IncidentsService struct{ client *Client }

// GET /incidents
// func (svc *IncidentsService) List(ctx context.Context, page, perPage int, sort string) ([]IncidentSummary, *Pagination, error) {
// 	params := gen.IncidentsGetParams{}
// 	res, err := svc.client.raw.IncidentsGet(ctx, params)
//
// 	if err != nil {
// 		return nil, nil, errors.Wrap(err, "err")
// 	}
//
// 	return res.Incidents, &res.ResponseMetadata, nil
// }

// func (c *Client) GetInsidents(ctx context.Context, page, perPage int, sort string) ([]IncidentSummary, *Pagination, error) {
// 	params := gen.IncidentsGetParams{}
// 	res, err := c.raw.IncidentsGet(ctx, params)
//
// 	if err != nil {
// 		return nil, nil, errors.Wrap(err, "err")
// 	}
//
// 	return res.Incidents, &res.ResponseMetadata, nil
// }

// GET /incidents/{uuid}
func (c *Client) GetIncident(ctx context.Context, uuid uuid.UUID) (*gen.IncidentsIncidentUUIDGetRes, error) {
	params := gen.IncidentsIncidentUUIDGetParams{IncidentUUID: uuid}

	res, err := c.raw.IncidentsIncidentUUIDGet(ctx, params)

	if err != nil {
		return nil, errors.Wrap(err, "err")
	}

	return &res, nil
}

// // 便利メソッド例（下は概念。生成された名前に合わせて薄くラップ）
// func (c *Client) CreateIncident(ctx context.Context, in CreateIncidentRequest) (*Incident, error) {
// 	out, err := c.raw.CreateIncident(ctx, gen.CreateIncidentRequest(in))
// 	return (*Incident)(out), err
// }

//
// // POST /incidents
// func (s *IncidentsService) Create(ctx context.Context, in CreateIncidentRequest) (*Incident, error) {
// 	out, err := s.c.raw.CreateIncident(ctx, gen.CreateIncidentRequest(in))
// 	if err != nil {
// 		return nil, wrapHTTPError(err)
// 	}
// 	return (*Incident)(out), nil
// }
//
// // GET /incidents/{uuid}
// func (s *IncidentsService) Get(ctx context.Context, uuid string) (*Incident, error) {
// 	out, err := s.c.raw.GetIncident(ctx, uuid)
// 	if err != nil {
// 		return nil, wrapHTTPError(err)
// 	}
// 	return (*Incident)(out), nil
// }
