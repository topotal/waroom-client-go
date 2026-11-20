package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// OnCallTeamService handles communication with the on-call team related
// methods of the Waroom API. Use it via Client.OnCallTeams.
type OnCallTeamService service

// List returns the on-call teams with their rotations and current on-duty
// participants.
func (s *OnCallTeamService) List(ctx context.Context, params ListOnCallTeamsParams) ([]OnCallTeamListItem, error) {
	res, err := s.client.raw.ListOnCallTeams(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list on-call teams: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListOnCallTeamsOK:
		return res.Teams, nil
	case *gen.ListOnCallTeamsUnauthorized:
		return nil, fmt.Errorf("waroom: list on-call teams: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListOnCallTeamsPaymentRequired:
		return nil, fmt.Errorf("waroom: list on-call teams: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list on-call teams: unhandled response %T", res)
	}
}

// Get returns the on-call team identified by params.UUID.
func (s *OnCallTeamService) Get(ctx context.Context, params GetOnCallTeamParams) (*OnCallTeam, error) {
	res, err := s.client.raw.GetOnCallTeam(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get on-call team: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetOnCallTeamOK:
		return &res.Team, nil
	case *gen.GetOnCallTeamNotFound:
		return nil, fmt.Errorf("waroom: get on-call team: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetOnCallTeamUnauthorized:
		return nil, fmt.Errorf("waroom: get on-call team: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetOnCallTeamPaymentRequired:
		return nil, fmt.Errorf("waroom: get on-call team: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get on-call team: unhandled response %T", res)
	}
}

// Create creates an on-call team with its rotations. req.Rotations must
// contain at least one rotation.
func (s *OnCallTeamService) Create(ctx context.Context, req CreateOnCallTeamReq) (*CreatedOnCallTeam, error) {
	res, err := s.client.raw.CreateOnCallTeam(ctx, &gen.CreateOnCallTeamReq{Team: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: create on-call team: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateOnCallTeamCreated:
		return &res.Team, nil
	case *gen.CreateOnCallTeamBadRequest:
		return nil, fmt.Errorf("waroom: create on-call team: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.CreateOnCallTeamUnauthorized:
		return nil, fmt.Errorf("waroom: create on-call team: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateOnCallTeamPaymentRequired:
		return nil, fmt.Errorf("waroom: create on-call team: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreateOnCallTeamUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create on-call team: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: create on-call team: unhandled response %T", res)
	}
}

// Update updates only the fields set in req. When req.Rotations is set, it
// replaces all the rotations of the team: a rotation with a UUID is updated,
// one without is added, and rotations missing from the slice are deleted.
func (s *OnCallTeamService) Update(ctx context.Context, req UpdateOnCallTeamReq, params UpdateOnCallTeamParams) (*UpdatedOnCallTeam, error) {
	res, err := s.client.raw.UpdateOnCallTeam(ctx, &gen.UpdateOnCallTeamReq{Team: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update on-call team: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateOnCallTeamOK:
		return &res.Team, nil
	case *gen.UpdateOnCallTeamBadRequest:
		return nil, fmt.Errorf("waroom: update on-call team: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallTeamNotFound:
		return nil, fmt.Errorf("waroom: update on-call team: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallTeamUnauthorized:
		return nil, fmt.Errorf("waroom: update on-call team: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateOnCallTeamPaymentRequired:
		return nil, fmt.Errorf("waroom: update on-call team: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateOnCallTeamUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update on-call team: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: update on-call team: unhandled response %T", res)
	}
}

// Delete deletes the on-call team and its rotations. A team used by a step
// of an escalation policy cannot be deleted and results in a 422 APIError.
func (s *OnCallTeamService) Delete(ctx context.Context, params DeleteOnCallTeamParams) error {
	res, err := s.client.raw.DeleteOnCallTeam(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete on-call team: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteOnCallTeamNoContent:
		return nil
	case *gen.DeleteOnCallTeamNotFound:
		return fmt.Errorf("waroom: delete on-call team: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.DeleteOnCallTeamUnauthorized:
		return fmt.Errorf("waroom: delete on-call team: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteOnCallTeamPaymentRequired:
		return fmt.Errorf("waroom: delete on-call team: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.DeleteOnCallTeamUnprocessableEntity:
		return fmt.Errorf("waroom: delete on-call team: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return fmt.Errorf("waroom: delete on-call team: unhandled response %T", res)
	}
}
