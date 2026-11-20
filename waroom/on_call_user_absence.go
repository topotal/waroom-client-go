package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// OnCallUserAbsenceService handles communication with the on-call user
// absence related methods of the Waroom API. Use it via Client.OnCallUserAbsences.
type OnCallUserAbsenceService service

// List returns the on-call absences of the user identified by params.Sub,
// ordered by start time. By default only active absences are returned; set
// params.IncludeExpired to also return finished one-off absences.
func (s *OnCallUserAbsenceService) List(ctx context.Context, params ListOnCallUserAbsencesParams) ([]OnCallUserAbsenceListItem, error) {
	res, err := s.client.raw.ListOnCallUserAbsences(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list absences of a user: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListOnCallUserAbsencesOK:
		return res.Absences, nil
	case *gen.ListOnCallUserAbsencesUnauthorized:
		return nil, fmt.Errorf("waroom: list absences of a user: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListOnCallUserAbsencesPaymentRequired:
		return nil, fmt.Errorf("waroom: list absences of a user: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.ListOnCallUserAbsencesNotFound:
		return nil, fmt.Errorf("waroom: list absences of a user: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: list absences of a user: unhandled response %T", res)
	}
}

// Get returns an on-call absence of a user.
func (s *OnCallUserAbsenceService) Get(ctx context.Context, params GetOnCallUserAbsenceParams) (*OnCallUserAbsence, error) {
	res, err := s.client.raw.GetOnCallUserAbsence(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get an absence of a user: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetOnCallUserAbsenceOK:
		return &res.Absence, nil
	case *gen.GetOnCallUserAbsenceUnauthorized:
		return nil, fmt.Errorf("waroom: get an absence of a user: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetOnCallUserAbsencePaymentRequired:
		return nil, fmt.Errorf("waroom: get an absence of a user: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.GetOnCallUserAbsenceNotFound:
		return nil, fmt.Errorf("waroom: get an absence of a user: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	default:
		return nil, fmt.Errorf("waroom: get an absence of a user: unhandled response %T", res)
	}
}

// Create registers an on-call absence of a user. While the absence lasts,
// the user's turns in every rotation are passed to the next person.
func (s *OnCallUserAbsenceService) Create(ctx context.Context, req CreateOnCallUserAbsenceReq, params CreateOnCallUserAbsenceParams) (*CreatedOnCallUserAbsence, error) {
	res, err := s.client.raw.CreateOnCallUserAbsence(ctx, &gen.CreateOnCallUserAbsenceReq{Absence: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: create an absence of a user: %w", err)
	}

	switch res := res.(type) {
	case *gen.CreateOnCallUserAbsenceCreated:
		return &res.Absence, nil
	case *gen.CreateOnCallUserAbsenceBadRequest:
		return nil, fmt.Errorf("waroom: create an absence of a user: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.CreateOnCallUserAbsenceUnauthorized:
		return nil, fmt.Errorf("waroom: create an absence of a user: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.CreateOnCallUserAbsencePaymentRequired:
		return nil, fmt.Errorf("waroom: create an absence of a user: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.CreateOnCallUserAbsenceNotFound:
		return nil, fmt.Errorf("waroom: create an absence of a user: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.CreateOnCallUserAbsenceUnprocessableEntity:
		return nil, fmt.Errorf("waroom: create an absence of a user: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: create an absence of a user: unhandled response %T", res)
	}
}

// Update updates an on-call absence of a user. Only the fields set in req
// are changed.
func (s *OnCallUserAbsenceService) Update(ctx context.Context, req UpdateOnCallUserAbsenceReq, params UpdateOnCallUserAbsenceParams) (*UpdatedOnCallUserAbsence, error) {
	res, err := s.client.raw.UpdateOnCallUserAbsence(ctx, &gen.UpdateOnCallUserAbsenceReq{Absence: req}, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: update an absence of a user: %w", err)
	}

	switch res := res.(type) {
	case *gen.UpdateOnCallUserAbsenceOK:
		return &res.Absence, nil
	case *gen.UpdateOnCallUserAbsenceBadRequest:
		return nil, fmt.Errorf("waroom: update an absence of a user: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallUserAbsenceUnauthorized:
		return nil, fmt.Errorf("waroom: update an absence of a user: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.UpdateOnCallUserAbsencePaymentRequired:
		return nil, fmt.Errorf("waroom: update an absence of a user: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.UpdateOnCallUserAbsenceNotFound:
		return nil, fmt.Errorf("waroom: update an absence of a user: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.UpdateOnCallUserAbsenceUnprocessableEntity:
		return nil, fmt.Errorf("waroom: update an absence of a user: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	default:
		return nil, fmt.Errorf("waroom: update an absence of a user: unhandled response %T", res)
	}
}

// Delete deletes an on-call absence of a user.
func (s *OnCallUserAbsenceService) Delete(ctx context.Context, params DeleteOnCallUserAbsenceParams) error {
	res, err := s.client.raw.DeleteOnCallUserAbsence(ctx, params)
	if err != nil {
		return fmt.Errorf("waroom: delete an absence of a user: %w", err)
	}

	switch res := res.(type) {
	case *gen.DeleteOnCallUserAbsenceNoContent:
		return nil
	case *gen.DeleteOnCallUserAbsenceUnauthorized:
		return fmt.Errorf("waroom: delete an absence of a user: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.DeleteOnCallUserAbsencePaymentRequired:
		return fmt.Errorf("waroom: delete an absence of a user: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.DeleteOnCallUserAbsenceNotFound:
		return fmt.Errorf("waroom: delete an absence of a user: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	default:
		return fmt.Errorf("waroom: delete an absence of a user: unhandled response %T", res)
	}
}
