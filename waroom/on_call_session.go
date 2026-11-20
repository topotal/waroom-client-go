package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// OnCallSessionService handles communication with the on-call session related
// methods of the Waroom API. Use it via Client.OnCallSessions.
type OnCallSessionService service

func (s *OnCallSessionService) List(ctx context.Context, params ListOnCallSessionsParams) ([]OnCallSessionListItem, error) {
	res, err := s.client.raw.ListOnCallSessions(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list on-call sessions: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListOnCallSessionsOK:
		return res.Sessions, nil
	case *gen.ListOnCallSessionsBadRequest:
		return nil, fmt.Errorf("waroom: list on-call sessions: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.ListOnCallSessionsUnauthorized:
		return nil, fmt.Errorf("waroom: list on-call sessions: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListOnCallSessionsPaymentRequired:
		return nil, fmt.Errorf("waroom: list on-call sessions: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list on-call sessions: unhandled response %T", res)
	}
}

func (s *OnCallSessionService) Get(ctx context.Context, params GetOnCallSessionParams) (*OnCallSession, error) {
	res, err := s.client.raw.GetOnCallSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get on-call session: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetOnCallSessionOK:
		return &res.Session, nil
	case *gen.GetOnCallSessionNotFound:
		return nil, fmt.Errorf("waroom: get on-call session: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetOnCallSessionUnauthorized:
		return nil, fmt.Errorf("waroom: get on-call session: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetOnCallSessionPaymentRequired:
		return nil, fmt.Errorf("waroom: get on-call session: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get on-call session: unhandled response %T", res)
	}
}
