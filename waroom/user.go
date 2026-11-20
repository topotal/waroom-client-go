package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// UserService handles communication with the user related
// methods of the Waroom API. Use it via Client.Users.
type UserService service

// NewUserService returns client.Users.
//
// Deprecated: Use Client.Users instead.
func NewUserService(client *Client) *UserService {
	return client.Users
}

func (s *UserService) List(ctx context.Context, params ListUsersParams) ([]User, error) {
	res, err := s.client.raw.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list users: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListUsersOK:
		return res.Users, nil
	case *gen.ListUsersUnauthorized:
		return nil, fmt.Errorf("waroom: list users: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListUsersPaymentRequired:
		return nil, fmt.Errorf("waroom: list users: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list users: unhandled response %T", res)
	}
}

func (s *UserService) Get(ctx context.Context, params GetUserParams) (*UserDetail, error) {
	res, err := s.client.raw.GetUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: get user: %w", err)
	}

	switch res := res.(type) {
	case *gen.GetUserOK:
		return &res.User, nil
	case *gen.GetUserNotFound:
		return nil, fmt.Errorf("waroom: get user: %w", &APIError{StatusCode: http.StatusNotFound, Messages: optMessages(res.Error)})
	case *gen.GetUserUnauthorized:
		return nil, fmt.Errorf("waroom: get user: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.GetUserPaymentRequired:
		return nil, fmt.Errorf("waroom: get user: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: get user: unhandled response %T", res)
	}
}

func (s *UserService) Invite(ctx context.Context, req InviteUsersReq) ([]InvitedUser, error) {
	res, err := s.client.raw.InviteUsers(ctx, &gen.InviteUsersReq{Invitation: req})
	if err != nil {
		return nil, fmt.Errorf("waroom: invite users: %w", err)
	}

	switch res := res.(type) {
	case *gen.InviteUsersAccepted:
		return res.Users, nil
	case *gen.InviteUsersBadRequest:
		return nil, fmt.Errorf("waroom: invite users: %w", &APIError{StatusCode: http.StatusBadRequest, Messages: optMessages(res.Error)})
	case *gen.InviteUsersUnauthorized:
		return nil, fmt.Errorf("waroom: invite users: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.InviteUsersPaymentRequired:
		return nil, fmt.Errorf("waroom: invite users: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	case *gen.InviteUsersConflict:
		return nil, fmt.Errorf("waroom: invite users: %w", &APIError{StatusCode: http.StatusConflict, Messages: res.Errors})
	case *gen.InviteUsersUnprocessableEntity:
		return nil, fmt.Errorf("waroom: invite users: %w", &APIError{StatusCode: http.StatusUnprocessableEntity, Messages: res.Errors})
	case *gen.InviteUsersTooManyRequests:
		return nil, fmt.Errorf("waroom: invite users: %w", &APIError{StatusCode: http.StatusTooManyRequests})
	default:
		return nil, fmt.Errorf("waroom: invite users: unhandled response %T", res)
	}
}
